// Code generated by "stringer -type=byte"; DO NOT EDIT.

package format

var (
	nameTable = map[byte]string{
		Int8:     "Int8",
		Int16:    "Int16",
		Int32:    "Int32",
		Int64:    "Int64",
		Float:    "Float",
		Float64:  "Float64",
		Array16:  "Array16",
		Array32:  "Array32",
		Str8:     "Str8",
		Str16:    "Str16",
		Str32:    "Str32",
		Bin8:     "Bin8",
		Bin16:    "Bin16",
		Bin32:    "Bin32",
		Map16:    "Map16",
		Map32:    "Map32",
		Ext8:     "Ext8",
		Ext16:    "Ext16",
		Ext32:    "Ext32",
		Fixext1:  "Fixext1",
		Fixext2:  "Fixext2",
		Fixext4:  "Fixext4",
		Fixext8:  "Fixext8",
		Fixext16: "Fixext16",
	}
)

func StringOf(b byte) string {
	if s, ok := nameTable[b]; ok {
		return s
	}

	switch {
	case b >= FixmapLow && b <= FixmapHigh:
		return "Fixmap"
	case b >= FixarrayLow && b <= FixarrayHigh:
		return "Fixarray"
	case b >= FixstrLow && b <= FixstrHigh:
		return "Fixstr"
	case b <= PositiveFixintHigh:
		return "PositiveFixint"
	case b >= NegativeFixintLow:
		return "NegativeFixint"
	default:
		return "Never"
	}
}
