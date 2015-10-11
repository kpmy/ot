package conv

import (
	"container/list"
	"fmt"
	"github.com/kpmy/ot/ir"
	"github.com/kpmy/ot/ir/types"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/trigo"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/fn"
	"github.com/kpmy/ypk/halt"
	"hash/adler32"
	"reflect"
	"strconv"
)

type class struct {
	tpl, cls string
}

func (c *class) Qualident() otm.Qualident {
	return otm.Qualident{Template: c.tpl, Class: c.cls}
}

type object struct {
	tpl, cls, id string
	om           map[uint32]*object
	em           map[uint32]*object
	vl           []interface{}
	up           *object
	clazz        otm.Class
}

type link struct {
	to *object
}

func (l *link) Object() otm.Object {
	return l.to
}

type futureLink struct {
	to string
	o  *object
	up *object
}

func (l *futureLink) Object() otm.Object {
	if l.o == nil {
		assert.For(!fn.IsNil(l.up), 20)
		l.o = otm.RootOf(l.up).FindById(l.to).(*object)
		assert.For(l.o != nil, 60, "object not found ", l.to)
	}
	return l.o
}

type stack struct {
	list.List
}

func newStack() *stack {
	ret := &stack{}
	ret.List = *list.New()
	return ret
}

func (s *stack) push(o *object) {
	s.PushFront(o)
}

func (s *stack) pop() (ret *object) {
	if s.Len() > 0 {
		ret = s.Remove(s.Front()).(*object)
	}
	return
}

func (s *stack) top() (ret *object) {
	if s.Len() > 0 {
		ret = s.Front().Value.(*object)
	}
	return
}

func (o *object) init() {
	o.om = make(map[uint32]*object)
	o.em = make(map[uint32]*object)
	o.clazz = &class{tpl: o.tpl, cls: o.cls}
}

func (o *object) omQualident() uint32 {
	return adler32.Checksum([]byte(fmt.Sprint(o.tpl, ".", o.cls, "(", o.id, ")")))
}

func (o *object) Parent() otm.Object { return o.up }

func (o *object) InstanceOf(override ...otm.Class) otm.Class {
	if len(override) > 0 {
		_, ok := o.clazz.(*class)
		assert.For(ok, 40, "already instantiated")
		o.clazz = override[0]
	}
	return o.clazz
}

func copyOf(o otm.Object) otm.Producer {
	return func(...otm.Modifier) otm.Object {
		return o.CopyOf(otm.DEEP)
	}
}
func (o *object) CopyOf(deep otm.CopyMode) (ret otm.Object) {
	b := Begin(o.Qualident())
	for _x := range o.Children() {
		switch x := _x.(type) {
		case otm.Object:
			if deep == otm.DEEP {
				b.Child(copyOf(x))
			}
		case string, int64, float64, rune, tri.Trit:
			b.Value(_x)
		case *futureLink:
			if deep == otm.DEEP {
				b.Value(x)
			}
		case *link:
			if deep == otm.DEEP {
				b.Value(&futureLink{to: x.Object().Qualident().Identifier})
			}
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	}
	return b.End()
}

func (o *object) Children() (c chan interface{}) {
	c = make(chan interface{})
	go func() {
		for _, _v := range o.vl {
			switch v := _v.(type) {
			case *futureLink:
				v.up = o
			}
			c <- _v
		}
		close(c)
	}()
	return
}

func (o *object) ChildrenObjects() (c chan otm.Object) {
	c = make(chan otm.Object)
	go func() {
		for _, _v := range o.vl {
			switch v := _v.(type) {
			case otm.Object:
				c <- v
			}
		}
		close(c)
	}()
	return
}

func (o *object) Qualident() otm.Qualident {
	return otm.Qualident{Template: o.tpl, Class: o.cls, Identifier: o.id}
}

func (o *object) ChildrenCount() uint {
	return uint(len(o.vl))
}

func (o *object) FindById(id string) (ret otm.Object) {
	var this func(o *object) (ret *object)
	this = func(o *object) (ret *object) {
		if o.id == id {
			ret = o
		}
		for _, _x := range o.vl {
			if x, ok := _x.(*object); ok {
				if z := this(x); z != nil {
					ret = z
				}
			}
		}
		return
	}
	ret = this(o)
	return
}

func omQualident(s *ir.Emit) uint32 {
	return adler32.Checksum([]byte(fmt.Sprint(s.Template, ".", s.Class, "(", s.Ident, ")")))
}

func omQualidentString(s *ir.Emit) string {
	return fmt.Sprint(otm.Qualident{Template: s.Template, Class: s.Class, Identifier: s.Ident})
}
func Map(t *ir.Template) (ret otm.Object) {
	st := newStack()
	uids := make(map[string]*object)
	var (
		emit  func() *object
		reuse func()
	)
	for _, _s := range t.Stmt {
		switch s := _s.(type) {
		case *ir.Emit:
			emit = func() *object {
				var o *object
				if _, ok := uids[s.Ident]; s.Ident != "" && ok {
					halt.As(100, "non-unique identifier ", s.Ident)
				}
				o = &object{tpl: s.Template, cls: s.Class, id: s.Ident}
				o.init()
				if o.id != "" {
					uids[o.id] = o
				}
				if parent := st.top(); parent != nil {
					if _, ok := parent.om[omQualident(s)]; !ok {
						parent.om[o.omQualident()] = o //remember first
					}
					if _, ok := parent.em[omQualident(s)]; ok {
						halt.As(100, "need reuse `", o.Qualident(), "`")
					}
					parent.vl = append(parent.vl, o)
					o.up = parent
				} else {
					ret = o
				}
				if s.ChildCount > 0 {
					st.push(o)
				}
				emit = nil
				reuse = nil
				return o
			}
			reuse = func() {
				if parent := st.top(); parent != nil {
					if old, ok := parent.em[omQualident(s)]; ok {
						st.push(old)
					} else if _, ok := parent.om[omQualident(s)]; !ok {
						old = emit()
						parent.em[omQualident(s)] = old
					} else {
						halt.As(100, "cannot reuse ", omQualidentString(s))
					}
				} else {
					halt.As(100, "nothing to reuse")
				}
				emit = nil
				reuse = nil
			}
			if s.ChildCount == 0 {
				assert.For(emit != nil, 20, "emitter is nil")
				emit()
			}
		case *ir.Dive:
			if s.Reuse {
				assert.For(reuse != nil, 20, "emitter is nil")
				reuse()
			} else {
				assert.For(emit != nil, 20, "emitter is nil")
				emit()
			}
		case *ir.Rise:
			st.pop()
		case *ir.Put:
			top := st.top()
			switch s.Type {
			case types.LINK:
				if o, ok := uids[s.Value.(string)]; ok {
					l := &link{to: o}
					top.vl = append(top.vl, l)
				} else {
					l := &futureLink{to: s.Value.(string), up: top}
					top.vl = append(top.vl, l)
				}
			case types.STRING:
				top.vl = append(top.vl, s.Value.(string))
			case types.INTEGER:
				if x, err := strconv.ParseInt(s.Value.(string), 10, 64); err == nil {
					top.vl = append(top.vl, x)
				} else {
					halt.As(100, err)
				}
			case types.REAL:
				if x, err := strconv.ParseFloat(s.Value.(string), 64); err == nil {
					top.vl = append(top.vl, x)
				} else {
					halt.As(100, err)
				}
			case types.CHAR:
				top.vl = append(top.vl, s.Value.(rune))
			case types.TRILEAN:
				top.vl = append(top.vl, tri.This(s.Value))
			default:
				halt.As(100, s.Type)
			}
		default:
			halt.As(100, reflect.TypeOf(s))
		}
	}
	var clean func(o *object)
	clean = func(o *object) {
		o.em = nil
		o.om = nil
		for x := range o.ChildrenObjects() {
			clean(x.(*object))
		}
	}
	return
}
