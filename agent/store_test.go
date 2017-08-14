package agent

import "testing"

func TestNewAgentStore(t *testing.T) {
	store := NewAgentStore()
	agent1 := New("kube.*_default_*sekai*.log", "agent1", 12204)
	agent2 := New("kube.*_default_*world*.log", "agent2", 12204)
	store.addAgent(agent1)
	store.addAgent(agent2)

	agent := store.Take("kube.abc_default_afdf_sekaidfd.log")
	if agent.Host != agent1.Host {
		t.Fail()
	}

	agentx := store.Take("notinlist")
	if agentx != nil {
		t.Fail()
	}
}
