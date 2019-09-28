package components

import (
	"github.com/EndCrystal/Server/user"
)

type UserControl struct{ Owner *user.UserInfo }

var UserControlId = RegisterRuntime("core:user_control", (*UserControl)(nil), nil)
