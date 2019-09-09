package actor

type System interface {
	Update()
	Remove(Id)
	Add(Actor)
}

type SystemInitializer interface{ Init() }
