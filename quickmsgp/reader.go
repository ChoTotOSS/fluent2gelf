package quickmsgpack

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r}
}

// NextByte Read a byte of data from io
func (r *Reader) NextByte() byte {
	b := make([]byte, 1)
	_, _ = r.Read(b) //Read 1 byte, skip error
	return b[0]
}

//NextBytes read n bytes of data from io
func (r *Reader) NextBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = r.Read(b)
	return b
}

func (r *Reader) NextFormat() (uint16, byte) {
	b := r.NextByte()
	return familyOf[b], b
}

func (r *Reader) NextInt8() int8 {
	var i int8
	err := binary.Read(r, binary.BigEndian, &i)
	if err != nil {
		panic(err)
	}
	return i
}
