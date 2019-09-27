package user

import (
	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/chunk"
	"github.com/EndCrystal/Server/world/dim"
	"github.com/EndCrystal/Server/world/system"
)

type UserInfo struct {
	Username         string
	UserLabel        string
	Dimension        dim.Dimension
	ControllingActor actor.Id
	Pos              chunk.ChunkPos
}

func (info UserInfo) GetChunkPosition() chunk.ChunkPos {
	if info.ControllingActor == actor.Invalid {
		return info.Pos
	}
	pos := info.Dimension.Systems["core:position"].(system.PositionSystem)[info.ControllingActor].GetPosition()
	return chunk.ChunkPos{
		X: int32(pos.X / 16.),
		Z: int32(pos.Z / 16.),
	}
}
