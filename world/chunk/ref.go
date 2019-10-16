package chunk

import (
	"sync"
	"time"
)

// Ref reference for chunk
type Ref struct {
	mtx *sync.Mutex
	*Chunk
	lastAccess time.Time
}

// Update Used for tick_area
func (ref *Ref) Update() {
	ref.lastAccess = time.Now()
}

// Access make access
func (ref *Ref) Access() func() {
	ref.mtx.Lock()
	ref.Update()
	return ref.mtx.Unlock
}

// Since calculate duration between now to last access
func (ref Ref) Since() time.Duration { return time.Since(ref.lastAccess) }
