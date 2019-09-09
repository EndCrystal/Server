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
	maxId Id
	mtx   sync.Mutex
)

func (id Identifier) ID() Id { return id.Id }

func AllocIdentifier() Id {
	mtx.Lock()
	defer func() {
		maxId++
		mtx.Unlock()
	}()
	return maxId
}
