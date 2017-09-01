package agent

import "testing"

func TestNewAgentWithMatch(t *testing.T) {

	agent := New(NewConfig("kube.*_default_*sekai*.log", "localhost", 12204, false))
	if !agent.Match.MatchString("kube.abc_default_tangtang_sekailmao.log") {
		t.Fail()
	}
}

func TestResetAgent(t *testing.T) {
	agent := New(NewConfig("kube.*_default_*sekai*.log", "localhost", 12204, false))
	done := make(chan bool)
	go agent.Run(done)
	agent.SendAndReset()

	done <- true
	if len(agent.chunks) != 0 {
		t.Fail()
	}
}
