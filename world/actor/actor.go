package actor

import (
	"sync"

	. "github.com/EndCrystal/Server/types"
)

type (
	Identifier          struct{ Id }
	RuntimeComponentMap map[Id]interface{}
	Basic               struct {
		Identifier
		runtime RuntimeComponentMap
	}
	Actor interface {
		ID() Id
		RuntimeComponentMap() RuntimeComponentMap
	}
)

var (
	Invalid  = Id(0)
	maxId    = Id(1)
	freelist = make(map[Id]struct{})
	mtx      = new(sync.Mutex)
)

func (id Identifier) ID() Id                             { return id.Id }
func (b Basic) RuntimeComponentMap() RuntimeComponentMap { return b.runtime }

func AllocIdentifier() Id {
	mtx.Lock()
	defer mtx.Unlock()
	for id := range freelist {
		defer delete(freelist, id)
		return id
	}
	defer func() { maxId++ }()
	return maxId
}

func DeleteIdentifier(id Id) {
	mtx.Lock()
	defer mtx.Unlock()
	freelist[id] = struct{}{}
}
