package otm

import (
	"fmt"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/fn"
)

type CopyMode bool

const (
	DEEP    CopyMode = true
	SHALLOW CopyMode = false
)

type Qualident struct {
	Template   string
	Class      string
	Identifier string
}

func (q Qualident) String() string {
	return fmt.Sprint("", fn.MaybeString(q.Template, "~"), "", fn.MaybeString(q.Class), "", fn.MaybeString("(", q.Identifier, ")"))
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

	FindById(string) Object
	FindByQualident(Qualident) []Object
	CopyOf(deep CopyMode) Object
}

type Link interface {
	Object() Object
}

func RootOf(o Object) Object {
	assert.For(!fn.IsNil(o), 20)
	if fn.IsNil(o.Parent()) {
		return o
	} else {
		return RootOf(o.Parent())
	}
}

type Modifier func(Object) Object
type Producer func(...Modifier) Object

type Builder interface {
	End(...Modifier) Object
	Prod() Producer

	Value(...interface{}) Builder
	Child(Producer, ...Modifier) Builder
}
