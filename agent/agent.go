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

	"github.com/ChoTotOSS/fluent2gelf/entry"
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
	in        chan chan *entry.Entry
}

func New(cfg Config) *Agent {
	expr := strings.Replace(cfg.Match, "*", ".*", -1)

	rx, err := regexp.Compile(expr)

	if err != nil {
		panic(err)
	}

	var f func([]byte) bool
	if cfg.Multiline != nil {

		switch {
		case len(cfg.Multiline.Regexp) > 0:
			rx, err := regexp.Compile(cfg.Multiline.Regexp)
			if err != nil {
				panic(err)
			}

			f = func(log []byte) bool {
				return rx.Match(log)
			}

		case cfg.Multiline.IndexLT > 0 && len(cfg.Multiline.Begin) > 0:
			maxLengthToCheck := len(cfg.Multiline.Begin) + cfg.Multiline.IndexLT
			sep := []byte(cfg.Multiline.Begin)
			f = func(log []byte) bool {
				//should handle error
				if len(log) > maxLengthToCheck {
					return bytes.Index(log[:maxLengthToCheck], sep) >= 0
				}
				return bytes.Index(log, sep) >= 0
			}
		case len(cfg.Multiline.Begin) > 0:
			sep := []byte(cfg.Multiline.Begin)
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
		Multiline: cfg.Multiline != nil,
		firstline: f,
		chunks:    make([][]byte, 0),
		in:        make(chan chan *entry.Entry),
	}
}

func (a *Agent) appendChunk(chunk []byte) {
	a.chunks = append(a.chunks, chunk)
}

func (a *Agent) Ship(logs chan *entry.Entry) {
	a.in <- logs
}

func (a *Agent) Run(done chan bool) {

	_gelf := gelf.New("Start agent for: "+a.Match.String(), time.Now().Unix(), 1, gelf.DefaultHost)

	commit := func() {
		chunks := _gelf.ToChunks()

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
		_gelf = gelf.New(m["log"], t.Unix(), streamLevels[m["stream"]], gelf.DefaultHost)
	}

	for {
		select {
		case entries := <-a.in:
			for e := range entries {
				var logm map[string]string
				err := json.Unmarshal(e.Log, &logm)
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
