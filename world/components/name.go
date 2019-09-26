package components

import packed "github.com/EndCrystal/PackedIO"

type Nameable struct {
	Name string
}

func (n Nameable) GetName() string      { return n.Name }
func (n *Nameable) SetName(name string) { n.Name = name }

type NameComponent interface {
	GetName() string
	SetName(string)
}

type nameInfo struct{}

func (nameInfo) Secure() bool { return false }
func (nameInfo) LoadComponent(i packed.Input, obj interface{}) {
	obj.(NameComponent).SetName(i.ReadString())
}
func (nameInfo) SaveComponent(o packed.Output, obj interface{}) {
	o.WriteString(obj.(NameComponent).GetName())
}

func init() {
	Register("core:name", (*NameComponent)(nil), nameInfo{})
}
