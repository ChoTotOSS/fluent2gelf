package fluentd

import (
	"bytes"
	"net"

	"go.uber.org/zap"

	"github.com/ChoTotOSS/fluent2gelf/agent"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/family"
)

func ForwardHandle(conn net.Conn, agentStore *agent.AgentStore) {
	reader := quickmsgpack.NewReader(conn)
	defer func() { _ = conn.Close() }()

	f, b := reader.NextFormat() //type should be fixarray

	if f != family.Array || !quickmsgpack.IsFixedFormat(b) {
		return
	}

	switch quickmsgpack.FixedValueOf(b) {
	case 2:
		handleWithoutZip(reader, agentStore)
	case 3:

	default:
	}
}

func handleWithoutZip(r *quickmsgpack.Reader, agentStore *agent.AgentStore) {
	f, b := r.NextFormat()
	if f != family.String {
		return
	}

	var tag string

	if quickmsgpack.IsFixedFormat(b) {
		//logger.Debug("Fixed string tag")
		tag = string(r.NextBytes(uint(quickmsgpack.FixedValueOf(b))))
	} else {
		//logger.Debug("String tag")
		extra := quickmsgpack.ExtraOf(b)
		tag = string(r.NextBytes(r.NextLength(extra)))
	}

	//logger.Debug("Handle message", zap.String("tag", tag))

	a := agentStore.Take(tag)

	if a == nil {
		return
	}

	//Handle for entries block
	f, b = r.NextFormat()

	if f != family.Array {
		logger.Warn("msgp entries wrong format", zap.Any("family", f))
		return
	}

	count := r.NextLengthOf(b)
	c := make(chan []byte)
	a.Ship(c)
	for ; count > 0; count-- {
		if e := nextEntry(r); e != nil {
			c <- e.Log
		}
	}
	close(c)
}

type Entry struct {
	Log []byte
}

func nextEntry(r *quickmsgpack.Reader) *Entry {
	f, b := r.NextFormat()
	if f != family.Array || quickmsgpack.FixedValueOf(b) != 2 {
		logger.Warn("Entry is wrong format", zap.Uint16("family", f), zap.Any("fixedValue", quickmsgpack.FixedValueOf(b)))
		return nil
	}

	f, b = r.NextFormat()

	switch f {
	case family.Integer:
		r.NextBytes(4)
	case family.Extension:
		logger.Warn("Does not implements timestamp ext")
		return nil
	default:
		logger.Warn("Wrong time format entry")
		return nil
	}

	// then read record
	f, b = r.NextFormat()

	if f != family.Map {
		logger.Warn("Wrong record format")
		return nil
	}

	count := r.NextLengthOf(b)

	var e *Entry

	for ; count > 0; count-- {
		// Read for key
		f, b = r.NextFormat()
		if f != family.String {
			logger.Warn("Wrong record format")
			return nil
		}

		key := r.NextBytes(r.NextLengthOf(b))

		// Read for value
		f, b = r.NextFormat()
		if f != family.String {
			logger.Warn("Wrong record format")
			return nil
		}

		value := r.NextBytes(r.NextLengthOf(b))

		//logger.Debug("record map", zap.ByteString("key", key), zap.ByteString("value", value))

		if bytes.Compare(key, []byte("log")) == 0 {
			e = new(Entry)
			e.Log = value
		}
	}
	return e
}
