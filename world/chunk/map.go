package chunk

import (
	"sync"
	"time"
)

// Map map data
type Map struct {
	mtx       sync.Mutex
	inited    bool
	storage   Storage
	generator Generator
	loaded    map[CPos]*Ref
}

// Init init map data
func (m *Map) Init(storage Storage, generator Generator) {
	if m.inited {
		return
	}
	m.inited = true
	m.storage = storage
	m.generator = generator
	m.loaded = make(map[CPos]*Ref)
}

// GetChunk get chunk from pos
func (m *Map) GetChunk(pos CPos) (ret *Ref, err error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	// lookup cache
	var ok bool
	if ret, ok = m.loaded[pos]; ok {
		return
	}
	// try load from storage
	var chk *Chunk
	defer func() {
		if err == nil {
			ret = &Ref{
				mtx:        new(sync.Mutex),
				Chunk:      chk,
				lastAccess: time.Now(),
			}
			m.loaded[pos] = ret
		}
	}()
	chk, err = m.storage.LoadChunk(pos)
	if err == nil || err != ErrChunkNotFound {
		return
	}
	// try to generate on the fly
	err = nil // ignore EChunkNotFound
	chk = m.generator.Generate(pos)
	return
}
