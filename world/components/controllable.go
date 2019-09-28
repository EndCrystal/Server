package components

import (
	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/user"
)

type (
	ControlRequest struct {
		Source *user.UserInfo
	}
	Controllable struct {
		ControlRequest chan ControlRequest
	}
	ControllableComponent interface{ Controllable() Controllable }
	controllableInfo      struct{}
)

func (c Controllable) Controllable() Controllable { return c }

func (controllableInfo) Secure() bool { return false }
func (controllableInfo) CreateComponent(i packed.Input) interface{} {
	return Controllable{make(chan ControlRequest)}
}
func (controllableInfo) SaveComponent(o packed.Output, obj interface{}) {}

var ControllableId = RegisterRuntime("core:controllable", (*ControllableComponent)(nil), controllableInfo{})
