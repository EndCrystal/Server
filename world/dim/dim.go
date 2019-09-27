package dim

import (
	"sync"

	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/chunk"
)

type Dimension struct {
	*sync.Mutex
	actor.Systems
	chunk.Map
	tags []string
}

var dims = make(map[string]*Dimension)

func (d Dimension) Tags() []string {
	return d.tags
}

func LookupDimension(name string) (*Dimension, bool) {
	d, ok := dims[name]
	return d, ok
}

func ForEachDimension(fn func(name string, d *Dimension) error) error {
	for name, d := range dims {
		err := fn(name, d)
		if err != nil {
			return err
		}
	}
	return nil
}
