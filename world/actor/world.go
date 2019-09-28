package actor

import . "github.com/EndCrystal/Server/types"

type Systems struct {
	main          map[string]System
	list          []System
	recheckable   []SystemRecheckable
	preUpdateable []SystemPreUpdatable
	updateable    []SystemUpdatable
	simpleAddable []SystemSimpleAdd
}

func MakeSystems() (ret Systems) {
	ret.main = make(map[string]System)
	return
}

func (s Systems) GetByName(name string) System {
	if ret, ok := s.main[name]; ok {
		return ret
	}
	return nil
}

func (s Systems) AddActorSystem(sys System) {
	s.main[sys.Name()] = sys
	s.list = append(s.list, sys)
	if recheck, ok := sys.(SystemRecheckable); ok {
		s.recheckable = append(s.recheckable, recheck)
	}
	if preupd, ok := sys.(SystemPreUpdatable); ok {
		s.preUpdateable = append(s.preUpdateable, preupd)
	}
	if upd, ok := sys.(SystemUpdatable); ok {
		s.updateable = append(s.updateable, upd)
	}
	if add, ok := sys.(SystemSimpleAdd); ok {
		s.simpleAddable = append(s.simpleAddable, add)
	}
}

func (s Systems) Update() {
	for _, sys := range s.preUpdateable {
		sys.PreUpdate()
	}
	needRecheck := make(map[Actor]struct{})
	for _, sys := range s.updateable {
		for _, act := range sys.Update() {
			needRecheck[act] = struct{}{}
		}
	}
	for act := range needRecheck {
		s.Recheck(act)
	}
}

func (s Systems) AddActor(actor Actor) {
	for _, sys := range s.simpleAddable {
		sys.Add(actor)
	}
}

func (s Systems) RemoveActor(id Id) {
	if id == Invalid {
		return
	}
	for _, sys := range s.list {
		sys.Remove(id)
	}
}

func (s Systems) Recheck(actor Actor) {
	for _, sys := range s.recheckable {
		sys.Recheck(actor)
	}
}
