package packet

import (
	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/chunk"
)

type ChunkDataPacket struct {
	Pos  chunk.ChunkPos
	Data *chunk.Chunk
}

func (pkt *ChunkDataPacket) Load(in packed.Input) {
	pkt.Pos.Load(in)
	pkt.Data = new(chunk.Chunk)
	pkt.Data.Load(in)
}

func (pkt ChunkDataPacket) Save(out packed.Output) {
	pkt.Pos.Save(out)
	pkt.Data.Save(out)
}

func (ChunkDataPacket) PacketId() PacketId            { return IdChunkData }
func (ChunkDataPacket) Check(pctx *ParseContext) bool { return pctx.Check(ServerSide, 0) }
