package gelf

import (
	"bytes"
	"encoding/json"

	"github.com/duythinht/zaptor"
	"go.uber.org/zap"
)

var logger = zaptor.Default()

const (
	GELF_DEFAUL_VERSION = "1.1"
	SIXTY_FOUR_KiB      = 64000 //max chunk size
	THIRTY_TWO_KiB      = 32000 //max size es can handle
)

type Gelf struct {
	Version      string `json:"version"`
	Host         string `json:"host"`
	ShortMessage string `json:"short_message"`
	FullMessage  string `json:"full_message"`
	Timestamp    int64  `json:"timestamp"`
	Level        int    `json:"level"`
	buf          bytes.Buffer
}

func CreateGelf(short string, timestamp int64, level int) *Gelf {
	g := Gelf{
		Version:      GELF_DEFAUL_VERSION,
		Host:         "default",
		ShortMessage: short,
		Timestamp:    timestamp,
		Level:        level,
	}
	g.buf.WriteString(short)
	return &g
}

func (g *Gelf) ToJSON() []byte {
	g.FullMessage = g.buf.String()
	msg, err := json.Marshal(g)
	if err != nil {
		logger.Warn("gelf#tojson", zap.Error(err))
		return msg
	}
	return msg
}

func (g *Gelf) Append(message string) {
	g.buf.WriteString(message)
}

func (g *Gelf) Len() int {
	return g.buf.Len()
}
