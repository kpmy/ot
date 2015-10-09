package conv

import (
	"errors"
	"fmt"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/fn"
	"github.com/kpmy/ypk/halt"
	"log"
)

type ForeignTemplate struct {
	TemplateName string
	Classes      map[string]*ForeignClass
}

type ForeignClass struct {
	Template   *ForeignTemplate
	Class      string
	Applicator func(*ForeignClass, otm.Object) (func(*ForeignClass, otm.Object) error, error)
	Entity     interface{}
}

func (c *ForeignClass) copyOf() *ForeignClass {
	assert.For(fn.IsNil(c.Entity), 20)
	ret := &ForeignClass{}
	*ret = *c
	return ret
}

func (c *ForeignClass) Qualident() otm.Qualident {
	return otm.Qualident{Template: c.Template.TemplateName, Class: c.Class}
}

type pfunc func() (pfunc, error)

func (c *ForeignClass) apply(o otm.Object) pfunc {
	return func() (pret pfunc, err error) {
		if c.Applicator != nil {
			var ret func(*ForeignClass, otm.Object) error
			if ret, err = c.Applicator(c, o); err == nil && ret != nil {
				pret = func() (pfunc, error) {
					return nil, ret(c, o)
				}
			}
		}
		return
	}
}

type TemplateEntity struct {
}

type ImportEntity struct {
	Imports map[otm.Qualident]otm.Object
	Ref     otm.Object
}

var Core *ForeignTemplate
var Html *ForeignTemplate

var tm map[string]*ForeignTemplate

func initCore() {
	Core = &ForeignTemplate{TemplateName: "core"}
	Core.Classes = make(map[string]*ForeignClass)
	Core.Classes["template"] = &ForeignClass{
		Template: Core,
		Class:    "template",
		Applicator: func(c *ForeignClass, o otm.Object) (_ func(*ForeignClass, otm.Object) error, err error) {
			c.Entity = &TemplateEntity{}
			x := otm.RootOf(o)
			if x != o {
				err = errors.New("core.template must be root")
			}
			return
		}}
	Core.Classes["import"] = &ForeignClass{
		Template: Core,
		Class:    "import",
		Applicator: func(c *ForeignClass, o otm.Object) (post func(*ForeignClass, otm.Object) error, err error) {
			e := &ImportEntity{Imports: make(map[otm.Qualident]otm.Object), Ref: o}
			c.Entity = e
			for imp := range o.ChildrenObjects() {
				e.Imports[imp.Qualident()] = imp
			}
			for k, _ := range e.Imports {
				if k.Template != "" || k.Identifier != "" || k.Class == Core.TemplateName {
					err = errors.New(fmt.Sprintln("cannot import ", k))
				}
			}
			post = func(c *ForeignClass, o otm.Object) (err error) {
				e := c.Entity.(*ImportEntity)
				for k, _ := range e.Imports {
					if t, ok := tm[k.Class]; ok {
						if p := o.Parent(); !fn.IsNil(p) {
							resolve(t, p)
						} else {
							halt.As(100, "cannot apply import to nil")
						}
					} else {
						err = errors.New(fmt.Sprintln("unknown template ", k))
					}
				}
				return
			}
			return
		}}
}

func initHtml() {
	Html = &ForeignTemplate{TemplateName: "html"}
	tm["html"] = Html
	Html.Classes = make(map[string]*ForeignClass)
	Html.Classes["html"] = &ForeignClass{Template: Html, Class: "html"}
	Html.Classes["body"] = &ForeignClass{Template: Html, Class: "body"}
	Html.Classes["p"] = &ForeignClass{Template: Html, Class: "p"}
	Html.Classes["br"] = &ForeignClass{Template: Html, Class: "br"}
}

func init() {
	tm = make(map[string]*ForeignTemplate)
	initCore()
	initHtml()
}
func resolve(t *ForeignTemplate, o otm.Object) (err error) {
	assert.For(!fn.IsNil(o), 20)
	var processList []pfunc

	var upd func(t *ForeignTemplate, o otm.Object)
	upd = func(t *ForeignTemplate, o otm.Object) {
		if clazz, ok := t.Classes[o.Qualident().Class]; ok && (o.Qualident().Template == t.TemplateName || o.Qualident().Template == "") {
			inst := clazz.copyOf()
			o.InstanceOf(inst)
			if fn := inst.apply(o); fn != nil {
				processList = append(processList, fn)
			}
			log.Println("class updated for", o.Qualident(), " set ", clazz.Qualident())
		}
		for x := range o.ChildrenObjects() {
			upd(t, x)
		}
	}

	upd(t, o)

	for tmp := processList; len(tmp) > 0; {
		var _tmp []pfunc
		for _, f := range tmp {
			var p pfunc
			if p, err = f(); err == nil && p != nil {
				_tmp = append(_tmp, p)
			} else if err != nil {
				_tmp = nil
				break
			}
		}
		tmp = _tmp
	}
	return
}

func Resolve(o otm.Object) (err error) {
	assert.For(!fn.IsNil(o), 20)

	switch tpl := o.Qualident().Template; tpl {
	case Core.TemplateName:
		err = resolve(Core, o)
	default:
		err = errors.New("nothing to resolve")
	}
	return
}
