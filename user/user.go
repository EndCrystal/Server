package user

import (
	"sync"

	"github.com/EndCrystal/Server/world/chunk"
	"github.com/EndCrystal/Server/world/dim"

	. "github.com/EndCrystal/Server/types"
)

type UserInfo struct {
	Username         string
	UserLabel        string
	Dimension        *dim.Dimension
	ControllingActor Id
	Pos              chunk.ChunkPos
}

func (info UserInfo) GetChunkPosition() chunk.ChunkPos {
	return info.Pos
}

var users = make(map[string]*UserInfo)
var mtx sync.RWMutex

func FindUser(name string) *UserInfo {
	mtx.RLock()
	defer mtx.RUnlock()
	if found, ok := users[name]; ok {
		return found
	}
	return nil
}

func AddUser(info *UserInfo) {
	mtx.Lock()
	defer mtx.Unlock()
	users[info.Username] = info
}

func RemoveUser(name string) {
	mtx.Lock()
	defer mtx.Unlock()
	delete(users, name)
}
