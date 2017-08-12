package gelf

import (
	"bytes"
	"compress/gzip"

	"go.uber.org/zap"
)

func ZipMessage(message []byte) []byte {

	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err := writer.Write(message)

	if err != nil {
		logger.Warn("zip#error", zap.Error(err))
	}

	_ = writer.Close()
	return buffer.Bytes()
}
