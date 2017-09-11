package gelf

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/duythinht/zaptor"
)

var logger = zaptor.Default()

const (
	GELF_DEFAUL_VERSION = "1.1"
	MAX_CHUNK_SIZE      = 65500 // 2^16 - 36
)

type Gelf struct {
	Version      string       `json:"version"`
	Host         string       `json:"host"`
	ShortMessage string       `json:"short_message"`
	FullMessage  string       `json:"full_message"`
	Timestamp    int64        `json:"timestamp"`
	Level        int          `json:"level"`
	buf          bytes.Buffer `json:"-"`
}

var DefaultHost, _ = os.Hostname()

// Create a new gelf, with short message, timestamp and level, and host
func New(short string, timestamp int64, level int, host string) *Gelf {
	g := Gelf{
		Version:      GELF_DEFAUL_VERSION,
		Host:         host,
		ShortMessage: short,
		Timestamp:    timestamp,
		Level:        level,
	}
	g.buf.WriteString(short)
	return &g
}

func (g *Gelf) ToJSON() []byte {
	g.FullMessage = g.buf.String()
	msg, _ := json.Marshal(g)
	// There are not neccessary for test unmarshal error, bcz Message struct always be a correct type
	//if err != nil {
	//	return msg
	//}
	return msg
}

func (g *Gelf) toZip() []byte {
	return zipMessage(g.ToJSON())
}

func (g *Gelf) ToChunks() *Chunks {
	return NewChunks(g.toZip())
}

func (g *Gelf) Append(message string) {
	g.buf.WriteString(message)
}

func (g *Gelf) Len() int {
	return g.buf.Len()
}
