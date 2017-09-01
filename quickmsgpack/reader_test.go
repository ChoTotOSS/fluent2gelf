package quickmsgpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/family"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/format"
)

func TestReader(t *testing.T) {
	r := bytes.NewReader([]byte{0xce, 0x61, 0x62, 0x63, 0x22})

	reader := NewReader(r)

	t.Run("Test Reader can read a byte", func(t *testing.T) {
		b := reader.NextByte()
		if b != 0xce {
			t.Log("Reader read a file is incorrect")
			t.Fail()
		}
	})

	t.Run("Test Reader can read 3 bytes in a momment", func(t *testing.T) {
		b := reader.NextBytes(3)
		if bytes.Compare(b, []byte("abc")) != 0 {
			t.Log("Reader should return correct n bytes")
			t.Fail()
		}
	})
}

func TestReaderReadQuickMsp(t *testing.T) {
	r := bytes.NewReader([]byte{format.FixarrayLow + 5, 1, 2, 3, 4, 5})
	reader := NewReader(r)

	t.Run("Reader read array[5]", func(t *testing.T) {
		f, b := reader.NextFormat()

		if f != family.Array {
			t.Fail()
		}

		if !IsFixedFormat(b) {
			t.Logf("StringOf(b): %s", format.StringOf(b))
			t.Fail()
		}

		for i := 0; i < 5; i++ {
			if ValueOrExtraOf(reader.NextByte()) != i+1 {
				t.Fail()
			}
		}
	})
}

func TestReadNextLength(t *testing.T) {
	buf := make([]byte, 4)
	bufWriter := bytes.NewBuffer(buf)
	readerOf := func(i interface{}) *Reader {
		bufWriter.Reset()
		_ = binary.Write(bufWriter, binary.BigEndian, i)
		return NewReader(bytes.NewReader(buf))
	}

	t.Run("Read length by 1 byte", func(t *testing.T) {
		traver(0, 255, func(x uint8) {
			result := readerOf(x).NextLength(1)
			if uint(x) != result {
				t.Logf("result = %v, expected = %v\n", result, x)
				t.Fail()
			}
		})
	})

	t.Run("Read length up to 2 byte", func(t *testing.T) {
		for i := 0; i < 0xffff; i++ {
			x := uint16(i)
			result := readerOf(x).NextLength(2)
			if result != uint(x) {
				t.Logf("result = %v, expected = %v\n", result, x)
				t.Logf("buf = %#x\n", buf)
				t.Fail()
			}
		}
	})
	t.Run("Read length up to 4 byte", func(t *testing.T) {
		var i int64
		for i = 0; i < int64(0xffff); i++ {
			x := uint32(i)
			result := readerOf(x).NextLength(4)
			if result != uint(x) {
				t.Logf("result = %v, expected = %v\n", result, x)
				t.Logf("buf = %#x\n", buf)
				t.Fail()
			}
		}
	})

	t.Run("Panic", func(t *testing.T) {

		test := func(size uint8) {
			t.Run(fmt.Sprintf("Test panic for size: %v", size), func(t *testing.T) {
				defer func() {
					if x := recover(); x == nil {
						t.Log("Should be fail because of empty buf")
						t.Fail()
					}
				}()

				reader := NewReader(bytes.NewReader([]byte{}))
				reader.NextLength(size)
			})
		}

		test(1)
		test(2)
		test(4)
		test(6)
	})
}

func TestNextLengthOf(t *testing.T) {
	buf := make([]byte, 4)
	bufWriter := bytes.NewBuffer(buf)
	readerOf := func(i interface{}) *Reader {
		bufWriter.Reset()
		_ = binary.Write(bufWriter, binary.BigEndian, i)
		return NewReader(bytes.NewReader(buf))
	}

	t.Run("Test for fixarray", func(t *testing.T) {
		traver(format.FixarrayLow, format.FixarrayHigh, func(x uint8) {
			result := readerOf(x).NextLengthOf(x)
			if result != uint(FixedValueOf(x)) {
				t.Logf("result = %v, expected = %v\n", result, x)
				t.Fail()
			}
		})
	})
	test := func(extra uint8, formats ...byte) {
		for x := int64(math.Pow(8, float64(extra)/2)); x < int64(math.Pow(8, float64(extra))); x++ {
			for _, f := range formats {
				t.Run("Test read lenght of format", func(t *testing.T) {
					var r *Reader
					switch extra {
					case 1:
						r = readerOf(uint8(x))
					case 2:
						r = readerOf(uint16(x))
					case 4:
						r = readerOf(uint32(x))
					}

					result := r.NextLengthOf(f)
					if result != uint(x) {
						t.Logf("result = %v, expected = %v\n", result, x)
						t.Fail()
					}
				})
			}
		}
	}

	test(1, format.Int8, format.Bin8, format.Str8)
	test(2, format.Int16, format.Bin16, format.Str16)
	test(2, format.Array16, format.Map16, format.Ext16)
	test(4, format.Int32, format.Bin32, format.Str32)
	test(4, format.Array32, format.Map32, format.Ext32)
}
