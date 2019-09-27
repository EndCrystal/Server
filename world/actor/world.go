package actor

type Systems map[string]System

func (s Systems) AddActorSystem(sys System) {
	s[sys.Name()] = sys
}

func (s Systems) Update() {
	for _, sys := range s {
		if upd, ok := sys.(SystemUpdatable); ok {
			upd.Update()
		}
	}
}

func (s Systems) AddActor(actor Actor) {
	for _, sys := range s {
		sys.Add(actor)
	}
}

func (s Systems) RemoveActor(id Id) {
	if id == Invalid {
		return
	}
	for _, sys := range s {
		sys.Remove(id)
	}
}
