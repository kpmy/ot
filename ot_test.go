package ot

import (
	"bufio"
	"bytes"
	"github.com/kpmy/ot/otm/conv"
	"github.com/kpmy/ot/otp"
	"github.com/kpmy/ot/ots"
	"log"
	"testing"
)

func init() {
	log.SetFlags(0)
}

func TestScanner(t *testing.T) {
	const scannerTestTemplate = `(* test template no semantic rules applied *)
		CORE.TEMPLATE:
			import :: context html;
			html.body:
			br: tere; 123 "fas


			df" 'f'
				prop:;

			0.125
			0DU
			seek: fas blab(bubub) "dfsdf";
			;
		;
	`
	sc := ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(scannerTestTemplate)))
	for i := 0; sc.Error() == nil && i < 100; i++ {
		log.Println(sc.Get())
	}
	log.Println(sc.Error())
}

func TestParser(t *testing.T) {
	const testTemplate = `
	block:
		blob.child0(par):
			unique ::
				блаб
			;
			"стринг"
			3323
			1.333
			0DU
			child1:
				child3
			;
		;
		child2:
			@par
		;
		child3:
			@par
		;
	;`
	p := otp.ConnectTo(ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(testTemplate))))
	if tpl, err := p.Template(); err == nil {
		prettyPrint(tpl)
	} else {
		t.Fatal(err)
	}
}

func TestModel(t *testing.T) {
	const testTemplate = `
		root:
			node0: a b c d: d0 d1 d2;;
			node1: x(x) y z;
			node2: @x "a" "b" "c" "012345";
			attr.uniq0 :: u0 u1 1 2 3;
			uniq1 :: u2 u3 0.1 0.2 0.3;
			attr.uniq0 :: u4 u5 0U 1U 2U;
			uniq2(blab) :: x 0;
			uniq2(blab) :: y 0;
			u
		;
	`
	p := otp.ConnectTo(ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(testTemplate))))
	if tpl, err := p.Template(); err == nil {
		m := conv.Map(tpl)
		prettyPrintObject(m)
	} else {
		t.Fatal(err)
	}
}

func TestModules(t *testing.T) {
	const testTemplate = `
		core.template:
			import :: блаб;
			import :: хуй;
			br br br
			br: test;
		;
	`
	p := otp.ConnectTo(ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(testTemplate))))
	if tpl, err := p.Template(); err == nil {
		m := conv.Map(tpl)
		if err := conv.Resolve(m); err == nil {

		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
