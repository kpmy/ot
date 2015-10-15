package otp

import (
	"fmt"
	"github.com/kpmy/ot/ir/types"
	"github.com/kpmy/ot/ots"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/halt"
	"log"
	"strconv"
)

type mark struct {
	rd        int
	line, col int
	marker    Marker
}

func (m *mark) Mark(msg ...interface{}) {
	m.marker.(*common).m = m
	m.marker.Mark(msg...)
}

func (m *mark) FutureMark() Marker { halt.As(100); return nil }

type Common interface {
	Sym() ots.Symbol
	Next() ots.Symbol
	Expect(ots.SymCode)
}

type ext struct {
	common
}

func (e *ext) Expect(sym ots.SymCode) {
	e.expect(sym, "unexpected sym")
}

func (e *ext) Next() ots.Symbol {
	return e.next()
}

func (e *ext) Sym() ots.Symbol {
	return e.sym
}

func Commons(sc ots.Scanner) Common {
	e := &ext{}
	e.sc = sc
	e.next()
	return e
}

type common struct {
	sc    ots.Scanner
	sym   ots.Symbol
	done  bool
	debug bool
	m     *mark
}

func (p *common) Mark(msg ...interface{}) {
	p.mark(msg...)
}

func (p *common) FutureMark() Marker {
	rd := p.sc.Read()
	str, pos := p.sc.Pos()
	m := &mark{marker: p, rd: rd, line: str, col: pos}
	return m
}

func (p *common) mark(msg ...interface{}) {
	rd := p.sc.Read()
	str, pos := p.sc.Pos()
	if len(msg) == 0 {
		p.m = &mark{rd: rd, line: str, col: pos}
	} else if p.m != nil {
		rd, str, pos = p.m.rd, p.m.line, p.m.col
		p.m = nil
	}
	if p.m == nil {
		panic(ots.Err("parser", rd, str, pos, msg...))
	}
}

func (p *common) next() ots.Symbol {
	p.done = true
	if p.sym.Code != ots.None {
		//		fmt.Print("this ")
		//		fmt.Print("`" + fmt.Sprint(p.sym) + "`")
	}
	p.sym = p.sc.Get()
	//fmt.Print(" next ")
	if p.debug {
		log.Println("`" + fmt.Sprint(p.sym) + "`")
	}
	return p.sym
}

//expect is the most powerful step forward runner, breaks the compilation if unexpected sym found
func (p *common) expect(sym ots.SymCode, msg string, skip ...ots.SymCode) {
	assert.For(p.done, 20)
	if !p.await(sym, skip...) {
		p.mark(msg)
	}
	p.done = false
}

//await runs for the sym through skip list, but may not find the sym
func (p *common) await(sym ots.SymCode, skip ...ots.SymCode) bool {
	assert.For(p.done, 20)
	skipped := func() (ret bool) {
		for _, v := range skip {
			if v == p.sym.Code {
				ret = true
			}
		}
		return
	}

	for sym != p.sym.Code && skipped() && p.sc.Error() == nil {
		p.next()
	}
	p.done = p.sym.Code != sym
	return p.sym.Code == sym
}

//pass runs through skip list
func (p *common) pass(skip ...ots.SymCode) {
	skipped := func() (ret bool) {
		for _, v := range skip {
			if v == p.sym.Code {
				ret = true
			}
		}
		return
	}
	for skipped() && p.sc.Error() == nil {
		p.next()
	}
}

//run runs to the first sym through any other sym
func (p *common) run(sym ots.SymCode) {
	if p.sym.Code != sym {
		for p.sc.Error() == nil && p.next().Code != sym {
			if p.sc.Error() != nil {
				p.mark("not found")
				break
			}
		}
	}
}

func (p *common) ident() string {
	assert.For(p.sym.Code == ots.Ident, 20, "identifier expected")
	//добавить валидацию идентификаторов
	return p.sym.Value
}

func (p *common) number() (t types.Type, v interface{}) {
	assert.For(p.is(ots.Number), 20, "number expected here")
	switch p.sym.NumberOpts.Modifier {
	case "":
		if p.sym.NumberOpts.Period {
			t, v = types.REAL, p.sym.Value
		} else {
			//x, err := strconv.Atoi(p.sym.Str)
			//assert.For(err == nil, 40)
			t, v = types.INTEGER, p.sym.Value
		}
	case "U":
		if p.sym.NumberOpts.Period {
			p.mark("hex integer value expected")
		}
		//fmt.Println(p.sym)
		if r, err := strconv.ParseUint(p.sym.Value, 16, 64); err == nil {
			t = types.CHAR
			v = rune(r)
		} else {
			p.mark("error while reading integer")
		}
	case "H":
		if p.sym.NumberOpts.Period {
			p.mark("hex integer value expected")
		}
		t, v = types.INTEGER, p.sym.Value
	default:
		p.mark("unknown number format `", p.sym.NumberOpts.Modifier, "`")
	}
	return
}

func (p *common) is(sym ots.SymCode) bool {
	return p.sym.Code == sym
}
