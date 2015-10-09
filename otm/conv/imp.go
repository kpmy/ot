package conv

import (
	"errors"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/fn"
	"log"
)

type ForeignTemplate struct {
	TemplateName string
	Classes      map[string]*BuiltInClass
}

type BuiltInClass struct {
	tpl, cls string
}

func (c *BuiltInClass) Qualident() otm.Qualident {
	return otm.Qualident{Template: c.tpl, Class: c.cls}
}

var Core *ForeignTemplate

func initCore() {
	Core = &ForeignTemplate{TemplateName: "core"}
	Core.Classes = make(map[string]*BuiltInClass)
	Core.Classes["template"] = &BuiltInClass{tpl: Core.TemplateName, cls: "template"}
	Core.Classes["import"] = &BuiltInClass{tpl: Core.TemplateName, cls: "import"}
}

func init() {
	initCore()
}

func Resolve(o otm.Object) (err error) {
	assert.For(!fn.IsNil(o), 20)

	var upd func(t *ForeignTemplate, o otm.Object)
	upd = func(t *ForeignTemplate, o otm.Object) {
		if clazz, ok := Core.Classes[o.Qualident().Class]; ok && (o.Qualident().Template == t.TemplateName || o.Qualident().Template == "") {
			o.InstanceOf(clazz)
			log.Println("class updated for", o.Qualident(), " set ", clazz.Qualident())
		}
		for x := range o.ChildrenObjects() {
			upd(t, x)
		}
	}
	switch tpl := o.Qualident().Template; tpl {
	case Core.TemplateName:
		upd(Core, o)
	default:
		err = errors.New("nothing to resolve")
	}
	return
}
