package actor

import (
	"sync"
)

type (
	Id         uint32
	Identifier struct{ Id }
	Actor      interface {
		ID() Id
	}
)

var (
	maxId    Id
	freelist map[Id]struct{}
	mtx      sync.Mutex
)

func (id Identifier) ID() Id { return id.Id }

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
