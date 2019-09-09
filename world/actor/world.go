package actor

var systems []System

func AddSystem(sys System) {
	systems = append(systems, sys)
	systems[len(systems)-1].Update()
}

func AddActor(actor Actor) {
	for _, sys := range systems {
		sys.Add(actor)
	}
}

func RemoveAcctor(id Id) {
	for _, sys := range systems {
		sys.Remove(id)
	}
}
