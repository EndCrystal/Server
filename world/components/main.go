package components

import (
	"reflect"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/types"
	"github.com/EndCrystal/Server/world/actor"
)

// Component base
type Component struct {
	id   types.ID
	Type reflect.Type
	Name string
	Info ComponentInfo
}

// ComponentInfo info for component
type ComponentInfo interface {
	LoadComponent(packed.Input, interface{})
	SaveComponent(packed.Output, interface{})
	Secure() bool
}

// RuntimeComponent runtime component
type RuntimeComponent struct {
	id   types.ID
	Type reflect.Type
	Name string
	Info RuntimeComponentInfo
}

// RuntimeComponentInfo runtime component info
type RuntimeComponentInfo interface {
	CreateComponent(packed.Input) interface{}
	SaveComponent(packed.Output, interface{})
	Secure() bool
}

var (
	registry        = make([]*Component, 0)
	runtimeRegistry = make([]*RuntimeComponent, 0)
	index           = make(map[string]types.ID)
	runtimeIndex    = make(map[string]types.ID)
)

// Register static component type
// (if info == nil they won't be saved or sent)
func Register(name string, pointer interface{}, info ComponentInfo) (id types.ID) {
	t := reflect.TypeOf(pointer).Elem()
	if t.Kind() != reflect.Interface {
		panic("Invalid component: not an interface")
	}
	id = types.ID(len(registry))
	registry = append(registry, &Component{
		id:   id,
		Type: t,
		Name: name,
		Info: info,
	})
	index[name] = id
	return
}

// RegisterRuntime register runtime component type
func RegisterRuntime(name string, pointer interface{}, info RuntimeComponentInfo) (id types.ID) {
	t := reflect.TypeOf(pointer).Elem()
	id = types.ID(len(runtimeRegistry))
	runtimeRegistry = append(runtimeRegistry, &RuntimeComponent{
		id:   id,
		Type: t,
		Name: name,
		Info: info,
	})
	return
}

// GetByID get component by id
func GetByID(id types.ID) *Component {
	return registry[id]
}

// GetByName get component by name
func GetByName(name string) (types.ID, *Component) {
	if id, ok := index[name]; ok {
		return id, registry[id]
	}
	return 0, nil
}

// GetRuntimeByID get runtime component by id
func GetRuntimeByID(id types.ID) *RuntimeComponent {
	return runtimeRegistry[id]
}

// GetRuntimeByName get runtime component by name
func GetRuntimeByName(name string) (types.ID, *RuntimeComponent) {
	if id, ok := runtimeIndex[name]; ok {
		return id, runtimeRegistry[id]
	}
	return 0, nil
}

// ComponentSet set for components
type ComponentSet map[*Component]types.ID

// Has check set has components
func (set ComponentSet) Has(comps ...*Component) bool {
	for _, comp := range comps {
		if _, ok := set[comp]; !ok {
			return false
		}
	}
	return true
}

var actorCache map[reflect.Type]ComponentSet
var actorCacheSend map[reflect.Type]ComponentSet

// GetComponentSetForSend get components for send
func GetComponentSetForSend(act actor.Actor) ComponentSet {
	t := reflect.TypeOf(act)
	if ret, ok := actorCacheSend[t]; ok {
		return ret
	}
	ret := make(ComponentSet)
	temp := GetComponentSet(act)
	for comp, id := range temp {
		if comp.Info != nil && !comp.Info.Secure() {
			ret[comp] = id
		}
	}
	actorCacheSend[t] = ret
	return ret
}

// GetComponentSet get components
func GetComponentSet(act actor.Actor) ComponentSet {
	t := reflect.TypeOf(act)
	if ret, ok := actorCache[t]; ok {
		return ret
	}
	ret := make(ComponentSet)
	for id, comp := range registry {
		if comp.Info != nil && t.Implements(comp.Type) {
			ret[comp] = types.ID(id)
		}
	}
	actorCache[t] = ret
	return ret
}

// DescribeComponents describe components
func DescribeComponents(o packed.Output) {
	o.WriteVarUint32(uint32(len(registry)))
	for _, comp := range registry {
		o.WriteString(comp.Name)
	}
}

// SaveActor save actor
func SaveActor(o packed.Output, act actor.Actor) {
	set := GetComponentSet(act)
	o.WriteVarUint32(uint32(len(set)))
	for comp := range set {
		o.WriteString(comp.Name)
		comp.Info.SaveComponent(o, act)
	}
	count := uint32(0)
	for id := range act.RuntimeComponentMap() {
		if info := runtimeRegistry[id].Info; info != nil {
			count++
		}
	}
	o.WriteVarUint32(count)
	for id, comp := range act.RuntimeComponentMap() {
		entry := runtimeRegistry[id]
		if info := entry.Info; info != nil {
			o.WriteString(entry.Name)
			info.SaveComponent(o, comp)
		}
	}
}

// SendActor send actor
func SendActor(o packed.Output, act actor.Actor) {
	set := GetComponentSetForSend(act)
	o.WriteVarUint32(uint32(len(set)))
	for comp := range set {
		o.WriteUint32(uint32(comp.id))
		comp.Info.SaveComponent(o, act)
	}
	count := uint32(0)
	for id := range act.RuntimeComponentMap() {
		if info := runtimeRegistry[id].Info; info != nil && !info.Secure() {
			count++
		}
	}
	o.WriteVarUint32(count)
	for id, comp := range act.RuntimeComponentMap() {
		if info := runtimeRegistry[id].Info; info != nil && !info.Secure() {
			o.WriteVarUint32(uint32(id))
			info.SaveComponent(o, comp)
			count--
			if count == 0 {
				break // Early exit
			}
		}
	}
}

// LoadActor load actor
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
