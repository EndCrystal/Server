package chunk

import (
	"errors"
	"sort"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/block"
)

// CPos position of chunk
type CPos struct{ X, Z int32 }

// Load load from data
func (cp *CPos) Load(in packed.Input) {
	cp.X = in.ReadInt32()
	cp.Z = in.ReadInt32()
}

// Save save to data
func (cp CPos) Save(out packed.Output) {
	out.WriteInt32(cp.X)
	out.WriteInt32(cp.Z)
}

func abs(x int32) uint32 {
	r := x >> 31
	return uint32((x ^ r) - r)
}

func max(a, b uint32) (r uint32) {
	if a > b {
		return a
	}
	return b
}

// Distance calculate distance between two pos
func (cp CPos) Distance(oth CPos) uint32 {
	return max(abs(cp.X-oth.X), abs(cp.Z-oth.Z))
}

// BPos block pos
type BPos uint16

// Pos make block pos from x y z
func Pos(x, y, z uint8) BPos {
	return BPos(x) + 16*(BPos(z)+16*BPos(y))
}

// X extract x from block pos
func (p BPos) X() uint8 { return uint8(p) % 16 }

// Z extract z from block pos
func (p BPos) Z() uint8 { return uint8(p/16) % 16 }

// Y extract y from block pos
func (p BPos) Y() uint8 { return uint8(p / 256) }

// Chunk chunk data
// Format:
// <Map local id to block.Instance>
// (length: varuint32) * { (block name: string) [aux: varuint32] [color: uint32] }
// <Block data (constant 65536)>
// 65536 * (local id: varuint32)
// <Background Block data (constant 65536)>
// 65536 * (local id: varuint32)
// <Extra data>
// (length: varuint32) * { (position: uint16) (extradata: data) }
type Chunk struct {
	Blocks     [65536]block.Instance
	Background [65536]block.Instance
	ExtraData  map[BPos]packed.Serializable
}

// ErrCorruptedData corrupted data
var ErrCorruptedData = errors.New("Failed to load chunk data: corrupted data")

// ErrUnknownBlock unknown block
var ErrUnknownBlock = errors.New("Failed to load chunk data: unknown block")

// Load load from data
func (chk *Chunk) Load(in packed.Input) {
	var maps []block.Instance
	in.IterateArray(func(length int) { maps = make([]block.Instance, length) }, func(i int) {
		name := in.ReadString()
		blk, ok := block.LookupID(name)
		if !ok {
			panic(ErrCorruptedData)
		}
		maps[i].ID = blk
		b := block.Get(blk)
		if b == nil {
			panic(ErrUnknownBlock)
		}
		if b.HasAux() {
			maps[i].Aux = in.ReadVarUint32()
		}
		if b.HasColor() {
			maps[i].Color = in.ReadUint32()
		}
	})
	lmaps := uint32(len(maps))
	for i := BPos(0); i <= ^BPos(0); i++ {
		idx := in.ReadVarUint32()
		if idx >= lmaps {
			panic(ErrCorruptedData)
		}
		chk.Blocks[i] = maps[idx]
	}
	for i := BPos(0); i <= ^BPos(0); i++ {
		idx := in.ReadVarUint32()
		if idx >= lmaps {
			panic(ErrCorruptedData)
		}
		chk.Background[i] = maps[idx]
	}
	in.IterateArray(nil, func(int) {
		pos := BPos(in.ReadUint16())
		id := chk.Blocks[pos].ID
		loader := block.GetExtraLoader(id)
		if loader == nil {
			panic(ErrCorruptedData)
		}
		chk.ExtraData[pos] = loader(in)
	})
}

// Save save to data
func (chk Chunk) Save(out packed.Output) {
	maps := make(map[block.Instance]int)
	for _, ins := range chk.Blocks {
		maps[ins]++
	}
	for _, ins := range chk.Background {
		maps[ins]++
	}
	num := len(maps)
	list := make([]struct {
		ins   block.Instance
		count int
	}, num)
	i := 0
	for ins, count := range maps {
		list[i].ins = ins
		list[i].count = count
		i++
	}
	sort.Slice(list, func(a int, b int) bool { return list[a].count > list[b].count })
	rmap := make(map[block.Instance]int)
	out.WriteVarUint32(uint32(num))
	for i, tpl := range list {
		ins := tpl.ins
		blk := ins.GetBlock()
		out.WriteString(blk.Name)
		if blk.HasAux() {
			out.WriteVarUint32(ins.Aux)
		}
		if blk.HasColor() {
			out.WriteUint32(ins.Color)
		}
		rmap[ins] = i
	}
	for _, ins := range chk.Blocks {
		idx := rmap[ins]
		out.WriteVarUint32(uint32(idx))
	}
	for _, ins := range chk.Background {
		idx := rmap[ins]
		out.WriteVarUint32(uint32(idx))
	}
	out.WriteVarUint32(uint32(len(chk.ExtraData)))
	for _, data := range chk.ExtraData {
		data.Save(out)
	}
}
