package chunk

import packed "github.com/EndCrystal/PackedIO"

// Generator chunk generator
type Generator interface {
	packed.Serializable
	Generate(pos CPos) *Chunk
}
