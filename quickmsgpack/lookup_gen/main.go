package main

import (
	"bytes"
	"fmt"
	gofmt "go/format"
	"io"
	"io/ioutil"
	"log"

	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/family"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/format"
)

const (
	bitmask = 0x0f
)

var header = []byte(`package quickmsgpack

const (
	yes = 1 == 1
	nop = 0 == 1
)
`)

func main() {
	buf := bytes.NewBuffer(header)
	gen_family(buf, "familyOf")
	gen_extra(buf, "extraBytesOf")
	gen_fixedcheck(buf, "isFixed")
	gen_fixedvalues(buf, "fixedValueOf")

	fmt.Printf("%s\n", buf.Bytes())

	out, err := gofmt.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("format_lookup.go", out, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func gen_family(w io.Writer, name string) {
	fmt.Fprintf(w, "var %s = [256]uint16{", name)
	for i := 0; i < 256; i++ {
		if i%8 == 0 {
			fmt.Fprint(w, "\n  ")
		} else {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, "%#04x,", familyOf(uint8(i)))
	}
	fmt.Fprint(w, "\n}\n\n")
}

func gen_extra(w io.Writer, name string) {
	fmt.Fprintf(w, "var %s = [256]uint8{", name)
	for i := 0; i < 256; i++ {
		if i%16 == 0 {
			fmt.Fprint(w, "\n  ")
		} else {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, "%#02x,", extraBytesOf(uint8(i)))
	}
	fmt.Fprint(w, "\n}\n\n")
}

func gen_fixedcheck(w io.Writer, name string) {
	fmt.Fprintf(w, "var %s = [256]bool{", name)
	for i := 0; i < 256; i++ {
		if i%16 == 0 {
			fmt.Fprint(w, "\n  ")
		} else {
			fmt.Fprint(w, " ")
		}
		if isFixed(uint8(i)) {
			fmt.Fprintf(w, "yes,")
		} else {
			fmt.Fprintf(w, "nop,")
		}
	}
	fmt.Fprint(w, "\n}\n\n")
}

func gen_fixedvalues(w io.Writer, name string) {
	fmt.Fprintf(w, "var %s = [256]int8{", name)
	for i := 0; i < 256; i++ {
		if i%16 == 0 {
			fmt.Fprint(w, "\n  ")
		} else {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, "%d,", fixedValueOf(uint8(i)))
	}
	fmt.Fprint(w, "\n}\n\n")
}

func familyOf(b byte) uint16 {
	switch b {
	case format.Int8, format.Int16, format.Int32, format.Int64, format.Uint8, format.Uint16, format.Uint32, format.Uint64:
		return family.Integer
	case format.Array16, format.Array32:
		return family.Array
	case format.Map16, format.Map32:
		return family.Map
	case format.Float, format.Float64:
		return family.Float
	case format.True, format.False:
		return family.Boolean
	case format.Str8, format.Str16, format.Str32:
		return family.String
	case format.Bin8, format.Bin16, format.Bin32:
		return family.Binary
	case format.Fixext1, format.Fixext2, format.Fixext4, format.Fixext8, format.Fixext16, format.Ext8, format.Ext16, format.Ext32:
		return family.Extension
	default: //Check for fixtype
		switch {
		case b <= format.PositiveFixintHigh || b >= format.NegativeFixintLow:
			return family.Integer
		case b >= format.FixarrayLow && b <= format.FixarrayHigh:
			return family.Array
		case b >= format.FixmapLow && b <= format.FixmapHigh:
			return family.Map
		case b >= format.FixstrLow && b <= format.FixstrHigh:
			return family.String
		}
	}
	return family.Unknown
}

func extraBytesOf(b byte) uint8 {
	switch b {
	case format.Uint8, format.Int8, format.Str8, format.Bin8, format.Ext8:
		return 1
	case format.Uint16, format.Int16, format.Str16, format.Bin16, format.Array16, format.Map16, format.Ext16:
		return 2
	case format.Uint32, format.Int32, format.Str32, format.Bin32, format.Array32, format.Map32, format.Ext32:
		return 4
	case format.Float:
		return 4
	case format.Uint64, format.Int64, format.Float64:
		return 8
	}
	//Fixed data does not need extra bytes
	return 0
}

func isFixed(b byte) bool {
	switch b {
	case format.Nil, format.True, format.False, format.Never:
		return false
	default:
		return extraBytesOf(b) == 0
	}
	return false
}

func fixedValueOf(b byte) int8 {
	if !isFixed(b) {
		return 0
	}
	switch familyOf(b) {
	case family.Integer:
		return int8(b)
	case family.String:
		return int8(b & mask(format.FixstrLow, format.FixstrHigh))
	case family.Array:
		return int8(b & mask(format.FixarrayLow, format.FixarrayHigh))
	case family.Map:
		return int8(b & mask(format.FixmapLow, format.FixmapHigh))
	case family.Extension:
		switch b {
		case format.Fixext1:
			return 1
		case format.Fixext2:
			return 2
		case format.Fixext4:
			return 4
		case format.Fixext8:
			return 8
		case format.Fixext16:
			return 16
		}
	}
	return 0
}

func mask(low byte, high byte) byte {
	return high - low
}
