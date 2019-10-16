package block

import (
	"errors"

	packed "github.com/EndCrystal/PackedIO"
)

// Block block struct
type Block struct {
	Name       string
	Attributes Attribute
}

// Attribute attribute
type Attribute uint8

const (
	// AttributeNone none
	AttributeNone Attribute = 0
	// AttributeHasAux has aux
	AttributeHasAux = 1 << iota
	// AttributeHasColor has color
	AttributeHasColor
	// AttributeSolid solid
	AttributeSolid
	// AttributeFluid fluid
	AttributeFluid
)

// HasAux check if block has aux
func (b Block) HasAux() bool { return b.Attributes&AttributeHasAux != 0 }

// HasColor check if block has color
func (b Block) HasColor() bool { return b.Attributes&AttributeHasColor != 0 }

// IsSolid check if block is solid
func (b Block) IsSolid() bool { return b.Attributes&AttributeSolid != 0 }

// IsFluid check if block is fluid
func (b Block) IsFluid() bool { return b.Attributes&AttributeFluid != 0 }

// Instance block instance
type Instance struct {
	ID    BID
	Aux   uint32
	Color uint32
}

// GetBlock get block from instance
func (i *Instance) GetBlock() *Block {
	return Get(i.ID)
}

// ErrInvalidInstance Invalid block instance
var ErrInvalidInstance = errors.New("Invalid block instance")

// Normalize fix aux data
func (i *Instance) Normalize() {
	blk := Get(i.ID)
	if blk == nil {
		panic(ErrInvalidInstance)
	}
	if !blk.HasAux() {
		i.Aux = 0
	}
	if !blk.HasColor() {
		i.Color = 0
	}
}

// ExtraLoader loader for extra info
type ExtraLoader func(in packed.Input) packed.Serializable

// BID block id
type BID uint16

var (
	registry [^BID(0)]Block
	loaders      = make(map[BID]ExtraLoader)
	maxID    BID = 0
	index        = make(map[string]BID)
)

// ErrConflict Conflicted block
var ErrConflict = errors.New("Conflicted block")

// Register register block
func Register(b Block) BID {
	if _, ok := index[b.Name]; ok {
		panic(ErrConflict)
	}
	defer func() { maxID++ }()
	registry[maxID] = b
	index[b.Name] = maxID
	return maxID
}

// Get get block by id
func Get(id BID) *Block {
	if id < maxID {
		return &registry[id]
	}
	return nil
}

// LookupID lookup block id by name
func LookupID(name string) (BID, bool) {
	ret, ok := index[name]
	return ret, ok
}

// Lookup lookup block by name
func Lookup(name string) *Block {
	if ret, ok := index[name]; ok {
		return &registry[ret]
	}
	return nil
}

// GetExtraLoader get extra loader by id
func GetExtraLoader(id BID) ExtraLoader {
	loader, ok := loaders[id]
	if !ok {
		return nil
	}
	return loader
}

// RegisterExtraLoader register extra data loader
func RegisterExtraLoader(id BID, loader ExtraLoader) {
	loaders[id] = loader
}

// DescribeBlocks serialize block
func DescribeBlocks(o packed.Output) {
	o.WriteVarUint32(uint32(maxID))
	for i := BID(0); i < maxID; i++ {
		o.WriteString(registry[i].Name)
		o.WriteUint8(uint8(registry[i].Attributes))
	}
}
