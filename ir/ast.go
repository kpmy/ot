package ir

import (
	"github.com/kpmy/ot/ir/types"
)

type Template struct {
	Stmt []Statement
}

type Statement interface {
	Process() Statement
}

type Emit struct {
	Template   string
	Class      string
	Ident      string
	ChildCount uint
}

func (i *Emit) Process() Statement { return i }

type Dive struct {
	Reuse bool
}

func (i *Dive) Process() Statement { return i }

type Rise struct{}

func (i *Rise) Process() Statement { return i }

type Put struct {
	Type  types.Type
	Value interface{}
}

func (i *Put) Process() Statement { return i }
