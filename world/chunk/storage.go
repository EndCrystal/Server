package chunk

import "errors"

var EChunkNotFound = errors.New("Chunk not found")

type Storage interface {
	LoadChunk(pos ChunkPos) (*Chunk, error)
	SaveChunk(pos ChunkPos, data *Chunk) error
}
