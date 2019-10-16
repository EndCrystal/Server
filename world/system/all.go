package system

import (
	"github.com/EndCrystal/Server/types"
	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/dim"
)

// AllSystem all actor
type AllSystem map[types.ID]actor.Actor

// Name name
func (AllSystem) Name() string { return "core:all" }

// Add add actor
func (s AllSystem) Add(act actor.Actor) {
	s[act.ID()] = act
}

// Remove remove actor
func (s AllSystem) Remove(id types.ID) {
	delete(s, id)
}

func init() {
	dim.AddPreloadActorSystem(func() actor.System { return make(AllSystem) })
}
