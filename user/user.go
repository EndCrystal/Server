package user

import (
	"sync"

	"github.com/EndCrystal/Server/types"
	"github.com/EndCrystal/Server/world/chunk"
	"github.com/EndCrystal/Server/world/dim"
)

// Info userinfo
type Info struct {
	Username         string
	UserLabel        string
	Dimension        *dim.Dimension
	ControllingActor types.ID
	Pos              chunk.CPos
}

// GetChunkPosition get user chunk position
func (info Info) GetChunkPosition() chunk.CPos {
	return info.Pos
}

var users = make(map[string]*Info)
var mtx sync.RWMutex

// FindUser find user by name
func FindUser(name string) *Info {
	mtx.RLock()
	defer mtx.RUnlock()
	if found, ok := users[name]; ok {
		return found
	}
	return nil
}

// AddUser add userinfo
func AddUser(info *Info) {
	mtx.Lock()
	defer mtx.Unlock()
	users[info.Username] = info
}

// RemoveUser purge user info
func RemoveUser(name string) {
	mtx.Lock()
	defer mtx.Unlock()
	delete(users, name)
}
