package otp

import (
	"errors"
	"github.com/kpmy/ot/ir"
	"github.com/kpmy/ot/ir/types"
	"github.com/kpmy/ot/ots"
	"github.com/kpmy/ypk/assert"
	"log"
)

type Marker interface {
	Mark(...interface{})
	FutureMark() Marker
}

type Parser interface {
	Template() (*ir.Template, error)
}

type pr struct {
	common
	target
}

func (p *pr) qualident() (string, string) {
	assert.For(p.is(ots.Ident), 20, "identifier expected here")
	id := ""
	tid := p.ident()
	p.next()
	if p.is(ots.Period) {
		p.expect(ots.Ident, "identifier expected", ots.Period)
		id = p.ident()
		p.next()
	} else {
		id = tid
		tid = ""
	}
	return tid, id
}

func (p *pr) block() {
	p.expect(ots.Ident, "identifier expected", ots.Delimiter)
	tid, id := p.qualident()
	uid := ""
	if p.await(ots.Lparen, ots.Delimiter) {
		p.next()
		p.expect(ots.Ident, "identifier expected", ots.Delimiter)
		uid = p.ident()
		p.next()
		p.expect(ots.Rparen, ") expected", ots.Delimiter)
		p.next()
	}
	this := &ir.Emit{Template: tid, Class: id, Ident: uid}
	p.emit(this)
	inner := func(e *ir.Emit) {
		for stop := false; !stop; {
			e.ChildCount++
			p.pass(ots.Delimiter)
			switch p.sym.Code {
			case ots.Ident:
				p.block()
			case ots.Number:
				st := &ir.Put{}
				st.Type, st.Value = p.number()
				p.emit(st)
				p.next()
			case ots.String:
				p.emit(&ir.Put{Type: types.STRING, Value: p.sym.Value})
				p.next()
			case ots.Link:
				p.next()
				p.expect(ots.Ident, "identifier expected")
				p.emit(&ir.Put{Type: types.LINK, Value: p.ident()})
				p.next()
			case ots.Semicolon:
				e.ChildCount--
				stop = true
			default:
				p.mark("unexpected ", p.sym)
			}
		}
	}
	down := func(reuse bool) {
		p.emit(&ir.Dive{Reuse: reuse})
		p.next()
		inner(this)
		p.expect(ots.Semicolon, "semicolon expected")
		p.emit(&ir.Rise{})
		if this.ChildCount == 0 {
			p.mark("empty block :/:: is redundant")
		}
		p.next()
	}
	if p.await(ots.Colon, ots.Delimiter) {
		down(false)
	} else if p.is(ots.Square) {
		down(true)
	} else {
		//empty
	}
}

func (p *pr) Template() (ret *ir.Template, err error) {
	if err = p.sc.Error(); err != nil {
		return nil, err
	}
	if !p.debug {
		defer func() {
			if x := recover(); x != nil {
				log.Println(x) // later errors from parser
			}
		}()
	}
	err = errors.New("compiler error")
	p.block()
	ret = p.tpl
	err = nil
	return
}

func ConnectTo(sc ots.Scanner) Parser {
	ret := &pr{}
	ret.sc = sc
	ret.debug = false
	ret.next()
	ret.init()
	return ret
}
