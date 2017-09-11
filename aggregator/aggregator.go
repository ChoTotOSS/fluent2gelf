package aggregator

import (
	"net"

	"go.uber.org/zap"

	"github.com/ChoTotOSS/fluent2gelf/agent"
	"github.com/ChoTotOSS/fluent2gelf/entry"
	"github.com/ChoTotOSS/fluent2gelf/fluentd"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack"
	"github.com/duythinht/zaptor"
)

var logger = zaptor.Default()

type Aggregator int

func New() Aggregator {
	return Aggregator(0)
}

func (a *Aggregator) Process(conn net.Conn, s *agent.Store) {

	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Warn("Error close connection", zap.Error(err))
		}
	}()
	reader := quickmsgpack.NewReader(conn)

	fw := fluentd.NewForwardReader(reader)

	if tag, ok := fw.ReadTag(); ok {
		count := fw.CountSegments()
		if count > 3 || count < 2 {
			return
		}
		agent := s.Take(string(tag))
		if agent != nil {
			entries := make(chan *entry.Entry)
			agent.Ship(entries)
			defer func() {
				close(entries)
			}()
			for countEntry := fw.CountEntry(); countEntry > 0; countEntry-- {
				e, err := fw.ReadEntry()
				if err != nil {
					logger.Error("Error read entry", zap.Error(err))
					return
				}
				entries <- e
			}
		} else {
			//Handle if agent == nil
		}
	}

}

var Default = New()
