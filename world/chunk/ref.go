package chunk

import (
	"sync"
	"time"
)

type ChunkRef struct {
	mtx *sync.Mutex
	*Chunk
	lastAccess time.Time
}

func (ref *ChunkRef) Access() func() {
	ref.mtx.Lock()
	ref.lastAccess = time.Now()
	return ref.mtx.Unlock
}
func (ref ChunkRef) Since() time.Duration { return time.Since(ref.lastAccess) }
