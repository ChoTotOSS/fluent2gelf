package gelf

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestChunks(t *testing.T) {
	s := []byte(randRunes(160000))
	chunks := NewChunks(s)

	next := chunks.Next()

	t.Run("Chunks should be hasNext", func(t *testing.T) {
		if chunks.HasNext() != true {
			t.Fail()
		}
	})

	t.Run("chunks lenght should less than 2 * 16", func(t *testing.T) {
		if len(next) > 65536 {
			t.Logf("expected = %d, result = %d\n", MAX_CHUNK_SIZE+12, len(next))
			t.Fail()
		}
	})

	t.Run("chunks data should match with runes data", func(t *testing.T) {
		if bytes.Compare(next[12:], s[:len(next)-12]) != 0 {
			t.Fail()
		}
	})

	t.Run("After all, chunks does not has next", func(t *testing.T) {
		chunks.Next()
		chunks.Next()
		if chunks.HasNext() == true {
			t.Fail()
		}
	})
}

func TestLittleChunks(t *testing.T) {
	data := []byte(randRunes(MAX_CHUNK_SIZE))
	c := NewChunks(data)

	t.Run("1. Chunks must has next", func(t *testing.T) {
		if c.HasNext() == false {
			t.Fail()
		}
	})
	t.Run("2. Chunks should return only a chunk", func(t *testing.T) {
		c.Next()
		if c.HasNext() != false {
			t.Fail()
		}
	})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghij klmnopqrstuvwxyz ABCDEFGHIJKLMN OPQRSTUVWXYZ")

func randRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
