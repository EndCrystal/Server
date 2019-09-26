package dim

import (
	"sync"

	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/chunk"
)

type PluginDimensionHost struct {
	mtx                     *sync.Mutex
	pendingActorSystemAdder []func(tags []string) actor.System
}

func (host PluginDimensionHost) AddDimension(name string, tags []string, storage chunk.Storage, generator chunk.Generator) {
	host.mtx.Lock()
	defer host.mtx.Unlock()
	d := Dimension{
		Mutex:   new(sync.Mutex),
		Systems: make(actor.Systems, 0),
		tags:    tags,
	}
	d.Map.Init(storage, generator)
	for _, adder := range host.pendingActorSystemAdder {
		sys := adder(tags)
		if sys != nil {
			d.AddActorSystem(sys)
		}
	}
	dims[name] = &d
	return
}

func (host PluginDimensionHost) AddActorSystem(adder func(tags []string) actor.System) {
	host.mtx.Lock()
	defer host.mtx.Unlock()
	for _, d := range dims {
		sys := adder(d.tags)
		if sys != nil {
			d.AddActorSystem(sys)
		}
	}
	host.pendingActorSystemAdder = append(host.pendingActorSystemAdder, adder)
}
