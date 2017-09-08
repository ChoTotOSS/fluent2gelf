package gelf

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGelfMessage(t *testing.T) {
	for i := 1; i < 100; i++ {
		timestamp := time.Now().Unix()
		short := randRunes(10)
		host := randRunes(15)
		message := CreateGelf(short, timestamp, 1, host)

		x := make(map[string]interface{}, 0)
		_ = json.Unmarshal(message.ToJSON(), &x)

		t.Run("full message, should be equal with short message", func(t *testing.T) {
			if x["full_message"] != short {
				t.Logf("expected=%s | result=%s\n", short, x["full_message"])
				t.Fail()
			}
		})

		t.Run("host name should be correct", func(t *testing.T) {
			if x["host"] != host {
				t.Logf("expected=%s | result=%s\n", host, x["host"])
				t.Fail()
			}
		})

		t.Run("timestamp after unmarshal", func(t *testing.T) {
			if int64(x["timestamp"].(float64)) != timestamp {
				t.Logf("expected=%v | result=%v\n", timestamp, int64(x["timestamp"].(float64)))
				t.Fail()
			}
		})

		t.Run("Len should be len of short message", func(t *testing.T) {
			if message.Len() != len(short) {
				t.Fail()
			}
		})
	}
}

func TestGelfAppend(t *testing.T) {

	for i := 0; i < 100; i++ {
		timestamp := time.Now().Unix()
		short := randRunes(10)
		more1 := randRunes(100)
		more2 := randRunes(120)

		message := CreateGelf(short, timestamp, 1, "default")
		message.Append(more1)
		message.Append(more2)

		msg := message.ToJSON()

		x := make(map[string]interface{}, 0)
		_ = json.Unmarshal(msg, &x)
		if x["full_message"] != short+more1+more2 {
			t.Log("Full message is not correct")
			t.Fail()
		}
	}
}

func TestGelfToChunks(t *testing.T) {
	test := func() {
		timestamp := time.Now().Unix()
		short := randRunes(10)
		more1 := randRunes(100)
		more2 := randRunes(120)

		message := CreateGelf(short, timestamp, 1, "default")
		message.Append(more1)
		message.Append(more2)

		c := message.ToChunks()

		if c.HasNext() {
			// ok
			_ = c.Next()
		} else {
			t.Fail()
		}

	}

	test()
}
