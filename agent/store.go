package agent

import "sync"

type AgentStore struct {
	AgentList     []*Agent
	AgentQuickMap map[string]*Agent
	mx            *sync.RWMutex
}

func NewAgentStore() *AgentStore {
	return &AgentStore{
		AgentList:     make([]*Agent, 0),
		AgentQuickMap: make(map[string]*Agent, 0),
		mx:            new(sync.RWMutex),
	}
}

func (as *AgentStore) addAgent(agent *Agent) {
	as.AgentList = append(as.AgentList, agent)
}

func (as *AgentStore) Take(tag string) *Agent {
	as.mx.RLock()
	if a, ok := as.AgentQuickMap[tag]; ok {
		defer as.mx.RUnlock()
		return a
	}
	as.mx.RUnlock()
	//Lock to write
	as.mx.Lock()
	defer as.mx.Unlock()
	// Double check
	if _, ok := as.AgentQuickMap[tag]; !ok {
		as.AgentQuickMap[tag] = nil
		for _, agent := range as.AgentList {
			if agent.Match.MatchString(tag) {
				as.AgentQuickMap[tag] = agent
			}
		}
	}

	return as.AgentQuickMap[tag]
}
