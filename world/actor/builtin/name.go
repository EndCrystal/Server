package builtin

type Nameable struct {
	Name string
}

func (n Nameable) GetName() string      { return n.Name }
func (n *Nameable) SetName(name string) { n.Name = name }

type NameComponent interface {
	GetName() string
	SetName(string)
}
