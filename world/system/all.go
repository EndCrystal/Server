package system

import "github.com/EndCrystal/Server/world/actor"

type AllSystem map[actor.Id]actor.Actor

func (AllSystem) Name() string { return "core:all" }

func (s AllSystem) Add(act actor.Actor) {
	s[act.ID()] = act
}

func (s AllSystem) Remove(id actor.Id) {
	delete(s, id)
}

func init() {
	PreloadedSystems = append(PreloadedSystems, func() actor.System { return make(AllSystem) })
}
