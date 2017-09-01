package format

const (
	//Fix int 0:127
	PositiveFixintLow  byte = 0x00 // 0
	PositiveFixintHigh      = 0x7f // 127
	NegativeFixintLow       = 0xe0
	NegativeFixintHigh      = 0xff

	//What is nil?
	Nil = 0xc0

	Never = 0xc1
	// Bool
	False = 0xc2
	True  = 0xc3

	// Float
	Float   = 0xca
	Float64 = 0xcb

	// Unsigned Int
	Uint8  = 0xcc
	Uint16 = 0xcd
	Uint32 = 0xce
	Uint64 = 0xcf

	// Int
	Int8  = 0xd0
	Int16 = 0xd1
	Int32 = 0xd2
	Int64 = 0xd3

	//String
	FixstrLow  = 0xa0
	FixstrHigh = 0xbf
	Str8       = 0xd9
	Str16      = 0xda
	Str32      = 0xdb

	// bin
	Bin8  = 0xc4
	Bin16 = 0xc5
	Bin32 = 0xc6

	//Array
	FixarrayLow  = 0x90
	FixarrayHigh = 0x9f
	Array16      = 0xdc
	Array32      = 0xdd

	FixmapLow  = 0x80
	FixmapHigh = 0x8f
	Map16      = 0xde
	Map32      = 0xdf

	Fixext1  = 0xd4
	Fixext2  = 0xd5
	Fixext4  = 0xd6
	Fixext8  = 0xd7
	Fixext16 = 0xd8
	Ext8     = 0xc7
	Ext16    = 0xc8
	Ext32    = 0xc9
)
