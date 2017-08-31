package quickmsgpack

import (
	"bytes"
	"testing"

	"github.com/ChoTotOSS/fluent2gelf/quickmsgp/family"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgp/format"
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

func TestReadInt8(t *testing.T) {
	data := []uint8{}
	for i := 0; i < 256; i++ {
		data = append(data, uint8(i))
	}
	r := bytes.NewReader(data)
	reader := NewReader(r)
	t.Run("Correct", func(t *testing.T) {
		for i := 0; i < 256; i++ {
			if reader.NextInt8() != int8(i) {
				t.Fail()
			}
		}
	})
	t.Run("Panic", func(t *testing.T) {
		defer func() {
			if x := recover(); x != nil {
				return
			}
			t.Log("Read empty reader should return panic")
			t.Fail()
		}()
		reader.NextInt8()
	})
}
