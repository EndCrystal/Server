package components

import (
	"github.com/EndCrystal/Server/user"
)

// UserControl control component
type UserControl struct{ Owner *user.Info }

// UserControlID id
var UserControlID = RegisterRuntime("core:user_control", (*UserControl)(nil), nil)
