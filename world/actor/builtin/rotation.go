package builtin

import packed "github.com/EndCrystal/PackedIO"

type (
	Rotation          struct{ RotationAngle float32 }
	RotationComponent interface {
		GetRotation() *Rotation
		LoadRotation(in packed.Input)
		SaveRotation(out packed.Output)
	}
)

func (rot *Rotation) GetRotation() *Rotation        { return rot }
func (rot *Rotation) LoadRotation(in packed.Input)  { rot.RotationAngle = in.ReadFloat32() }
func (rot Rotation) SaveRotation(out packed.Output) { out.WriteFloat32(rot.RotationAngle) }
