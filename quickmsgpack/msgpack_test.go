package quickmsgpack

import (
	"fmt"
	"math"
	"testing"

	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/family"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/format"
)

func TestFamilyOf(t *testing.T) {

	test := func(begin byte, end byte, f uint16, name string) {
		t.Run("Test lookup "+name, func(t *testing.T) {
			name := format.StringOf(begin)
			t.Logf("Test %s[%#2x:%#2x]\n", name, begin, end)
			traver(begin, end, func(x byte) {
				if FamilyOf(x) != f {
					t.Logf("x = %#2x, family = %v, f = %v\n", x, FamilyOf(x), f)
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

func TestIsFixedFormat(t *testing.T) {
	test := func(expected bool, formats ...byte) {
		for _, f := range formats {
			t.Run(fmt.Sprintf("Test %s is fixed or not\n", format.StringOf(f)), func(t *testing.T) {
				if IsFixedFormat(f) != expected {
					t.Logf("Result = %v, expected = %v", IsFixedFormat(f), expected)
					t.Fatal()
				}
			})
		}
	}

	test(false, format.Array16, format.Array32,
		format.Bin8, format.Bin16, format.Bin32,
		format.Ext16, format.Ext32, format.Ext8,
		format.False, format.True, format.Float,
		format.Float64, format.Int8, format.Int16,
		format.Int32, format.Int64, format.Uint8,
		format.Uint16, format.Uint32, format.Uint64,
		format.Str8, format.Str16, format.Str32,
	)

	traver(format.PositiveFixintLow, format.PositiveFixintHigh, func(x byte) {
		test(true, x)
	})
	traver(format.NegativeFixintLow, format.NegativeFixintHigh, func(x byte) {
		test(true, x)
	})
	traver(format.FixarrayLow, format.FixarrayHigh, func(x byte) {
		test(true, x)
	})

	traver(format.FixmapLow, format.FixmapHigh, func(x byte) {
		test(true, x)
	})

	traver(format.FixstrLow, format.FixmapHigh, func(x byte) {
		test(true, x)
	})

	traver(format.Fixext1, format.Fixext16, func(x byte) {
		test(true, x)
	})
}

func TestExtraOrValue(t *testing.T) {
	test := func(expected int, formats ...byte) {
		for _, f := range formats {
			t.Run(fmt.Sprintf("Test %s extra for value", format.StringOf(f)), func(t *testing.T) {
				if ValueOrExtraOf(f) != expected {
					t.Logf("Result = %v, expected = %v", ValueOrExtraOf(f), expected)
					t.Fail()
				}
			})
		}
	}

	test(1, format.Bin8, format.Ext8, format.Int8, format.Uint8, format.Str8)
	test(2, format.Bin16, format.Ext16, format.Int16, format.Uint16, format.Str16)
	test(2, format.Array16, format.Map16)
	test(4, format.Bin32, format.Ext32, format.Int32, format.Uint32, format.Str32)
	test(4, format.Float, format.Array32, format.Map32)
	test(8, format.Float64, format.Int64)

	traver(format.PositiveFixintLow, format.PositiveFixintHigh, func(x byte) {
		test(int(x-format.PositiveFixintLow), x)
	})
	traver(format.NegativeFixintLow, format.NegativeFixintHigh, func(x byte) {
		test(int(x)-256, x)
	})
	traver(format.FixarrayLow, format.FixarrayHigh, func(x byte) {
		test(int(x-format.FixarrayLow), x)
	})

	traver(format.FixmapLow, format.FixmapHigh, func(x byte) {
		test(int(x-format.FixmapLow), x)
	})

	traver(format.FixstrLow, format.FixmapHigh, func(x byte) {
		test(int(x-format.FixstrLow), x)
	})

	traver(format.Fixext1, format.Fixext16, func(x byte) {
		test(int(math.Pow(2.0, float64(x-format.Fixext1))), x)
	})
}

func TestFixedValueOf(t *testing.T) {
	test := func(expected int8, formats ...byte) {
		for _, f := range formats {
			t.Run(fmt.Sprintf("Test %s extra for value", format.StringOf(f)), func(t *testing.T) {
				if FixedValueOf(f) != expected {
					t.Logf("Result = %v, expected = %v", FixedValueOf(f), expected)
					t.Fail()
				}
			})
		}
	}

	traver(format.PositiveFixintLow, format.PositiveFixintHigh, func(x byte) {
		test(int8(x-format.PositiveFixintLow), x)
	})
	traver(format.NegativeFixintLow, format.NegativeFixintHigh, func(x byte) {
		test(int8(x), x)
	})
	traver(format.FixarrayLow, format.FixarrayHigh, func(x byte) {
		test(int8(x-format.FixarrayLow), x)
	})

	traver(format.FixmapLow, format.FixmapHigh, func(x byte) {
		test(int8(x-format.FixmapLow), x)
	})

	traver(format.FixstrLow, format.FixmapHigh, func(x byte) {
		test(int8(x-format.FixstrLow), x)
	})

	traver(format.Fixext1, format.Fixext16, func(x byte) {
		test(int8(math.Pow(2.0, float64(x-format.Fixext1))), x)
	})
}

func TestExtraOf(t *testing.T) {
	test := func(expected uint8, formats ...byte) {
		for _, f := range formats {
			t.Run(fmt.Sprintf("Test %s extra for value", format.StringOf(f)), func(t *testing.T) {
				if ExtraOf(f) != expected {
					t.Logf("Result = %v, expected = %v", ExtraOf(f), expected)
					t.Fail()
				}
			})
		}
	}

	test(1, format.Bin8, format.Ext8, format.Int8, format.Uint8, format.Str8)
	test(2, format.Bin16, format.Ext16, format.Int16, format.Uint16, format.Str16)
	test(2, format.Array16, format.Map16)
	test(4, format.Bin32, format.Ext32, format.Int32, format.Uint32, format.Str32)
	test(4, format.Float, format.Array32, format.Map32)
	test(8, format.Float64, format.Int64)
}
