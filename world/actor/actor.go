package actor

import (
	"sync"

	"github.com/EndCrystal/Server/types"
)

type (
	// Identifier Identifier
	Identifier struct{ id types.ID }
	// RuntimeComponentMap runtime component map
	RuntimeComponentMap map[types.ID]interface{}
	// Basic actor
	Basic struct {
		Identifier
		runtime RuntimeComponentMap
	}
	// Actor actor interface
	Actor interface {
		ID() types.ID
		RuntimeComponentMap() RuntimeComponentMap
	}
)

var (
	// Invalid invalid ID
	Invalid  = types.ID(0)
	maxID    = types.ID(1)
	freelist = make(map[types.ID]struct{})
	mtx      = new(sync.Mutex)
)

// ID retrive id
func (id Identifier) ID() types.ID { return id.id }

// RuntimeComponentMap get runtime component map
func (b Basic) RuntimeComponentMap() RuntimeComponentMap { return b.runtime }

// AllocIdentifier alloc identifier
func AllocIdentifier() types.ID {
	mtx.Lock()
	defer mtx.Unlock()
	for id := range freelist {
		defer delete(freelist, id)
		return id
	}
	defer func() { maxID++ }()
	return maxID
}

// DeleteIdentifier delete identifier
func DeleteIdentifier(id types.ID) {
	mtx.Lock()
	defer mtx.Unlock()
	freelist[id] = struct{}{}
}
