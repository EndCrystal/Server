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

// Used for tick_area
func (ref *ChunkRef) Update() {
	ref.lastAccess = time.Now()
}

func (ref *ChunkRef) Access() func() {
	ref.mtx.Lock()
	ref.Update()
	return ref.mtx.Unlock
}
func (ref ChunkRef) Since() time.Duration { return time.Since(ref.lastAccess) }
