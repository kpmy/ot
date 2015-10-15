package trav

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/ot/otm/path"
	"github.com/kpmy/ot/otp"
	"github.com/kpmy/ot/ots"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/halt"
	"log"
	"strings"
)

type level struct {
	q otm.Qualident
}

func (l *level) Qualident() otm.Qualident {
	return l.q
}

type trav struct {
	l []*level
}

func (t *trav) Path() (ret []path.Level) {
	for _, l := range t.l {
		ret = append(ret, l)
	}
	return
}

func (t *trav) Run(o otm.Object, typ path.ResultType) (ret interface{}, err error) {
	this := o
	for i, l := range t.l {
		if this.Qualident() == l.q {
			if i+1 < len(t.l) {
				for o := range this.ChildrenObjects() {
					if o.Qualident() == t.l[i+1].Qualident() {
						this = o
					}
				}
			} else {
				switch typ {
				case path.OBJECT:
					ret = this
				default:
					halt.As(100, typ)
				}
			}
		} else {
			err = errors.New(fmt.Sprint("wrong path ", l.q))
			break
		}
	}
	return
}

func ident(s string) (ret otm.Qualident, err error) {
	sc := ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(s)))
	p := otp.Commons(sc)
	p.Expect(ots.Ident)
	ret.Class = p.Sym().Value
	if c := p.Next().Code; c == ots.Qualifier {
		p.Next()
		p.Expect(ots.Ident)
		ret.Template = ret.Class
		ret.Class = p.Sym().Value
		if c := p.Next().Code; c == ots.Lparen {
			p.Next()
			p.Expect(ots.Ident)
			ret.Identifier = p.Sym().Value
			p.Expect(ots.Rparen)
		}
	} else if c != ots.None {
		halt.As(100, c)
	}
	return
}

func (t *trav) compile(path string) (err error) {
	assert.For(path != "", 20)
	for _, s := range strings.Split(path, ":") {
		assert.For(s != "", 21)
		l := &level{}
		if l.q, err = ident(s); err == nil {
			log.Println(l.q)
			t.l = append(t.l, l)
		} else {
			break
		}
	}
	return
}

func Trav(path string) (ret path.Traverser, err error) {
	t := &trav{}
	if err = t.compile(path); err == nil {
		ret = t
	}
	return
}
