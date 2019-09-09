package builtin

import packed "github.com/EndCrystal/PackedIO"

type Position struct {
	X, Z float64
	Y    float32
}

func (p *Position) GetPosition() *Position { return p }
func (p *Position) LoadPosition(in packed.Input) {
	p.X = in.ReadFloat64()
	p.Y = in.ReadFloat32()
	p.Z = in.ReadFloat64()
}
func (p *Position) SavePosition(out packed.Output) {
	out.WriteFloat64(p.X)
	out.WriteFloat32(p.Y)
	out.WriteFloat64(p.Z)
}

type PositionComponent interface {
	GetPosition() *Position
	LoadPosition(packed.Input)
	SavePosition(packed.Output)
}
