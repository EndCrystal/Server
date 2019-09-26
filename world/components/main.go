package components

import (
	"reflect"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/world/actor"
)

type Component struct {
	Id   Id
	Type reflect.Type
	Name string
	Info ComponentInfo
}

type ComponentInfo interface {
	LoadComponent(packed.Input, interface{})
	SaveComponent(packed.Output, interface{})
	Secure() bool
}

type Id uint32

var registry = make([]*Component, 0)
var index = make(map[string]Id)

func Register(name string, pointer interface{}, info ComponentInfo) (id Id) {
	t := reflect.TypeOf(pointer).Elem()
	if t.Kind() != reflect.Interface {
		panic("Invalid component: not an interface")
	}
	id = Id(len(registry))
	registry = append(registry, &Component{
		Id:   id,
		Type: t,
		Name: name,
		Info: info,
	})
	index[name] = id
	return
}

func GetById(id Id) *Component {
	return registry[id]
}

func GetByName(name string) *Component {
	if id, ok := index[name]; ok {
		return registry[id]
	}
	return nil
}

type ComponentSet map[*Component]struct{}

func (set ComponentSet) Has(comps ...*Component) bool {
	for _, comp := range comps {
		if _, ok := set[comp]; !ok {
			return false
		}
	}
	return true
}

var actor_cache map[reflect.Type]ComponentSet
var actor_cache_send map[reflect.Type]ComponentSet

func GetComponentSetForSend(act actor.Actor) ComponentSet {
	t := reflect.TypeOf(act)
	if ret, ok := actor_cache_send[t]; ok {
		return ret
	}
	ret := make(ComponentSet)
	temp := GetComponentSet(act)
	for comp := range temp {
		if !comp.Info.Secure() {
			ret[comp] = struct{}{}
		}
	}
	actor_cache_send[t] = ret
	return ret
}

func GetComponentSet(act actor.Actor) ComponentSet {
	t := reflect.TypeOf(act)
	if ret, ok := actor_cache[t]; ok {
		return ret
	}
	ret := make(ComponentSet)
	for _, comp := range registry {
		if t.Implements(comp.Type) {
			ret[comp] = struct{}{}
		}
	}
	actor_cache[t] = ret
	return ret
}

func DescribeComponents(o packed.Output) {
	o.WriteVarUint32(uint32(len(registry)))
	for _, comp := range registry {
		o.WriteString(comp.Name)
	}
}

func SaveActor(o packed.Output, act actor.Actor) {
	set := GetComponentSet(act)
	o.WriteVarUint32(uint32(len(set)))
	for comp := range set {
		o.WriteString(comp.Name)
		comp.Info.SaveComponent(o, act)
	}
}

func SendActor(o packed.Output, act actor.Actor) {
	set := GetComponentSetForSend(act)
	o.WriteVarUint32(uint32(len(set)))
	for comp := range set {
		o.WriteUint32(uint32(comp.Id))
		comp.Info.SaveComponent(o, act)
	}
}

func LoadActor(i packed.Input, act actor.Actor) {
	log := logprefix.Get("[component loader] ")
	set := GetComponentSet(act)
	i.IterateObject(func(key string) {
		if comp := GetByName(key); comp != nil {
			if set.Has(comp) {
				comp.Info.LoadComponent(i, act)
			} else {
				log.Printf("Try load non-implemented component %s for actor %d", key, act.ID())
			}
		} else {
			log.Printf("Try load non-registered component %s for actor %d", key, act.ID())
		}
	})
}
