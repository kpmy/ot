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

type Class interface {
	Qualident() Qualident
}

type Object interface {
	Qualident() Qualident
	InstanceOf(...Class) Class

	Parent() Object

	Children() chan interface{}
	ChildrenObjects() chan Object
	ChildrenCount() uint
}

type Link interface {
	Object() Object
}

func RootOf(o Object) Object {
	if fn.IsNil(o.Parent()) {
		return o
	} else {
		return RootOf(o.Parent())
	}
}
