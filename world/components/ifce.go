package components

import (
	"github.com/EndCrystal/Server/types"
	"github.com/EndCrystal/Server/world/actor"
)

// PluginComponentsHost plugin host
type PluginComponentsHost struct{}

// RegisterComponent register component
func (PluginComponentsHost) RegisterComponent(name string, pointer interface{}, info ComponentInfo) types.ID {
	return Register(name, pointer, info)
}

// GetComponetByID get component by id
func (PluginComponentsHost) GetComponetByID(id types.ID) *Component {
	return GetByID(id)
}

// GetComponentByName get component by name
func (PluginComponentsHost) GetComponentByName(name string) (types.ID, *Component) {
	return GetByName(name)
}

// GetComponetSet get components set
func (PluginComponentsHost) GetComponetSet(act actor.Actor) ComponentSet {
	return GetComponentSet(act)
}

// RegisterRuntimeComponent register runtime component
func (PluginComponentsHost) RegisterRuntimeComponent(name string, pointer interface{}, info RuntimeComponentInfo) types.ID {
	return RegisterRuntime(name, pointer, info)
}

// GetRuntimeComponentByID get runtime component by id
func (PluginComponentsHost) GetRuntimeComponentByID(id types.ID) *RuntimeComponent {
	return GetRuntimeByID(id)
}

// GetRuntimeComponentByName get runtime component by name
func (PluginComponentsHost) GetRuntimeComponentByName(name string) (types.ID, *RuntimeComponent) {
	return GetRuntimeByName(name)
}
