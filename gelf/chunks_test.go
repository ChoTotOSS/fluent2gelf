package gelf

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestChunks(t *testing.T) {
	s := []byte(RandRunes(160000))
	chunks := NewChunks(s)
	next := chunks.Next()
	if len(next) > 65000 {
		t.Log("len(chunk) should < 65k")
		t.Fail()
	}
	if bytes.Compare(next[12:], s[:len(next)-12]) != 0 {
		t.Log("Wrong chunk data")
		t.Fail()
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghij klmnopqrstuvwxyz ABCDEFGHIJKLMN OPQRSTUVWXYZ")

func RandRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
