package common

import "github.com/EndCrystal/Server/network"

// PluginCommonHost plugin host
type PluginCommonHost struct{}

// UserLabelHandler user label handler
type UserLabelHandler func(username string, id network.Identifier) string
// ServerMotdHandler server motd handler
type ServerMotdHandler func(username string, id network.Identifier) string

// Value value
var Value struct {
	UserLabelHandler
	ServerMotdHandler
}

// SetUserLabelHandler set user label handler
func (PluginCommonHost) SetUserLabelHandler(handler UserLabelHandler) {
	Value.UserLabelHandler = handler
}

// SetServerMotdHandler server motd handler
func (PluginCommonHost) SetServerMotdHandler(handler ServerMotdHandler) {
	Value.ServerMotdHandler = handler
}

func init() {
	Value.UserLabelHandler = func(username string, id network.Identifier) string { return username }
	Value.ServerMotdHandler = func(username string, id network.Identifier) string { return "EndCrystal" }
}
