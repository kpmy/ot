package ots

import (
	"fmt"
	"github.com/kpmy/ypk/assert"
	"io"
	"strings"
	"unicode"
)

type SymCode int

const (
	None SymCode = iota
	Ident
	Number
	String

	Delimiter
	Qualifier
	Colon
	Semicolon
	Lparen
	Rparen
	Square
	Link

	True
	False
	Null
	Inf
	Minus
)

var kv map[string]SymCode

func init() {
	kv = map[string]SymCode{"true": True, "false": False, "null": Null, "inf": Inf}
}

type Symbol struct {
	Code       SymCode
	Value      string
	NumberOpts struct {
		Modifier string
		Period   bool
	}
	StringOpts struct {
		Apos bool
	}
}

func (sym SymCode) String() (s string) {
	switch sym {
	case None:
		s = "none"
	case Delimiter:
		s = "space"
	case Ident:
		s = "ident"
	case Qualifier:
		s = "~"
	case Colon:
		s = ":"
	case Square:
		s = "::"
	case Semicolon:
		s = ";"
	case Lparen:
		s = "("
	case Rparen:
		s = ")"
	case Number:
		s = "number"
	case String:
		s = "string"
	case Link:
		s = "@"
	case Minus:
		s = "-"
	case True:
		s = "true"
	case False:
		s = "false"
	case Null:
		s = "null"
	case Inf:
		s = "inf"
	default:
		s = fmt.Sprint(sym)
	}
	return
}

func (s Symbol) String() string {
	return fmt.Sprint("sym: `", s.Code, "` ", s.Value)
}

type Scanner interface {
	Get() Symbol
	Error() error
	Read() int
	Pos() (int, int)
}

type sc struct {
	rd  io.RuneReader
	err error

	ch  rune
	pos int

	lines struct {
		count int
		last  int
		crlf  bool
		lens  map[int]func() (int, int)
	}
}

func (s *sc) Read() int { return s.pos }
func (s *sc) Pos() (int, int) {
	return s.lines.count, s.lines.last
}
func (s *sc) Error() error { return s.err }

func (s *sc) mark(msg ...interface{}) {
	//log.Println("at pos ", s.pos, " ", fmt.Sprintln(msg...))
	l, c := s.Pos()
	panic(Err("scanner", s.Read(), l, c, msg...))
}

func (s *sc) next() rune {
	read := 0
	s.ch, read, s.err = s.rd.ReadRune()
	if s.err == nil {
		s.pos += read
	}
	if s.ch == '\r' || s.ch == '\n' {
		s.line()
	} else {
		s.lines.last++
	}
	//log.Println(Token(s.ch), s.err)
	return s.ch
}

func (s *sc) line() {
	if s.ch == '\r' {
		s.lines.crlf = true
	}
	if (s.lines.crlf && s.ch == '\r') || (!s.lines.crlf && s.ch == '\n') {
		s.lines.lens[s.lines.count] = func() (int, int) {
			return s.lines.count, s.pos
		}
		s.lines.count++
		s.lines.last = 1
	} else if s.lines.crlf && s.ch == '\n' {
		s.lines.last--
	}
}

// @#$%^&*-_=+,.?!/|\\
func isIdentLetter(r rune) bool {
	return isIdentFirstLetter(r) || unicode.IsDigit(r) || strings.ContainsRune(`@#$%^&*-_=+,.?!/|\\`, r)
}

func isIdentFirstLetter(r rune) bool {
	return unicode.IsLetter(r) || strings.ContainsRune(`$`, r)
}

func (s *sc) ident() (sym Symbol) {
	assert.For(isIdentFirstLetter(s.ch), 20, "letter must be first")
	buf := make([]rune, 0)
	for s.err == nil && isIdentLetter(s.ch) {
		buf = append(buf, s.ch)
		s.next()
	}
	if s.err == nil {
		sym.Value = string(buf)
		if code, ok := kv[sym.Value]; ok {
			sym.Code = code
		} else {
			sym.Code = Ident
		}
	} else {
		s.mark("error while reading ident")
	}
	return
}

func (s *sc) comment() {
	assert.For(s.ch == '*', 20, "expected * ", "got ", Token(s.ch))
	for {
		for s.err == nil && s.ch != '*' {
			if s.ch == '(' {
				if s.next() == '*' {
					s.comment()
				}
			} else {
				s.next()
			}
		}
		for s.err == nil && s.ch == '*' {
			s.next()
		}
		if s.err != nil || s.ch == ')' {
			break
		}
	}
	if s.err == nil {
		s.next()
	} else {
		s.mark("unclosed comment")
	}
}

const dec = "0123456789"
const hhex = "ABCDEF"
const hex = dec + hhex

const modifier = "U"

//first char always 0..9
func (s *sc) num() (sym Symbol) {
	assert.For(unicode.IsDigit(s.ch), 20, "digit expected")
	var buf []rune
	var mbuf []rune
	hasDot := false

	for {
		buf = append(buf, s.ch)
		s.next()
		if s.ch == '.' {
			if !hasDot {
				hasDot = true
			} else if hasDot {
				s.mark("dot unexpected")
			}
		}
		if s.err != nil || !(s.ch == '.' || strings.ContainsRune(hex, s.ch)) {
			break
		}
	}
	if strings.ContainsRune(modifier, s.ch) {
		mbuf = append(mbuf, s.ch)
		s.next()
	}
	if strings.ContainsAny(string(buf), hhex) && len(mbuf) == 0 {
		s.mark("modifier expected")
	}
	if s.err == nil {
		sym.Code = Number
		sym.Value = string(buf)
		sym.NumberOpts.Modifier = string(mbuf)
		sym.NumberOpts.Period = hasDot
	} else {
		s.mark("error reading number")
	}
	return
}

func (s *sc) str() string {
	assert.For(s.ch == '"' || s.ch == '\'' || s.ch == '`', 20, "quote expected")
	var buf []rune
	ending := s.ch
	s.next()
	for ; s.err == nil && s.ch != ending; s.next() {
		buf = append(buf, s.ch)
	}
	if s.err == nil {
		s.next()
	} else {
		s.mark("string expected")
	}
	return string(buf)
}

func (s *sc) get() (sym Symbol) {
	switch s.ch {
	case '~':
		sym.Code = Qualifier
		s.next()
	case '(':
		if s.next() == '*' {
			s.comment()
		} else {
			sym.Code = Lparen
		}
	case ')':
		sym.Code = Rparen
		s.next()
	case ':':
		if s.next() == ':' {
			sym.Code = Square
			s.next()
		} else {
			sym.Code = Colon
		}
	case '@':
		sym.Code = Link
		s.next()
	case ';':
		sym.Code = Semicolon
		s.next()
	case '-':
		sym.Code = Minus
		s.next()
	case '"', '\'', '`':
		sym.StringOpts.Apos = (s.ch == '\'' || s.ch == '`')
		sym.Value = s.str()
		sym.Code = String
	default:
		switch {
		case isIdentFirstLetter(s.ch):
			sym = s.ident()
		case unicode.IsSpace(s.ch):
			for unicode.IsSpace(s.ch) {
				s.next()
			}
			sym.Code = Delimiter
		case unicode.IsDigit(s.ch):
			sym = s.num()
		default:
			s.mark("unhandled ", "`", Token(s.ch), "`")
			s.next()
		}
	}
	return
}

func (s *sc) Get() (sym Symbol) {
	for stop := s.err != nil; !stop; {
		sym = s.get()
		stop = sym.Code != 0 || s.err != nil
	}
	return
}

func ConnectTo(rd io.RuneReader) Scanner {
	ret := &sc{}
	ret.rd = rd
	ret.lines.lens = make(map[int]func() (int, int))
	ret.lines.count++
	ret.next()
	return ret
}
