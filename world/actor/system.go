package actor

type System interface {
	Name() string
	Remove(Id)
	Add(Actor)
}

type SystemUpdatable interface {
	System
	Update()
}

type BaseSystem map[Id]Actor

