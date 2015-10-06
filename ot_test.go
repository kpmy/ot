package ot

import (
	"bufio"
	"bytes"
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
			#par
		;
		child3:
			#par
		;
	;`
	p := otp.ConnectTo(ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(testTemplate))))
	if t, err := p.Template(); err == nil {
		prettyPrint(t)
	}
}
