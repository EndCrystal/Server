package system

import (
	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/components"
)

type PositionSystem map[actor.Id]components.PositionComponent

func (PositionSystem) Name() string { return "core:position" }

func (s PositionSystem) Add(act actor.Actor) {
	if pos, ok := act.(components.PositionComponent); ok {
		s[act.ID()] = pos
	}
}

func (s PositionSystem) Remove(id actor.Id) {
	delete(s, id)
}

func init() {
	PreloadedSystems = append(PreloadedSystems, func() actor.System { return make(PositionSystem) })
}
