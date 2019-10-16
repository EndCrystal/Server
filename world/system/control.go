package system

import (
	"github.com/EndCrystal/Server/types"
	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/components"
	"github.com/EndCrystal/Server/world/dim"
)

type controlSystemElem struct {
	actor.Actor
	components.ControllableComponent
}

// ControlSystem control system
type ControlSystem map[types.ID]controlSystemElem

// Name name
func (ControlSystem) Name() string { return "core:control" }

// Add add actor
func (s ControlSystem) Add(act actor.Actor) {
	if comp, ok := act.(components.ControllableComponent); ok {
		s[act.ID()] = controlSystemElem{act, comp}
	}
}

// Remove remove actor
func (s ControlSystem) Remove(id types.ID) {
	delete(s, id)
}

// Update update system
func (s ControlSystem) Update() (list []actor.Actor) {
	for _, comp := range s {
		select {
		case mix := <-comp.Controllable().ControlRequest:
			comp.RuntimeComponentMap()[components.UserControlID] = components.UserControl{Owner: mix.Source}
			list = append(list, comp.Actor)
		default:
		}
	}
	return
}

func init() {
	dim.AddPreloadActorSystem(func() actor.System { return make(ControlSystem) })
}
