package packet

import (
	"bytes"
	"sync"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/chunk"
)

type ChunkDataPacket struct {
	Pos   chunk.ChunkPos
	cache []byte
}

var chunkdataPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func (pkt *ChunkDataPacket) SetData(data *chunk.Chunk) func() {
	buf := chunkdataPool.Get().(*bytes.Buffer)
	buf.Reset()
	o := packed.MakeOutput(buf)
	data.Save(o)
	pkt.cache = buf.Bytes()
	return func() { chunkdataPool.Put(buf) }
}

func (pkt ChunkDataPacket) Save(out packed.Output) {
	pkt.Pos.Save(out)
	out.WriteFixedBytes(pkt.cache)
}

func (ChunkDataPacket) PacketId() PacketId { return IdChunkData }
