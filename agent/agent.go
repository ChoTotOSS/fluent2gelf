package agent

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Agent struct {
	Match  *regexp.Regexp
	Host   string
	Port   int
	chunks [][]byte //stupid implements, thinking later
	in     chan []byte
}

func New(match string, host string, port int) *Agent {
	expr := strings.Replace(match, "*", ".*", -1)

	rx, err := regexp.Compile(expr)

	if err != nil {
		panic(err)
	}

	return &Agent{
		Match:  rx,
		Host:   host,
		Port:   port,
		chunks: make([][]byte, 0),
		in:     make(chan []byte),
	}
}

func (a *Agent) Append(entry []byte) {
	a.chunks = append(a.chunks, entry)
}

func (a *Agent) Run(done chan bool) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case entry := <-a.in:
			a.Append(entry)
		case <-ticker.C:
			a.Reset()
			fmt.Println("Reset agent")
		case <-done:
			break
		}
	}
}

func (a *Agent) Reset() {
	a.chunks = a.chunks[:0]
}
