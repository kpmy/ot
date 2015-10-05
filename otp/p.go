package otp

import (
	"errors"
	"github.com/kpmy/ot/ots"
	"github.com/kpmy/ypk/assert"
	"log"
)

type Marker interface {
	Mark(...interface{})
	FutureMark() Marker
}

type Parser interface {
	Template() error
}

type pr struct {
	common
}

func (p *pr) qualident() {
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
	log.Println(tid, id)
}

func (p *pr) block() {
	p.expect(ots.Ident, "identifier expected", ots.Delimiter)
	p.qualident()
	if p.await(ots.Lparen, ots.Delimiter) {
		p.run(ots.Rparen)
		p.next()
	}
	inner := func() {
		for stop := false; !stop; {
			p.pass(ots.Delimiter)
			switch p.sym.Code {
			case ots.Ident:
				p.block()
			case ots.Number:
				p.next()
			case ots.String:
				p.next()
			case ots.Link:
				p.next()
				p.expect(ots.Ident, "identifier expected")
				p.next()
			case ots.Semicolon:
				stop = true
			default:
				p.mark("unexpected ", p.sym)
			}
		}
	}
	if p.await(ots.Colon, ots.Delimiter) {
		p.next()
		inner()
		p.expect(ots.Semicolon, "semicolon expected")
		p.next()
	} else if p.is(ots.Square) {
		p.next()
		inner()
		p.expect(ots.Semicolon, "semicolon expected")
		p.next()
	}
}

func (p *pr) Template() (err error) {
	if err = p.sc.Error(); err != nil {
		return err
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
	err = nil
	return nil
}

func ConnectTo(sc ots.Scanner) Parser {
	ret := &pr{}
	ret.sc = sc
	ret.debug = true
	ret.next()
	return ret
}
