package path

import (
	"github.com/kpmy/ot/otm"
)

type ResultType int

const (
	OBJECT ResultType = iota
)

type Level interface {
	Qualident() otm.Qualident
}

type Traverser interface {
	Path() []Level
	Run(otm.Object, ResultType) (interface{}, error)
}
