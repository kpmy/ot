package loader

import (
	"github.com/kpmy/ot/otm"
	"github.com/kpmy/rng/schema"
	"github.com/kpmy/rng/schema/std"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/halt"
	"log"
	"reflect"
)

var Constructors map[string]func() schema.Guide

func init() {
	Constructors = make(map[string]func() schema.Guide)
	Constructors["choice"] = std.Choice
	Constructors["element"] = std.Element
	Constructors["interleave"] = std.Interleave
	Constructors["attribute"] = std.Attribute
	Constructors["zeroOrMore"] = std.ZeroOrMore
	Constructors["oneOrMore"] = std.OneOrMore
	Constructors["group"] = std.Group
	Constructors["list"] = std.List
	Constructors["except"] = std.Except
	Constructors["mixed"] = std.Mixed
	Constructors["optional"] = std.Optional
	Constructors["anyName"] = std.AnyName
	Constructors["nsName"] = std.NSName
	Constructors["data"] = std.Data
	Constructors["text"] = std.Text
	Constructors["name"] = std.Name
	Constructors["empty"] = std.Empty
	Constructors["value"] = std.Value
	Constructors["param"] = std.Param
	Constructors["ref"] = std.Ref
	Constructors["externalRef"] = std.ExternalRef
}

func Construct(name otm.Qualident) (ret schema.Guide) {
	assert.For(name.Template == "ng", 20, name)
	fn := Constructors[name.Class]
	assert.For(fn != nil, 40, name)
	ret = fn()
	assert.For(ret != nil, 60, name)
	return
}

type Cached struct {
	o    otm.Object
	root schema.Guide
}

type Walker struct {
	root  otm.Object
	cache map[string]*Cached
	start schema.Start
	pos   schema.Guide
}

func (w *Walker) Init() *Walker {
	w.cache = make(map[string]*Cached)
	w.start = std.Start()
	w.pos = w.start
	return w
}

func (w *Walker) GrowDown(g schema.Guide) {
	w.pos.Add(g)
	g.Parent(w.pos)
	w.pos = g
}

func (w *Walker) Grow(g schema.Guide) {
	w.pos.Add(g)
}

func (w *Walker) Up() {
	w.pos = w.pos.Parent()
}

func (w *Walker) forEach(o otm.Object, do func(w *Walker, o otm.Object)) {
	for v := range o.ChildrenObjects() {
		do(w, v)
	}
}

func traverseWrap() func(w *Walker, o otm.Object) {
	return func(w *Walker, o otm.Object) {
		w.traverse(o)
	}
}

func (w *Walker) traverse(o otm.Object) {
	var (
		this    schema.Guide
		skip    *bool
		skipped = func() {
			s := true
			skip = &s
		}
		important = func() {
			s := false
			skip = &s
		}
	)

	switch o.Qualident().Class {
	//structure elements
	case "grammar":
		skipped()
		if start := o.FindByQualident(otm.Qualident{Template: "ng", Class: "start"}); start != nil {
			w.forEach(start[0], traverseWrap())
		}
	case "ref":
		skipped()
		panic(0)
		/*		if ref := w.root.FindByName(n.Name); ref != nil {
					if cached := w.cache[n.Name]; cached == nil {
						this = Construct(n.XMLName)
						{
							std.NameAttr(this, n.Name)
							//std.RefAttr(this, cached.root)
						}
						w.cache[n.Name] = &Cached{node: ref, root: this}
						w.GrowDown(this)
						w.forEach(ref, traverseWrap())
						w.Up()
					} else {
						w.Grow(cached.root)
					}
				} else {
					halt.As(100, "ref not found", n.Name)
				}
		*/
	//content elements
	case "element", "attribute", "data", "text", "value", "name", "param":
		fallthrough
	//constraint elements
	case "choice", "interleave", "optional", "zeroOrMore", "oneOrMore", "group", "list", "mixed", "except", "anyName", "nsName", "empty", "externalRef":
		important()
		this = Construct(o.Qualident())
		{
			/*std.NameAttr(this, n.Name)
			std.CharDataAttr(this, n.Data())
			std.TypeAttr(this, n.Type)
			std.NSAttr(this, n.NS)
			std.DataTypeAttr(this, n.DataType)
			std.CombineAttr(this, n.Combine)
			std.HrefAttr(this, n.Href)
			*/
		}
		w.GrowDown(this)
		w.forEach(o, traverseWrap())
		w.Up()
	//skipped elements
	case "description": //descriprion do nothing
	default:
		halt.As(100, o.Qualident())
	}
	if skip != nil {
		assert.For(*skip || this != nil, 60, "no result for ", o.Qualident())
	} else if this == nil {
		log.Println("unhandled", o.Qualident())
	}
}

func Load(source interface{}) (ret schema.Start) {
	switch s := source.(type) {
	case otm.Object:
		w := &Walker{root: s}
		w.Init()
		w.traverse(s)
		ret = w.start
	default:
		halt.As(100, reflect.TypeOf(s))
	}
	return
}
