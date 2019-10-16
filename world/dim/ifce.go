package dim

import (
	"sync"

	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/chunk"
)

var (
	mtx                     = new(sync.Mutex)
	pendingActorSystemAdder []func(tags []string) actor.System
)

// PluginDimensionHost plugin host
type PluginDimensionHost struct{}

// AddDimension add dimension
func (PluginDimensionHost) AddDimension(name string, tags []string, storage chunk.Storage, generator chunk.Generator) {
	mtx.Lock()
	defer mtx.Unlock()
	d := Dimension{
		Systems: actor.MakeSystems(),
		tags:    tags,
	}
	d.Map.Init(storage, generator)
	for _, adder := range pendingActorSystemAdder {
		sys := adder(tags)
		if sys != nil {
			d.AddActorSystem(sys)
		}
	}
	dims[name] = &d
	return
}

// AddActorSystem add actor system
func (PluginDimensionHost) AddActorSystem(adder func(tags []string) actor.System) {
	mtx.Lock()
	defer mtx.Unlock()
	for _, d := range dims {
		sys := adder(d.tags)
		if sys != nil {
			d.AddActorSystem(sys)
		}
	}
	pendingActorSystemAdder = append(pendingActorSystemAdder, adder)
}

// AddPreloadActorSystem add preload actor system
func AddPreloadActorSystem(adder func() actor.System) {
	pendingActorSystemAdder = append(pendingActorSystemAdder, func([]string) actor.System { return adder() })
}
