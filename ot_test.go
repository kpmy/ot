package ot

import (
	"bufio"
	"bytes"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/ot/otm/conv"
	"github.com/kpmy/ot/otp"
	"github.com/kpmy/ot/ots"
	"github.com/kpmy/trigo"
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
			-1
			0DU -1 true false null inf -inf
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
				-1 true false null inf -inf
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
			node0: a b c d: d0 d1 d2; @x;
			node1: x(x) y z;
			node2: @x "a" "b" "c" "012345";
			attr.uniq0 :: u0 u1 1 2 3;
			uniq1 :: u2 u3 0.1 0.2 0.3;
			attr.uniq0 :: u4 u5 0U 1U 2U;
			uniq2(blab) :: x 0;
			uniq2(blab) :: y 0;
			u: -1 true false null inf -inf;
			li: kpmy@blab;
		;
	`
	p := otp.ConnectTo(ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(testTemplate))))
	if tpl, err := p.Template(); err == nil {
		m := conv.Map(tpl)
		prettyPrintObject(m)
		prettyPrintObject(m.CopyOf(otm.DEEP))
	} else {
		t.Fatal(err)
	}
}

func TestModules(t *testing.T) {
	/*
		<!DOCTYPE HTML>
		<html>
			<body>
				<p id="hello-teddy">превед, медвед</p>
				<br/><br/><br/>
				<p id="good-by-teddy">пока, медвед</p>
			</body>
		</html>
	*/
	const testTemplate = `
		core.template:
			import :: html;
			html:
				body:
					p(hello-teddy): "превед, медвед";
					br br br
					p(good-by-teddy): "пока, медвед";
				;
			;
		;
	`
	p := otp.ConnectTo(ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(testTemplate))))
	if tpl, err := p.Template(); err == nil {
		m := conv.Map(tpl)
		if err := conv.Resolve(m); err == nil {
			renderHtml(m)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}

func TestBuilder(t *testing.T) {
	root := otm.Qualident{Class: "root"}
	br := otm.Qualident{Class: "br"}
	br0 := otm.Qualident{Class: "br", Identifier: "null"}
	b := conv.Begin(root).Value("hello", 1945, 3.14, 2.71, "world", '!')
	b.Child(conv.Begin(br).Prod()).Child(conv.Begin(br).Child(conv.Begin(br).Prod()).Prod()).Child(conv.Begin(br0).Prod())
	b.Value(conv.Link("null"))
	o := b.End()
	prettyPrintObject(o)
	prettyPrintObject(o.CopyOf(otm.SHALLOW))
	prettyPrintObject(o.CopyOf(otm.DEEP))
}

func TestContext(t *testing.T) {
	const testTemplate = `
		core.template(бла-блабыч):
			import :: context;
			$(test) $(test-tri)
			$(test-path/test)
			$(test-list/0) $(test-list/1) "so template" $(test-list/2)
		;
	`
	p := otp.ConnectTo(ots.ConnectTo(bufio.NewReader(bytes.NewBufferString(testTemplate))))
	if tpl, err := p.Template(); err == nil {
		m := conv.Map(tpl)
		if err := conv.Resolve(m); err == nil {
			data := make(map[string]interface{})
			data["test"] = "test-string"
			data["test-tri"] = tri.TRUE
			data["test-path"] = data
			data["test-list"] = []interface{}{"one", "two", "three"}
			if err := conv.ResolveContext(m, data); err == nil {
				prettyPrintObject(m)
			} else {
				t.Fatal(err)
			}
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}
