package types

import (
	"strconv"
)

type Type uint

const (
	NONE Type = iota
	STRING
	REAL
	INTEGER
	CHAR

	LINK
)

func (t Type) String() string {
	switch t {
	case NONE:
		return "NONE"
	case STRING:
		return "STRING"
	case REAL:
		return "REAL"
	case INTEGER:
		return "INTEGER"
	case CHAR:
		return "CHAR"

	case LINK:
		return "LINK"
	default:
		return strconv.Itoa(int(t))
	}
}
