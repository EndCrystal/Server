package actor

import "github.com/EndCrystal/Server/types"

// System basic interface
type System interface {
	Name() string
	Remove(types.ID)
}

// SystemRecheckable recheckable system
type SystemRecheckable interface {
	System
	Recheck(Actor)
}

// SystemSimpleAdd simple add system
type SystemSimpleAdd interface {
	System
	Add(Actor)
}

// SystemUpdatable updatable system
type SystemUpdatable interface {
	System
	Update() []Actor
}

// SystemPreUpdatable pre updatable system
type SystemPreUpdatable interface {
	System
	PreUpdate()
}
