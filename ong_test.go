package ot

import (
	zfmt "fmt"
	"github.com/kpmy/rng/fn"
	"github.com/kpmy/rng/schema"
	"github.com/kpmy/rng/schema/std"
	"github.com/kpmy/ypk/halt"
	"reflect"
	"testing"
)

var level int
var passed map[string]schema.Guide = make(map[string]schema.Guide)

func tab() (ret string) {
	for i := 0; i < level; i++ {
		ret = zfmt.Sprint(ret, " ")
	}
	return
}

func verbose(_g interface{}, meta ...interface{}) (ret interface{}) {
	level++
	tfmt.Println(meta...)
	delim := " "
	switch g := _g.(type) {
	case schema.Start:
		tfmt.Println(tab(), "$start", g)
	case schema.Choice:
		tfmt.Println(tab(), "choice")
	case schema.Element:
		tfmt.Println(tab(), "$element", g.Name())
	case schema.Attribute:
		tfmt.Println(tab(), "$attribute", g.Name())
	case schema.Interleave:
		tfmt.Println(tab(), "interleave")
	case schema.ZeroOrMore:
		tfmt.Println(tab(), "zero-or-more")
	case schema.OneOrMore:
		tfmt.Println(tab(), "one-or-more")
	case schema.Optional:
		tfmt.Println(tab(), "optional")
	case schema.Group:
		tfmt.Println(tab(), "group")
	case schema.AnyName:
		tfmt.Println(tab(), "any-name")
	case schema.Except:
		tfmt.Println(tab(), "except")
	case schema.NSName:
		tfmt.Println(tab(), "ns-name", g.NS())
	case schema.Text:
		tfmt.Println(tab(), "text")
	case schema.Data:
		tfmt.Println(tab(), "data")
		if g.Type() != "" {
			tfmt.Println(" ", "type ", g.Type())
		}
		tfmt.Println()
	case schema.Value:
		tfmt.Println(tab(), "value", g.Data())
	case schema.Name:
		tfmt.Println(tab(), "name")
	case schema.Empty:
		tfmt.Println(tab(), "empty")
	case schema.List:
		tfmt.Println(tab(), "list")
	case schema.Mixed:
		tfmt.Println(tab(), "mixed")
	case schema.Param:
		tfmt.Println(tab(), "param")
	case schema.Ref:
		tfmt.Println(tab(), "ref", g.Name())
	case schema.ExternalRef:
		tfmt.Println(tab(), "externalRef", g.Href())
	default:
		halt.As(100, reflect.TypeOf(g))
	}
	if id := _g.(std.Identified).Id(); passed[id] == nil {
		passed[id] = _g.(schema.Guide)
		fn.Map(fn.Iterate(_g.(schema.Guide)), verbose, delim)
	}
	level--
	return _g
}

func print(_g interface{}, meta ...interface{}) interface{} {
	tfmt.Println(_g)
	return _g
}

func elementFilter(name string) fn.Bool {
	return func(g schema.Guide, _ ...interface{}) bool {
		e, ok := g.(schema.Element)
		return ok && e.Name() == name
	}
}

type formatter struct {
	printfn func(...interface{})
}

func (f *formatter) Println(x ...interface{}) {
	if f.printfn != nil {

	}
}

var tfmt *formatter

func init() {
	tfmt = &formatter{}
}

func testSchemaPrint(start schema.Start, log *testing.T) {
	tfmt.printfn = func(x ...interface{}) {
		log.Log(x...)
	}

	verbose(start)
	log.Log("---")
	fn.Traverse(start, print)
	log.Log("---")
}
