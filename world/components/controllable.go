package components

import (
	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/user"
)

type (
	// ControlRequest control request
	ControlRequest struct {
		Source *user.Info
	}
	// Controllable controllable
	Controllable struct {
		ControlRequest chan ControlRequest
	}
	// ControllableComponent controllable component
	ControllableComponent interface{ Controllable() Controllable }
	controllableInfo      struct{}
)

// Controllable return Controllable
func (c Controllable) Controllable() Controllable { return c }

func (controllableInfo) Secure() bool { return false }
func (controllableInfo) CreateComponent(i packed.Input) interface{} {
	return Controllable{make(chan ControlRequest)}
}
func (controllableInfo) SaveComponent(o packed.Output, obj interface{}) {}

// ControllableID id
var ControllableID = RegisterRuntime("core:controllable", (*ControllableComponent)(nil), controllableInfo{})
