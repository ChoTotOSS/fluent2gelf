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

func TestGelfAppend(t *testing.T) {
	timestamp := time.Now().Unix()
	message := CreateGelf("hello world", timestamp, 1)
	message.Append(", the world is flat")

	msg := message.ToJSON()

	x := make(map[string]interface{}, 0)
	_ = json.Unmarshal(msg, &x)
	if x["full_message"] != "hello world, the world is flat" {
		t.Log("Full message should be hello world")
		t.Fail()
	}

}
