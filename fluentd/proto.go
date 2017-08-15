package fluentd

import (
	"io"
	"net"
	"reflect"

	"github.com/ChoTotOSS/fluent2gelf/agent"
	"github.com/duythinht/zaptor"
	"github.com/ugorji/go/codec"
	"go.uber.org/zap"
)

var logger = zaptor.Default()
var mh codec.MsgpackHandle

func init() {
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

func ForwardHandle(conn net.Conn, agentStore *agent.AgentStore) {

	decoder := codec.NewDecoder(conn, &mh)

	v := []interface{}{nil, nil, nil}

	err := decoder.Decode(&v)
	if err != nil {
		if err != io.EOF {
			logger.Error("decode error", zap.Error(err))
		}
		return
	}

	tag, ok := v[0].([]byte)

	if ok {
		logger.Debug("entries", zap.ByteString("tag", tag))

		a := agentStore.Take(string(tag))

		if a == nil {
			return
		}

		entries, ok := v[1].([]interface{})

		if ok {
			logger.Debug("Deserialize entries for", zap.ByteString("tag", tag))
			a.Put(entries)
		} else {
			entries, ok := v[1].([]uint8)
			if ok {
				logger.Warn("Need implements unzip for entries", zap.Any("entries", entries))
			} else {
				logger.Warn("decode entries failed",
					zap.ByteString("tag", tag),
					zap.Any("entries", v[1]),
					zap.Any("options", v[2]),
				)
			}
		}
	}
}
