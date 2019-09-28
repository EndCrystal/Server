package user

import (
	"container/list"
	"sync"

	"github.com/EndCrystal/Server/world/chunk"
	"github.com/EndCrystal/Server/world/dim"

	. "github.com/EndCrystal/Server/types"
)

type Event interface {
	EventName() string
}

type UserInfo struct {
	Username         string
	UserLabel        string
	Dimension        *dim.Dimension
	ControllingActor Id
	Pos              chunk.ChunkPos
	events           list.List
}

func (info UserInfo) GetChunkPosition() chunk.ChunkPos {
	return info.Pos
}

func (info UserInfo) AddEvent(event Event) {
	info.events.PushBack(event)
}

func (info UserInfo) HandleEvent(fn func(Event) bool) {
	e := info.events.Front()
	for e != nil {
		if fn(e.Value.(Event)) {
			next := e.Next()
			info.events.Remove(e)
			e = next
		} else {
			e = e.Next()
		}
	}
}

func (info UserInfo) ClearEvents() {
	info.events.Init()
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
