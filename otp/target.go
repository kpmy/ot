package otp

import (
	"github.com/kpmy/ot/ir"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/fn"
)

type target struct {
	tpl *ir.Template
}

func (t *target) init() {
	t.tpl = new(ir.Template)
}

func (t *target) emit(_s ir.Statement) {
	assert.For(!fn.IsNil(_s), 20)
	switch s := _s.(type) {
	default:
		t.tpl.Stmt = append(t.tpl.Stmt, s)
	}
}
