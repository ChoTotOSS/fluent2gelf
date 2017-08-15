package agent

import (
	"encoding/json"
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
	Match       *regexp.Regexp
	Host        string
	Port        int
	Multiline   bool
	BeginRecord string
	chunks      [][]byte //stupid implements, thinking later
	in          chan []interface{}
}

func New(cfg AgentConfig) *Agent {
	expr := strings.Replace(cfg.Match, "*", ".*", -1)

	rx, err := regexp.Compile(expr)

	if err != nil {
		panic(err)
	}

	return &Agent{
		Match:       rx,
		Host:        cfg.Host,
		Port:        cfg.Port,
		Multiline:   cfg.Multiline,
		BeginRecord: cfg.BeginRecord,
		chunks:      make([][]byte, 0),
		in:          make(chan []interface{}),
	}
}

func (a *Agent) appendChunk(chunk []byte) {
	a.chunks = append(a.chunks, chunk)
}

func (a *Agent) Put(entries []interface{}) {
	a.in <- entries
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
		case entries := <-a.in:
			for _, entry := range entries {

				e := entry.([]interface{}) // [timestamp, record[string]interface{}
				record := e[1].(map[string]interface{})
				jsonBytes := record["log"].([]byte)

				var m map[string]string
				err := json.Unmarshal(jsonBytes, &m)
				if err != nil {
					logger.Warn("agent#put", zap.Error(err))
					continue
				}

				if a.Multiline {
					if strings.HasPrefix(m["log"], a.BeginRecord) {
						// commit old gelf
						commit()
						//Then Create a new once gelf
						createGelf(m)
					} else {
						_gelf.Append(m["log"])
					}
				} else {
					//commit old and create new once
					commit()
					createGelf(m)
				}
			}
			commit() //commit last message
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

	logger.Debug("len chunks", zap.Int("len", len(a.chunks)))
	for _, chunk := range a.chunks {
		n, err := conn.Write(chunk)
		if err != nil {
			logger.Warn("gelf#send", zap.Error(err))
		} else {
			logger.Debug("gelf#send", zap.Int("n", n))
		}
	}

	a.chunks = a.chunks[:0]
}
