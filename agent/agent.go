package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/ChoTotOSS/fluent2gelf/gelf"
	"github.com/duythinht/zaptor"
	"go.uber.org/zap"
)

var logger = zaptor.Default()

var streamLevels = map[string]int{
	"stdout": 1,
	"stderr": 3,
}

type Agent struct {
	Match     *regexp.Regexp
	Host      string
	Port      int
	Multiline bool
	firstline func([]byte) bool
	chunks    [][]byte //stupid implements, thinking later
	in        chan chan []byte
}

func New(cfg AgentConfig) *Agent {
	expr := strings.Replace(cfg.Match, "*", ".*", -1)

	rx, err := regexp.Compile(expr)

	if err != nil {
		panic(err)
	}

	if cfg.Multiline && cfg.Firstline == nil {
		panic(errors.New("Agent config is invalid, multiline support should specs firstline match"))
	}

	var f func([]byte) bool
	if cfg.Multiline {

		switch {
		case len(cfg.Firstline.Regexp) > 0:
			rx, err := regexp.Compile(cfg.Firstline.Regexp)
			if err != nil {
				panic(err)
			}

			f = func(log []byte) bool {
				return rx.Match(log)
			}

		case cfg.Firstline.IndexLT > 0 && len(cfg.Firstline.Begin) > 0:
			maxLengthToCheck := len(cfg.Firstline.Begin) + cfg.Firstline.IndexLT
			sep := []byte(cfg.Firstline.Begin)
			f = func(log []byte) bool {
				//should handle error
				if len(log) > maxLengthToCheck {
					return bytes.Index(log[:maxLengthToCheck], sep) >= 0
				}
				return bytes.Index(log, sep) >= 0
			}
		case len(cfg.Firstline.Begin) > 0:
			sep := []byte(cfg.Firstline.Begin)
			f = func(log []byte) bool {
				return bytes.HasPrefix(log, sep)
			}
		default:
			panic(errors.New("Missing firstline match"))
		}
	}

	return &Agent{
		Match:     rx,
		Host:      cfg.Host,
		Port:      cfg.Port,
		Multiline: cfg.Multiline,
		firstline: f,
		chunks:    make([][]byte, 0),
		in:        make(chan chan []byte),
	}
}

func (a *Agent) appendChunk(chunk []byte) {
	a.chunks = append(a.chunks, chunk)
}

func (a *Agent) Ship(logs chan []byte) {
	a.in <- logs
}

func (a *Agent) Run(done chan bool) {

	_gelf := gelf.CreateGelf("Start agent for: "+a.Match.String(), time.Now().Unix(), 1)

	commit := func() {
		zipped := gelf.ZipMessage(_gelf.ToJSON())
		chunks := gelf.NewChunks(zipped)

		for chunks.HasNext() {
			a.appendChunk(chunks.Next())
		}
	}

	createGelf := func(m map[string]string) {
		t, err := time.Parse(time.RFC3339, m["time"])
		if err != nil {
			logger.Warn("agent#multiline#timestamp",
				zap.Error(err),
				zap.String("time", m["time"]),
			)
		}
		logger.Debug("agent#gelf#create", zap.String("log", m["log"]), zap.Time("time", t))
		_gelf = gelf.CreateGelf(m["log"], t.Unix(), streamLevels[m["stream"]])
	}

	for {
		select {
		case logs := <-a.in:
			logger.Debug("Got logs")
			for log := range logs {
				var logm map[string]string
				err := json.Unmarshal(log, &logm)
				if err != nil {
					continue
				}
				if a.Multiline {
					if a.firstline([]byte(logm["log"])) {
						commit()
						createGelf(logm)
					} else {
						_gelf.Append(logm["log"])
					}
				} else {
					commit()
					createGelf(logm)
				}
			}
			commit()
			a.SendAndReset()
		case <-done:
			return
		}
	}
}

func (a *Agent) SendAndReset() {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", a.Host, a.Port))
	if err != nil {
		logger.Error("Agent can't dial graylog", zap.Error(err))
	}
	defer func() { _ = conn.Close() }()
	logger.Debug("gelf#send chunks", zap.Int("count", len(a.chunks)))
	for _, chunk := range a.chunks {
		n, err := conn.Write(chunk)
		if err != nil {
			logger.Warn("gelf#send", zap.Error(err))
		} else {
			logger.Debug("gelf#send", zap.Int("bytes", n))
		}
	}

	a.chunks = a.chunks[:0]
}
