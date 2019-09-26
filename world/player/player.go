package player

import (
	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/components"
)

type Player struct {
	actor.Identifier
	components.Nameable
	components.Position
	components.Rotation
}

var _ actor.Actor = new(Player)
