package gelf

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGelfMessage(t *testing.T) {
	timestamp := time.Now().Unix()
	message := CreateGelf("hello world", timestamp, 1)

	msg := message.ToJSON()

	x := make(map[string]interface{}, 0)
	_ = json.Unmarshal(msg, &x)
	if x["full_message"] != "hello world" {
		t.Log("Full message should be hello world")
		t.Fail()
	}
}
