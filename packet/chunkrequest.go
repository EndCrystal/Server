package packet

import (
	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/chunk"
)

// ChunkRequestPacket chunk request from client-side
type ChunkRequestPacket struct {
	Pos chunk.CPos
}

// PacketID id
func (ChunkRequestPacket) PacketID() PID { return IDChunkRequest }

// Load load from packet data
func (pkt *ChunkRequestPacket) Load(i packed.Input) {
	pkt.Pos.Load(i)
}

// Check check for parse
func (pkt ChunkRequestPacket) Check(pctx *ParseContext) bool { return pctx.Check(64) }
