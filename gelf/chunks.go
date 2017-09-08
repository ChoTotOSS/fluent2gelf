package gelf

import (
	"bytes"
	"encoding/binary"
	"time"
)

var MAGIC_BYTES = []byte{0x1e, 0x0f}

func randomChunkId() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, time.Now().UnixNano())
	return buf.Bytes()
}

type Chunks struct {
	seqCount int
	seqIndex int
	msgId    []byte
	data     []byte
}

func NewChunks(data []byte) *Chunks {
	dataLenght := len(data)
	var seqCount int

	// Use this bcz it's faster than math.Ceil
	if dataLenght%MAX_CHUNK_SIZE > 0 {
		seqCount = dataLenght/MAX_CHUNK_SIZE + 1
	} else {
		seqCount = dataLenght / MAX_CHUNK_SIZE
	}

	return &Chunks{
		seqCount: seqCount,
		seqIndex: 0,
		msgId:    randomChunkId(),
		data:     data,
	}
}

func (c *Chunks) HasNext() bool {
	return c.seqIndex < c.seqCount
}

func (c *Chunks) Next() []byte {
	defer func() {
		c.seqIndex++
	}()

	//logger.Debug("chunk", zap.Int("seqc", c.seqCount), zap.Int("seqindex", c.seqIndex))

	//quick response
	if c.seqCount == 1 {
		return c.data
	}

	left := c.seqIndex * MAX_CHUNK_SIZE
	right := (c.seqIndex + 1) * MAX_CHUNK_SIZE

	if right > len(c.data) {
		right = len(c.data)
	}

	data := c.data[left:right]

	var buffer bytes.Buffer

	// Write 12 bytes of header
	buffer.Write(MAGIC_BYTES)          // 2 bytes
	buffer.Write(c.msgId)              // 8 bytes for message id
	buffer.WriteByte(byte(c.seqIndex)) // 1 byte to seq index
	buffer.WriteByte(byte(c.seqCount)) // 1 byte for seq count
	// Then write data
	buffer.Write(data)

	return buffer.Bytes()
}
