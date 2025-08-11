package transport

import (
	"fmt"
	"github.com/lzkking/harle/server/internal/agent"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	DefaultHTTPTimeout = time.Minute
)

type ListenerHttpC2 struct {
	Agents agent.Agents //	管理连接到此处的所有被控端的信息
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
		//	处理匿名agent，只做获取其中的公钥的操作，没有的话返回错误码
	} else {
		//
	}

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

	server := &ListenerHttpC2{}

	test := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", 9087),
		Handler:      server.router(),
		WriteTimeout: DefaultHTTPTimeout,
		ReadTimeout:  DefaultHTTPTimeout,
		IdleTimeout:  DefaultHTTPTimeout,
	}
}
