package chunk

import (
	"errors"
	"sort"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/block"
)

type ChunkPos struct{ X, Z int32 }

func (cp *ChunkPos) Load(in packed.Input) {
	cp.X = in.ReadInt32()
	cp.Z = in.ReadInt32()
}

func (cp ChunkPos) Save(out packed.Output) {
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
	} else {
		return b
	}
}

func (cp ChunkPos) Distance(oth ChunkPos) uint32 {
	return max(abs(cp.X-oth.X), abs(cp.Z-oth.Z))
}

type BlockPos uint16

func Pos(x, y, z uint8) BlockPos {
	return BlockPos(x) + 16*(BlockPos(z)+16*BlockPos(y))
}

func (p BlockPos) X() uint8 { return uint8(p) % 16 }
func (p BlockPos) Z() uint8 { return uint8(p/16) % 16 }
func (p BlockPos) Y() uint8 { return uint8(p / 256) }

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
	Blocks     [^BlockPos(0)]block.Instance
	Background [^BlockPos(0)]block.Instance
	ExtraData  map[BlockPos]packed.Serializable
}

var ECorruptedData = errors.New("Failed to load chunk data: corrupted data")
var EUnknownBlock = errors.New("Failed to load chunk data: unknown block")

func (chk *Chunk) Load(in packed.Input) {
	var maps []block.Instance
	in.IterateArray(func(length int) { maps = make([]block.Instance, length) }, func(i int) {
		name := in.ReadString()
		blk, ok := block.LookupId(name)
		if !ok {
			panic(ECorruptedData)
		}
		maps[i].Id = blk
		b := block.Get(blk)
		if b == nil {
			panic(EUnknownBlock)
		}
		if b.HasAux() {
			maps[i].Aux = in.ReadVarUint32()
		}
		if b.HasColor() {
			maps[i].Color = in.ReadUint32()
		}
	})
	lmaps := uint32(len(maps))
	for i := BlockPos(0); i < ^BlockPos(0); i++ {
		idx := in.ReadVarUint32()
		if idx >= lmaps {
			panic(ECorruptedData)
		}
		chk.Blocks[i] = maps[idx]
	}
	for i := BlockPos(0); i < ^BlockPos(0); i++ {
		idx := in.ReadVarUint32()
		if idx >= lmaps {
			panic(ECorruptedData)
		}
		chk.Background[i] = maps[idx]
	}
	in.IterateArray(nil, func(int) {
		pos := BlockPos(in.ReadUint16())
		id := chk.Blocks[pos].Id
		loader := block.GetExtraLoader(id)
		if loader == nil {
			panic(ECorruptedData)
		}
		chk.ExtraData[pos] = loader(in)
	})
}

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
	sort.Slice(list, func(a int, b int) bool { return list[a].count < list[b].count })
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
