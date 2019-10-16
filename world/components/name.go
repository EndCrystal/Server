package components

import packed "github.com/EndCrystal/PackedIO"

// Nameable nameable
type Nameable struct {
	Name string
}

// GetName get name
func (n Nameable) GetName() string      { return n.Name }
// SetName set name
func (n *Nameable) SetName(name string) { n.Name = name }

// NameComponent name compoent
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

// NameID id
var NameID = Register("core:name", (*NameComponent)(nil), nameInfo{})
