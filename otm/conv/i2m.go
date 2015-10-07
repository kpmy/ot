package conv

import (
	"container/list"
	"fmt"
	"github.com/kpmy/ot/ir"
	"github.com/kpmy/ot/ir/types"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/halt"
	"hash/adler32"
	"reflect"
	"strconv"
)

type object struct {
	tpl, cls, id string
	om           map[uint32]*object
	vl           []interface{}
	up           *object
}

type link struct {
	to *object
}

func (l *link) Object() otm.Object {
	return l.to
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
}

func (o *object) omQualident() uint32 {
	return adler32.Checksum([]byte(fmt.Sprint(o.tpl, ".", o.cls, "(", o.id, ")")))
}

func (o *object) Parent() otm.Object { return o.up }

func (o *object) Children() (c chan interface{}) {
	c = make(chan interface{})
	go func() {
		for _, v := range o.vl {
			c <- v
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

func omQualident(s *ir.Emit) uint32 {
	return adler32.Checksum([]byte(fmt.Sprint(s.Template, ".", s.Class, "(", s.Ident, ")")))
}

func Map(t *ir.Template) (ret otm.Object) {
	st := newStack()
	uids := make(map[string]*object)
	var (
		emit  func()
		reuse func()
	)
	for _, _s := range t.Stmt {
		switch s := _s.(type) {
		case *ir.Emit:
			emit = func() {
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
			}
			reuse = func() {
				if parent := st.top(); parent != nil {
					if old, ok := parent.om[omQualident(s)]; ok {
						st.push(old)
					} else {
						emit()
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
					halt.As(100, "identifier not found ", s.Value)
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
			default:
				halt.As(100, s.Type)
			}
		default:
			halt.As(100, reflect.TypeOf(s))
		}
	}
	return
}
