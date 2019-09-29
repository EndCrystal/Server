package packet

import (
	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/chunk"
)

type ChunkRequestPacket struct {
	Pos chunk.ChunkPos
}

func (ChunkRequestPacket) PacketId() PacketId { return IdChunkRequest }

func (pkt *ChunkRequestPacket) Load(i packed.Input) {
	pkt.Pos.Load(i)
}

func (pkt ChunkRequestPacket) Check(pctx *ParseContext) bool { return pctx.Check(64) }
