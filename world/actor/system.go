package actor

import . "github.com/EndCrystal/Server/types"

type System interface {
	Name() string
	Remove(Id)
}

type SystemRecheckable interface {
	System
	Recheck(Actor)
}

type SystemSimpleAdd interface {
	System
	Add(Actor)
}

type SystemUpdatable interface {
	System
	Update() []Actor
}

type SystemPreUpdatable interface {
	System
	PreUpdate()
}

