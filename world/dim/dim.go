package dim

import (
	"sync"

	"github.com/EndCrystal/Server/world/actor"
	"github.com/EndCrystal/Server/world/chunk"
)

// Dimension dimension struct
type Dimension struct {
	sync.Mutex
	actor.Systems
	chunk.Map
	tags []string
}

var dims = make(map[string]*Dimension)

// Tags get dimension tags
func (d *Dimension) Tags() []string {
	return d.tags
}

// LookupDimension lookup dimension by name
func LookupDimension(name string) (*Dimension, bool) {
	d, ok := dims[name]
	return d, ok
}

// ForEachDimension for each dimension
func ForEachDimension(fn func(name string, d *Dimension) error) error {
	for name, d := range dims {
		err := fn(name, d)
		if err != nil {
			return err
		}
	}
	return nil
}
