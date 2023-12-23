package resp

type ArrayType struct {
	Value any
	Type  any
}

// CRLF constants
var (
	// \r
	CR byte = '\r'
	// \n
	LF byte = '\n'
	// \r\n
	breakPoint []byte = []byte{CR, LF}
)

// resp types
const (
	SIMPLE_STRING   string = "SIMPLE_STRING"
	BULK_STRING     string = "BULK_STRING"
	INTEGER         string = "INTEGER"
	ARRAY           string = "ARRAY"
	ERROR           string = "ERROR"
	UNSUPORTED_TYPE string = "UNSUPORTED_TYPE"
)

// resp prefixes
const (
	SIMPLE_STRING_PREFIX byte = '+'
	BULK_STRING_PREFIX   byte = '$'
	INTEGER_PREFIX       byte = ':'
	ARRAY_PREFIX         byte = '*'
	ERROR_PREFIX         byte = '-'
)
