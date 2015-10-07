package otm

import (
	"fmt"
	"github.com/kpmy/ypk/fn"
)

type Qualident struct {
	Template   string
	Class      string
	Identifier string
}

func (q Qualident) String() string {
	return fmt.Sprint("", fn.MaybeString(q.Template, "."), "", fn.MaybeString(q.Class), "", fn.MaybeString("(", q.Identifier, ")"))
}

type Object interface {
	Parent() Object
	Qualident() Qualident
	Children() chan interface{}
	ChildrenCount() uint
}

type Link interface {
	Object() Object
}
