package chunk

import "errors"

// ErrChunkNotFound Chunk not found
var ErrChunkNotFound = errors.New("Chunk not found")

// Storage chunk storage
type Storage interface {
	LoadChunk(pos CPos) (*Chunk, error)
	SaveChunk(pos CPos, data *Chunk) error
}
