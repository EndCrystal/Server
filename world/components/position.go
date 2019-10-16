package components

import packed "github.com/EndCrystal/PackedIO"

// Position position component
type Position struct {
	X, Z float64
	Y    float32
}

// GetPosition get component
func (p *Position) GetPosition() *Position { return p }
// LoadPosition load position
func (p *Position) LoadPosition(in packed.Input) {
	p.X = in.ReadFloat64()
	p.Y = in.ReadFloat32()
	p.Z = in.ReadFloat64()
}
// SavePosition save position
func (p *Position) SavePosition(out packed.Output) {
	out.WriteFloat64(p.X)
	out.WriteFloat32(p.Y)
	out.WriteFloat64(p.Z)
}

// PositionComponent component interface
type PositionComponent interface {
	GetPosition() *Position
	LoadPosition(packed.Input)
	SavePosition(packed.Output)
}

type positionInfo struct{}

func (positionInfo) Secure() bool { return false }
func (positionInfo) LoadComponent(i packed.Input, obj interface{}) {
	obj.(PositionComponent).LoadPosition(i)
}
func (positionInfo) SaveComponent(o packed.Output, obj interface{}) {
	obj.(PositionComponent).SavePosition(o)
}

// PositionID id
var PositionID = Register("core:position", (*PositionComponent)(nil), positionInfo{})
