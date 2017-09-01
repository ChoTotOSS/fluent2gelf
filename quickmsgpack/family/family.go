package family

const (
	Integer uint16 = 1 << iota
	Nil
	Boolean
	Float
	String
	Binary
	Array
	Map
	Extension
	Unknown
)
