package components

import (
	. "github.com/EndCrystal/Server/types"
	"github.com/EndCrystal/Server/world/actor"
)

type PluginComponentsHost struct{}

func (PluginComponentsHost) RegisterComponent(name string, pointer interface{}, info ComponentInfo) Id {
	return Register(name, pointer, info)
}

func (PluginComponentsHost) GetComponetById(id Id) *Component {
	return GetById(id)
}

func (PluginComponentsHost) GetComponentByName(name string) (Id, *Component) {
	return GetByName(name)
}

func (PluginComponentsHost) GetComponetSet(act actor.Actor) ComponentSet {
	return GetComponentSet(act)
}

func (PluginComponentsHost) RegisterRuntimeComponent(name string, pointer interface{}, info RuntimeComponentInfo) Id {
	return RegisterRuntime(name, pointer, info)
}

func (PluginComponentsHost) GetRuntimeComponentById(id Id) *RuntimeComponent {
	return GetRuntimeById(id)
}

func (PluginComponentsHost) GetRuntimeComponentByName(name string) (Id, *RuntimeComponent) {
	return GetRuntimeByName(name)
}
