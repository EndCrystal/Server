package chunk

import packed "github.com/EndCrystal/PackedIO"

type Generator interface {
	packed.Serializable
	Generate(pos ChunkPos) *Chunk
}
