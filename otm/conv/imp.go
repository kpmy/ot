package conv

import (
	"errors"
	"fmt"
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/trigo"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/fn"
	"github.com/kpmy/ypk/halt"
	"log"
	"reflect"
	"strconv"
	"strings"
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

type ContextRefEntity struct {
}

type ContextIncludeEntity struct {
}

var (
	Core, Html, Context *ForeignTemplate
	tm                  map[string]*ForeignTemplate
)

func (c *ContextRefEntity) find(data map[string]interface{}, root *object) (v interface{}, err error) {
	if root.ChildrenCount() == 1 {
		var id = ""
		for o := range root.ChildrenObjects() {
			id = o.Qualident().Class
		}
		if id != "" {
			path := strings.Split(id, "/")
			var x interface{}
			x = data
			for _, s := range path {
				if s != "" {
					switch v := x.(type) {
					case map[string]interface{}:
						x = v[s]
					case []interface{}:
						var i int64 = -1
						if i, err = strconv.ParseInt(s, 10, 64); err == nil {
							x = v[int(i)]
						}
					default:
						err = errors.New(fmt.Sprint("not indexable ", reflect.TypeOf(v)))
					}
				}
				if err != nil {
					break
				}
			}
			if err == nil {
				v = x
			}
		} else {
			err = errors.New("context reference empty")
		}
	} else {
		err = errors.New(fmt.Sprint("context reference wrong format ", root.Qualident()))
	}
	return
}

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

func initContext() {
	Context = &ForeignTemplate{TemplateName: "context"}
	tm["context"] = Context
	Context.Classes = make(map[string]*ForeignClass)
	Context.Classes["$"] = &ForeignClass{Template: Context, Class: "$",
		Applicator: func(c *ForeignClass, o otm.Object) (_ func(*ForeignClass, otm.Object) error, err error) {
			c.Entity = &ContextRefEntity{}
			return nil, nil
		}}
	Context.Classes["$include"] = &ForeignClass{Template: Context, Class: "$include",
		Applicator: func(c *ForeignClass, o otm.Object) (_ func(*ForeignClass, otm.Object) error, err error) {
			c.Entity = &ContextIncludeEntity{}
			return nil, nil
		}}
}

func init() {
	tm = make(map[string]*ForeignTemplate)
	initCore()
	initHtml() //TODO выпилить позже или перенести в отдельный модуль
	initContext()
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

type ResolverFunc func(otm.Qualident) (otm.Object, error)

func resolveContext(_o *object, resolver ResolverFunc, data map[string]interface{}) (err error) {
	for i, _v := range _o.vl {
		switch v := _v.(type) {
		case *object:
			handled := false
			if clazz, ok := v.clazz.(*ForeignClass); ok {
				switch e := clazz.Entity.(type) {
				case *ContextRefEntity:
					if _val, _err := e.find(data, v); _err == nil {
						switch val := _val.(type) {
						case int:
							_o.vl[i] = int64(val)
						case string, int64, float64, rune, tri.Trit:
							_o.vl[i] = _val
						default:
							halt.As(100, reflect.TypeOf(val))
						}
						handled = true
					} else {
						err = errors.New(fmt.Sprint("context object not found: ", v.id, " ", _err))
					}
				case *ContextIncludeEntity:
					if v.ChildrenCount() == 1 {
						var id otm.Qualident
						for o := range v.ChildrenObjects() {
							id = o.Qualident()
						}
						if incl, _err := resolver(id); _err == nil {
							_i := incl.CopyOf(otm.DEEP).(*object)
							_o.vl[i] = _i
							_i.up = _o
							handled = true
						} else {
							err = _err
						}
					} else {
						err = errors.New("$include empty")
					}
				}
			}
			if !handled && err == nil {
				err = resolveContext(v, resolver, data)
			}
		}
		if err != nil {
			break
		}
	}
	return
}

func ResolveContext(o otm.Object, resolver ResolverFunc, data map[string]interface{}) (err error) {
	assert.For(!fn.IsNil(o), 20)

	switch tpl := o.Qualident().Template; tpl {
	case Core.TemplateName:
		err = resolveContext(o.(*object), resolver, data)
	default:
		err = errors.New("nothing to resolve")
	}
	return
}
