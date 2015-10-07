package ot

import (
	"bytes"
	"fmt"
	"github.com/kpmy/ot/ir"
	"github.com/kpmy/ot/ir/types"
	"github.com/kpmy/ot/otm"
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
			case types.REAL, types.INTEGER:
				fmt.Fprint(wr, s.Value)
			case types.CHAR:
				fmt.Fprint(wr, "0", strings.ToUpper(strconv.FormatUint(uint64(s.Value.(rune)), 16)), "U")
			case types.LINK:
				fmt.Fprint(wr, "@", s.Value)
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
				log.Println("@", x.Object().Qualident())
			case string, float64, int64, rune:
				log.Print(_x)
			default:
				halt.As(100, reflect.TypeOf(x))
			}
		}
		log.Println(";")
	}
}
