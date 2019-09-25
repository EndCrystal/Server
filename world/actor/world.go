package actor

type Systems []System

func (s *Systems) AddActorSystem(sys System) {
	*s = append(*s, sys)
	(*s)[len(*s)-1].Update()
}

func (s Systems) AddActor(actor Actor) {
	for _, sys := range s {
		sys.Add(actor)
	}
}

func (s Systems) RemoveActor(id Id) {
	for _, sys := range s {
		sys.Remove(id)
	}
}
