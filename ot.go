package ot

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/kpmy/ot/ir"
	"github.com/kpmy/ot/ir/types"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/trigo"
	"github.com/kpmy/ypk/fn"
	"github.com/kpmy/ypk/halt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func prettyPrint(t *ir.Template) {
	wr := bytes.NewBufferString("")
	depth := 0
	tab := func() {
		for i := 0; i < depth; i++ {
			fmt.Fprint(wr, " ")
		}
	}
	tab()
	for _, _s := range t.Stmt {
		switch s := _s.(type) {
		case *ir.Emit:
			tab()
			fmt.Fprint(wr, fn.MaybeString(s.Template, "."), s.Class, fn.MaybeString("(", s.Ident, ")"))
			if s.ChildCount == 0 {
				fmt.Fprintln(wr)
			}
		case *ir.Dive:
			if !s.Reuse {
				fmt.Fprintln(wr, ":")
			} else {
				fmt.Fprintln(wr, " ::")
			}
			depth++
		case *ir.Rise:
			depth--
			tab()
			fmt.Fprintln(wr, ";")
		case *ir.Put:
			tab()
			switch s.Type {
			case types.STRING:
				fmt.Fprint(wr, "`", s.Value, "`")
			case types.REAL, types.INTEGER, types.TRILEAN:
				fmt.Fprint(wr, s.Value)
			case types.CHAR:
				fmt.Fprint(wr, "0", strings.ToUpper(strconv.FormatUint(uint64(s.Value.(rune)), 16)), "U")
			case types.LINK:
				fmt.Fprint(wr, "^", s.Value)
			default:
				halt.As(100, s.Type)
			}
			fmt.Fprintln(wr)
		default:
			halt.As(100, reflect.TypeOf(s))
		}
	}
	log.Print(wr.String())
}

func prettyPrintObject(o otm.Object) {
	parent := ""
	for x := o.Parent(); !fn.IsNil(x); x = x.Parent() {
		parent = fmt.Sprint(x.Qualident(), "<-", parent)
	}
	log.Println(parent, o.Qualident())
	if o.ChildrenCount() > 0 {
		log.Println(":")
		for _x := range o.Children() {
			switch x := _x.(type) {
			case otm.Object:
				prettyPrintObject(x)
			case otm.Link:
				log.Println("^", x.Object().Qualident())
			case string, float64, int64, rune, tri.Trit:
				log.Print(_x)
			default:
				halt.As(100, reflect.TypeOf(x))
			}
		}
		log.Println(";")
	}
}

func renderHtml(o otm.Object) {
	buf := bytes.NewBufferString("<!DOCTYPE HTML>")
	e := xml.NewEncoder(buf)
	var obj func(otm.Object)
	obj = func(o otm.Object) {
		clazz := o.InstanceOf().Qualident()
		if clazz.Template == "html" {
			start := xml.StartElement{}
			start.Name.Local = clazz.Class
			if id := o.Qualident().Identifier; id != "" {
				attr := xml.Attr{}
				attr.Name.Local = "id"
				attr.Value = id
				start.Attr = append(start.Attr, attr)
			}
			e.EncodeToken(start)
			for _x := range o.Children() {
				switch x := _x.(type) {
				case otm.Object:
					obj(x)
				case string:
					e.EncodeToken(xml.CharData([]byte(x)))
				default:
					halt.As(100, reflect.TypeOf(x))
				}
			}
			e.EncodeToken(start.End())
		}
	}

	for x := range o.ChildrenObjects() {
		if x.InstanceOf().Qualident().Template == "html" && x.InstanceOf().Qualident().Class == "html" {
			obj(x)
		}
	}
	e.Flush()
	log.Println(buf.String())
}
