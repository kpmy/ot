package ot

import (
	"bufio"
	"github.com/kpmy/ot/ir"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/ot/otm/conv"
	"github.com/kpmy/ot/otp"
	"github.com/kpmy/ot/ots"
	"io"
)

func Load(rd io.Reader) (o otm.Object, err error) {
	p := otp.ConnectTo(ots.ConnectTo(bufio.NewReader(rd)))
	var t *ir.Template
	if t, err = p.Template(); err == nil {
		o = conv.Map(t)
	}
	return
}
