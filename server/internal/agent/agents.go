package agent

import (
	"fmt"
	"sync"
)

var (
	IDDuplication = fmt.Errorf("ID duplication")
)

// Agents - 所有的被控者的信息
type Agents struct {
	Agents map[string]*Agent
	m      sync.Mutex
}

// GetAgentById - 根据id获取Agent
func (a *Agents) GetAgentById(id string) *Agent {
	a.m.Lock()
	defer a.m.Unlock()
	if agent, ok := a.Agents[id]; ok {
		return agent
	} else {
		return nil
	}
}

// Add - 添加一个agent到agents
func (a *Agents) Add(id string, agent *Agent) (err error) {
	a.m.Lock()
	defer a.m.Unlock()

	if _, ok := a.Agents[id]; ok {
		return IDDuplication
	}

	a.Agents[id] = agent

	return
}

// Delete - 根据id删除一个agent
func (a *Agents) Delete(id string) {
	a.m.Lock()
	defer a.m.Unlock()

	delete(a.Agents, id)

	return
}
