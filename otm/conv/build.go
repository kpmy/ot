package conv

import (
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/fn"
	"github.com/kpmy/ypk/halt"
	"reflect"
)

type builder struct {
	root *object
}

func Link(id string) otm.Link {
	return &futureLink{to: id}
}

func Begin(q otm.Qualident) otm.Builder {
	ret := &builder{}
	ret.root = &object{tpl: q.Template, cls: q.Class, id: q.Identifier}
	ret.root.init()
	return ret
}

func (b *builder) End(fl ...otm.Modifier) otm.Object {
	assert.For(b.root != nil, 20)
	return b.Prod()(fl...)
}

func (b *builder) Prod() otm.Producer {
	assert.For(b.root != nil, 20)
	return func(fl ...otm.Modifier) (ret otm.Object) {
		ret = b.root
		for _, fn := range fl {
			ret = fn(ret)
		}
		return
	}
}

func (b *builder) Value(vl ...interface{}) otm.Builder {
	assert.For(b.root != nil, 20)
	for _, _v := range vl {
		switch v := _v.(type) {
		case int:
			b.root.vl = append(b.root.vl, int64(v))
		case string, rune, int64, float64:
			b.root.vl = append(b.root.vl, _v)
		case *futureLink:
			v.o = nil
			v.up = nil
			b.root.vl = append(b.root.vl, v)
		default:
			halt.As(100, reflect.TypeOf(v))
		}
	}
	return b
}

func (b *builder) Child(prod otm.Producer, mod ...otm.Modifier) otm.Builder {
	assert.For(b.root != nil, 20)
	n := prod(mod...).(*object)
	if n.id != "" {
		old := b.root.FindById(n.id)
		assert.For(fn.IsNil(old), 40, "non unique id")
	}
	b.root.vl = append(b.root.vl, n)
	n.up = b.root
	return b
}
