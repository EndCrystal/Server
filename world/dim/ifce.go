package dim

import (
	"sync"

	"github.com/EndCrystal/Server/world/actor"
)

type PluginDimensionHost struct {
	mtx                     *sync.Mutex
	pendingActorSystemAdder []func(tags []string) actor.System
}

func (host PluginDimensionHost) AddDimension(name string, tags []string) {
	host.mtx.Lock()
	defer host.mtx.Unlock()
	var d Dimension
	for _, adder := range host.pendingActorSystemAdder {
		sys := adder(tags)
		if sys != nil {
			d.AddActorSystem(sys)
		}
	}
	d.tags = tags
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
