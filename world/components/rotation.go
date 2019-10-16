package components

import packed "github.com/EndCrystal/PackedIO"

type (
	// Rotation component
	Rotation          struct{ RotationAngle float32 }
	// RotationComponent component interface
	RotationComponent interface {
		GetRotation() *Rotation
		LoadRotation(in packed.Input)
		SaveRotation(out packed.Output)
	}
)

// GetRotation get component
func (rot *Rotation) GetRotation() *Rotation        { return rot }
// LoadRotation load component
func (rot *Rotation) LoadRotation(in packed.Input)  { rot.RotationAngle = in.ReadFloat32() }
// SaveRotation save component
func (rot Rotation) SaveRotation(out packed.Output) { out.WriteFloat32(rot.RotationAngle) }

type rotationInfo struct{}

func (rotationInfo) Secure() bool { return false }
func (rotationInfo) LoadComponent(i packed.Input, obj interface{}) {
	obj.(RotationComponent).LoadRotation(i)
}
func (rotationInfo) SaveComponent(o packed.Output, obj interface{}) {
	obj.(RotationComponent).SaveRotation(o)
}

// RotationID id
var RotationID = Register("core:rotation", (*RotationComponent)(nil), rotationInfo{})
