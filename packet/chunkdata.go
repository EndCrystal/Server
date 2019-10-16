package packet

import (
	"bytes"
	"compress/zlib"
	"sync"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/chunk"
)

// ChunkDataPacket chunk data packet
type ChunkDataPacket struct {
	Pos   chunk.CPos
	cache []byte
}

var chunkdataPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// SetData set chunk data
func (pkt *ChunkDataPacket) SetData(data *chunk.Chunk) func() {
	buf := chunkdataPool.Get().(*bytes.Buffer)
	buf.Reset()
	w, _ := zlib.NewWriterLevel(buf, zlib.BestCompression)
	o := packed.MakeOutput(w)
	data.Save(o)
	w.Close()
	pkt.cache = buf.Bytes()
	return func() { chunkdataPool.Put(buf) }
}

// Save save to data
func (pkt ChunkDataPacket) Save(out packed.Output) {
	pkt.Pos.Save(out)
	out.WriteBytes(pkt.cache)
}

// PacketID ID
func (ChunkDataPacket) PacketID() PID { return IDChunkData }
