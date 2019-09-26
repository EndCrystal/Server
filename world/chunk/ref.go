package chunk

import (
	"sync"
	"time"
)

type ChunkRef struct {
	*sync.Mutex
	*Chunk
	lastAccess time.Time
}

func (ref *ChunkRef) Access()             { ref.lastAccess = time.Now() }
func (ref ChunkRef) Since() time.Duration { return time.Since(ref.lastAccess) }
