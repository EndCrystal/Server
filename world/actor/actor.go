package actor

import (
	"sync"
)

type (
	Id    uint32
	Actor interface {
		ID() Id
	}
)

var (
	maxId Id
	mtx   sync.Mutex
)

func AllocIdentifier() Id {
	mtx.Lock()
	defer func() {
		maxId++
		mtx.Unlock()
	}()
	return maxId
}
