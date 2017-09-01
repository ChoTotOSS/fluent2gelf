package quickmsgpack

import (
	"encoding/binary"
	"errors"
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
func (r *Reader) NextBytes(n uint) []byte {
	b := make([]byte, n)
	_, _ = r.Read(b)
	return b
}

func (r *Reader) NextFormat() (uint16, byte) {
	b := r.NextByte()
	return familyOf[b], b
}

func (r *Reader) NextLength(size uint8) uint {
	switch size {
	case 1:
		var i uint8
		tryOrDie(binary.Read(r, binary.BigEndian, &i))
		return uint(i)
	case 2:
		var i uint16
		tryOrDie(binary.Read(r, binary.BigEndian, &i))
		return uint(i)
	case 4:
		var i uint32
		tryOrDie(binary.Read(r, binary.BigEndian, &i))
		return uint(i)
	default:
		panic(errors.New("Length size only support: 1, 2, 4."))
	}
}

//NextLengthOf return length of map, string, arry, bin and ext format
func (r *Reader) NextLengthOf(b byte) uint {
	if IsFixedFormat(b) {
		return uint(FixedValueOf(b))
	}
	return r.NextLength(ExtraOf(b))
}

func tryOrDie(err error) {
	if err != nil {
		panic(err)
	}
}
