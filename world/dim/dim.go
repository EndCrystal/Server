package dim

import "github.com/EndCrystal/Server/world/actor"

type Dimension struct {
	actor.Systems
	tags []string
}

var dims map[string]*Dimension

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
