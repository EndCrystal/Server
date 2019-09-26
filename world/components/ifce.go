package components

import "github.com/EndCrystal/Server/world/actor"

type PluginComponentsHost struct{}

func (PluginComponentsHost) RegisterComponent(name string, pointer interface{}, info ComponentInfo) Id {
	return Register(name, pointer, info)
}

func (PluginComponentsHost) GetComponetById(id Id) *Component {
	return GetById(id)
}

func (PluginComponentsHost) GetComponentByName(name string) *Component {
	return GetByName(name)
}

func (PluginComponentsHost) GetComponetSet(act actor.Actor) ComponentSet {
	return GetComponentSet(act)
}
