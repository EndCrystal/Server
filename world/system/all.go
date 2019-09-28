package system

import (
	. "github.com/EndCrystal/Server/types"
	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/dim"
)

type AllSystem map[Id]actor.Actor

func (AllSystem) Name() string { return "core:all" }

func (s AllSystem) Add(act actor.Actor) {
	s[act.ID()] = act
}

func (s AllSystem) Remove(id Id) {
	delete(s, id)
}

func init() {
	dim.AddPreloadActorSystem(func() actor.System { return make(AllSystem) })
}
