package quickmsgpack

import (
	"fmt"
	"testing"

	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/family"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/format"
)

func traver(begin byte, end byte, f func(byte)) {
	for i := int(begin); i < int(end); i++ {
		f(uint8(i))
	}
}

func TestFamilyLookup(t *testing.T) {

	test := func(begin byte, end byte, f uint16, name string) {
		t.Run("Test lookup "+name, func(t *testing.T) {
			t.Logf("Test in range %2x - %2x for %s\n", begin, end, name)
			traver(begin, end, func(x byte) {
				if familyOf[x] != f {
					t.Logf("x = %#2x, family = %v, f = %v\n", x, familyOf[x], f)
					t.Fatal()
				}
			})
		})
	}

	test(format.PositiveFixintLow, format.PositiveFixintHigh, family.Integer, "Postive fixed Integer")
	test(format.NegativeFixintLow, format.NegativeFixintHigh, family.Integer, "Negative fixed Integer")
	test(format.FixmapLow, format.FixmapHigh, family.Map, "Fixed Map")
	test(format.FixstrLow, format.FixstrHigh, family.String, "Fixed String")
	test(format.Fixext1, format.Fixext16, family.Extension, "Fixed Extension")
	test(format.Nil, format.Nil, family.Nil, "Nil")
	test(format.Never, format.Never, family.Unknown, "Never")
	test(format.False, format.True, family.Boolean, "Boolean")
	test(format.Bin8, format.Bin32, family.Binary, "Binary")
	test(format.Array16, format.Array32, family.Array, "Array")
	test(format.Map16, format.Map32, family.Map, "Map")
	test(format.Float, format.Float64, family.Float, "Float")
	test(format.Uint8, format.Int64, family.Integer, "Integer")
	test(format.Str8, format.Str32, family.String, "String")
	test(format.Ext8, format.Ext32, family.Extension, "Extension")
}

func TestExtraBytesLookupFixed(t *testing.T) {
	test := func(begin byte, end byte, expected uint8, name string) {
		t.Run("Test extra lookup "+name, func(t *testing.T) {
			t.Logf("Test in range %2x - %2x for %s\n", begin, end, name)
			traver(begin, end, func(x byte) {
				if extraBytesOf[x] != expected {
					t.Logf("x = %#2x, got = %v, expected = %v\n", x, extraBytesOf[x], expected)
					t.Fatal()
				}
			})
		})
	}

	test(format.PositiveFixintLow, format.PositiveFixintHigh, 0, "Postive fixed Integer")
	test(format.NegativeFixintLow, format.NegativeFixintHigh, 0, "Negative fixed Integer")
	test(format.FixmapLow, format.FixmapHigh, 0, "Fixed Map")
	test(format.FixstrLow, format.FixstrHigh, 0, "Fixed String")
	test(format.Fixext1, format.Fixext16, 0, "Fixed Extension")
	test(format.Nil, format.Nil, 0, "Nil")
	test(format.Never, format.Never, 0, "Never")
}

func TestExtraBytesLookup(t *testing.T) {
	test := func(expected uint8, formats ...byte) {
		for _, x := range formats {
			if extraBytesOf[x] != expected {
				t.Logf("x = %#2x, got = %v, expected = %v\n", x, extraBytesOf[x], expected)
				t.Fatal()
			}
		}
	}
	test(1, format.Int8, format.Bin8, format.Str8)
	test(2, format.Int16, format.Bin16, format.Str16)
	test(2, format.Array16, format.Map16, format.Ext16)
	test(4, format.Int32, format.Bin32, format.Str32)
	test(4, format.Array32, format.Map32, format.Ext32)
	test(8, format.Float64, format.Int64)
}

func TestFixedValues(t *testing.T) {
	test := func(low byte, high byte) {
		t.Run(fmt.Sprintf("Test_%s[%#2x:%#2x]", format.StringOf(low), low, high), func(t *testing.T) {
			traver(low, high, func(x byte) {
				if fixedValueOf[x] != int8(x-low) {
					t.Logf("%s, result = %v, expected = %v\n", format.StringOf(x), fixedValueOf[x], x-low)
					t.Fatal()
				}
			})
		})
	}

	test(format.PositiveFixintLow, format.PositiveFixintHigh)
	test(format.FixarrayLow, format.FixarrayHigh)
	test(format.FixmapLow, format.FixmapHigh)
	test(format.FixstrLow, format.FixstrHigh)

	t.Run("Test Fixext", func(t *testing.T) {
		if fixedValueOf[format.Fixext1] != 1 {
			t.Fail()
		}
		if fixedValueOf[format.Fixext2] != 2 {
			t.Fail()
		}
		if fixedValueOf[format.Fixext4] != 4 {
			t.Fail()
		}
		if fixedValueOf[format.Fixext8] != 8 {
			t.Fail()
		}
		if fixedValueOf[format.Fixext16] != 16 {
			t.Fail()
		}
	})
}
