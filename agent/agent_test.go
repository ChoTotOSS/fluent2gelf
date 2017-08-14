package agent

import "testing"

func TestNewAgentWithMatch(t *testing.T) {
	agent := New("kube.*_default_*sekai*.log", "localhost", 12204)
	if !agent.Match.MatchString("kube.abc_default_tangtang_sekailmao.log") {
		t.Fail()
	}
}

func TestResetAgent(t *testing.T) {
	agent := New("kube.*_default_*sekai*.log", "localhost", 12204)
	agent.Append([]byte("hello"))
	agent.Append([]byte("hello"))
	agent.Append([]byte("hello"))
	agent.Append([]byte("hello"))
	agent.Append([]byte("hello"))
	agent.Reset()
	if len(agent.chunks) != 0 {
		t.Fail()
	}
}
