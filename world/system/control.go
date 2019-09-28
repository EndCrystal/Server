package system

import (
	. "github.com/EndCrystal/Server/types"
	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/components"
	"github.com/EndCrystal/Server/world/dim"
)

type controlSystemElem struct {
	actor.Actor
	components.ControllableComponent
}
type ControlSystem map[Id]controlSystemElem

func (ControlSystem) Name() string { return "core:control" }

func (s ControlSystem) Add(act actor.Actor) {
	if comp, ok := act.(components.ControllableComponent); ok {
		s[act.ID()] = controlSystemElem{act, comp}
	}
}

func (s ControlSystem) Remove(id Id) {
	delete(s, id)
}

func (s ControlSystem) Update() {
	for _, comp := range s {
		select {
		case mix := <-comp.Controllable().ControlRequest:
			comp.RuntimeComponentMap()[components.UserControlId] = components.UserControl{Owner: mix.Source}
		default:
		}
	}
}

func init() {
	dim.AddPreloadActorSystem(func() actor.System { return make(ControlSystem) })
}
