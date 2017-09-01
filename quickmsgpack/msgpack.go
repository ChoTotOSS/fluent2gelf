package quickmsgpack

//go:generate go run lookup_gen/main.go

// FamilyOf return family format
func FamilyOf(b byte) uint16 {
	return familyOf[b]
}

func IsFixedFormat(b byte) bool {
	return isFixed[b]
}

func ValueOrExtraOf(b byte) int {
	if isFixed[b] {
		return int(fixedValueOf[b])
	}
	return int(extraBytesOf[b])
}

func FixedValueOf(b byte) int8 {
	return fixedValueOf[b]
}

func ExtraOf(b byte) uint8 {
	return extraBytesOf[b]
}
