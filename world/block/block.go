package block

import (
	"errors"

	packed "github.com/EndCrystal/PackedIO"
)

type Block struct {
	Name       string
	Attributes Attribute
}

type Attribute uint8

const (
	AttributeNone   Attribute = 0
	AttributeHasAux           = 1 << iota
	AttributeHasColor
	AttributeSolid
	AttributeFluid
)

func (b Block) HasAux() bool   { return b.Attributes&AttributeHasAux != 0 }
func (b Block) HasColor() bool { return b.Attributes&AttributeHasColor != 0 }
func (b Block) IsSolid() bool  { return b.Attributes&AttributeSolid != 0 }
func (b Block) IsFluid() bool  { return b.Attributes&AttributeFluid != 0 }

type Instance struct {
	Id    Id
	Aux   uint32
	Color uint32
}

func (i *Instance) GetBlock() *Block {
	return Get(i.Id)
}

var EInvalidInstance = errors.New("Invalid block instance")

func (i *Instance) Normalize() {
	blk := Get(i.Id)
	if blk == nil {
		panic(EInvalidInstance)
	}
	if !blk.HasAux() {
		i.Aux = 0
	}
	if !blk.HasColor() {
		i.Color = 0
	}
}

type ExtraLoader func(in packed.Input) packed.Serializable
type Id uint16

var (
	registry [^Id(0)]Block
	loaders     = make(map[Id]ExtraLoader)
	maxId    Id = 0
	index       = make(map[string]Id)
)

var EConflict = errors.New("Conflicted block")

func Register(b Block) Id {
	if _, ok := index[b.Name]; ok {
		panic(EConflict)
	}
	defer func() { maxId++ }()
	registry[maxId] = b
	index[b.Name] = maxId
	return maxId
}

func Get(id Id) *Block {
	if id < maxId {
		return &registry[id]
	}
	return nil
}

func LookupId(name string) (Id, bool) {
	ret, ok := index[name]
	return ret, ok
}

func Lookup(name string) *Block {
	if ret, ok := index[name]; ok {
		return &registry[ret]
	}
	return nil
}

func GetExtraLoader(id Id) ExtraLoader {
	loader, ok := loaders[id]
	if !ok {
		return nil
	}
	return loader
}

func RegisterExtraLoader(id Id, loader ExtraLoader) {
	loaders[id] = loader
}

func DescribeBlocks(o packed.Output) {
	o.WriteVarUint32(uint32(maxId))
	for i := Id(0); i < maxId; i++ {
		o.WriteString(registry[i].Name)
		o.WriteUint8(uint8(registry[i].Attributes))
	}
}
