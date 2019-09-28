package components

import (
	"reflect"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/logprefix"
	. "github.com/EndCrystal/Server/types"
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

type RuntimeComponent struct {
	Id   Id
	Type reflect.Type
	Name string
	Info RuntimeComponentInfo
}

type RuntimeComponentInfo interface {
	CreateComponent(packed.Input) interface{}
	SaveComponent(packed.Output, interface{})
	Secure() bool
}

var (
	registry         = make([]*Component, 0)
	runtime_registry = make([]*RuntimeComponent, 0)
	index            = make(map[string]Id)
	runtime_index    = make(map[string]Id)
)

// Register static component type
// (if info == nil they won't be saved or sent)
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

func RegisterRuntime(name string, pointer interface{}, info RuntimeComponentInfo) (id Id) {
	t := reflect.TypeOf(pointer).Elem()
	id = Id(len(runtime_registry))
	runtime_registry = append(runtime_registry, &RuntimeComponent{
		Id:   id,
		Type: t,
		Name: name,
		Info: info,
	})
	return
}

func GetById(id Id) *Component {
	return registry[id]
}

func GetByName(name string) (Id, *Component) {
	if id, ok := index[name]; ok {
		return id, registry[id]
	}
	return 0, nil
}

func GetRuntimeById(id Id) *RuntimeComponent {
	return runtime_registry[id]
}

func GetRuntimeByName(name string) (Id, *RuntimeComponent) {
	if id, ok := runtime_index[name]; ok {
		return id, runtime_registry[id]
	}
	return 0, nil
}

type ComponentSet map[*Component]Id

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
	for comp, id := range temp {
		if comp.Info != nil && !comp.Info.Secure() {
			ret[comp] = id
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
	for id, comp := range registry {
		if comp.Info != nil && t.Implements(comp.Type) {
			ret[comp] = Id(id)
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
	count := uint32(0)
	for id := range act.RuntimeComponentMap() {
		if info := runtime_registry[id].Info; info != nil {
			count++
		}
	}
	o.WriteVarUint32(count)
	for id, comp := range act.RuntimeComponentMap() {
		entry := runtime_registry[id]
		if info := entry.Info; info != nil {
			o.WriteString(entry.Name)
			info.SaveComponent(o, comp)
		}
	}
}

func SendActor(o packed.Output, act actor.Actor) {
	set := GetComponentSetForSend(act)
	o.WriteVarUint32(uint32(len(set)))
	for comp := range set {
		o.WriteUint32(uint32(comp.Id))
		comp.Info.SaveComponent(o, act)
	}
	count := uint32(0)
	for id := range act.RuntimeComponentMap() {
		if info := runtime_registry[id].Info; info != nil && !info.Secure() {
			count++
		}
	}
	o.WriteVarUint32(count)
	for id, comp := range act.RuntimeComponentMap() {
		if info := runtime_registry[id].Info; info != nil && !info.Secure() {
			o.WriteVarUint32(uint32(id))
			info.SaveComponent(o, comp)
			count--
			if count == 0 {
				break // Early exit
			}
		}
	}
}

func LoadActor(i packed.Input, act actor.Actor) {
	log := logprefix.Get("[component loader] ")
	set := GetComponentSet(act)
	i.IterateObject(func(key string) {
		if _, comp := GetByName(key); comp != nil {
			if set.Has(comp) {
				if comp.Info != nil {
					comp.Info.LoadComponent(i, act)
				}
			} else {
				log.Printf("Try load non-implemented static component %s for actor %d", key, act.ID())
			}
		} else {
			log.Printf("Try load non-registered static component %s for actor %d", key, act.ID())
		}
	})
	i.IterateObject(func(key string) {
		if id, comp := GetRuntimeByName(key); comp != nil {
			if info := comp.Info; info != nil {
				act.RuntimeComponentMap()[id] = info.CreateComponent(i)
			} else {
				log.Printf("Try to load non-implemented runtime component %s for actor %d", key, act.ID())
			}
		} else {
			log.Printf("Try to load non-registered runtime component %s for actor %d", key, act.ID())
		}
	})
}
