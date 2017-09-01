package gelf

import (
	"bytes"
	"encoding/binary"
	"time"
)

func RandomChunkId() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, time.Now().UnixNano())
	return buf.Bytes()
}

func GetByteValue(i int) byte {
	return byte(i)
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

	if dataLenght%SIXTY_FOUR_KiB > 0 {
		seqCount = dataLenght/SIXTY_FOUR_KiB + 1
	} else {
		seqCount = dataLenght / SIXTY_FOUR_KiB
	}

	return &Chunks{
		seqCount: seqCount,
		seqIndex: 0,
		msgId:    RandomChunkId(),
		data:     data,
	}
}

var MAGIC_BYTES = []byte{0x1e, 0x0f}

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

	var buffer bytes.Buffer
	buffer.Write(MAGIC_BYTES)
	buffer.Write(c.msgId)
	buffer.WriteByte(byte(c.seqIndex))
	buffer.WriteByte(byte(c.seqCount))

	left := c.seqIndex * SIXTY_FOUR_KiB
	right := (c.seqIndex + 1) * SIXTY_FOUR_KiB
	if right > len(c.data) {
		right = len(c.data)
	}

	buffer.Write(c.data[left:right])

	return buffer.Bytes()
}
