package otg

import (
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/trigo"
	"github.com/kpmy/ypk/fn"
	"github.com/kpmy/ypk/halt"
	"io"
	"log"
	"math"
	"reflect"
	"strconv"
)

type Formatter interface {
	Write(otm.Object)
}

type Writer interface {
	RawString(string)
	Char(rune)
	Ident(string)
}

type wr struct {
	Writer
	base io.Writer
}

func (w *wr) Char(r rune) {
	w.RawString(string([]rune{r}))
}

func (w *wr) Ident(s string) {
	w.RawString(s)
}

func (w *wr) RawString(s string) {
	w.base.Write([]byte(s))
}

func (w *wr) Trinary(t tri.Trit) {
	switch t {
	case tri.NIL:
		w.RawString("null")
	case tri.TRUE:
		w.RawString("true")
	case tri.FALSE:
		w.RawString("false")
	}
}

type fm struct {
	wr
}

func (f *fm) stringValue(s string) {
	f.Char('"')
	f.RawString(s)
	f.Char('"')
}

func (f *fm) object(o otm.Object) {
	q := o.Qualident()
	if q.Template != "" {
		f.Ident(q.Template)
		f.Char('~')
	}
	f.Ident(q.Class)
	if q.Identifier != "" {
		f.Char('(')
		f.Ident(q.Identifier)
		f.Char(')')
	}
	if o.ChildrenCount() > 0 {
		f.Char(':')
		for _x := range o.Children() {
			f.Char(' ')
			log.Println(_x)
			switch x := _x.(type) {
			case otm.Object:
				f.object(x)
			case string:
				f.stringValue(x)
			case int64:
				i := strconv.Itoa(int(x))
				f.RawString(i)
			case float64:
				if math.IsInf(x, 1) {
					f.RawString("inf")
				} else if math.IsInf(x, -1) {
					f.RawString("-inf")
				} else {
					f_ := strconv.FormatFloat(x, 'f', 8, 64)
					f.RawString(f_)
				}
			case rune:
				r := strconv.FormatUint(uint64(x), 16)
				f.RawString("0" + r + "U")
			case tri.Trit:
				f.Trinary(x)
			default:
				halt.As(100, reflect.TypeOf(x))
			}
		}
		f.Char(';')
	} else {
		f.Char(' ')
	}
}

func (f *fm) Write(o otm.Object) {
	if !fn.IsNil(o) {
		f.object(o)
	} else {
		f.Trinary(tri.NIL)
	}
}

func ConnectTo(w io.Writer) Formatter {
	ret := &fm{}
	ret.wr.base = w
	return ret
}
