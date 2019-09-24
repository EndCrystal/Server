package common

import "github.com/EndCrystal/Server/network"

type PluginCommonHost struct{}

type UserLabelHandler func(username string, id network.NetworkIdentifier) string
type ServerMotdHandler func(username string, id network.NetworkIdentifier) string

var Value struct {
	UserLabelHandler
	ServerMotdHandler
}

func (PluginCommonHost) SetUserLabelHandler(handler UserLabelHandler) {
	Value.UserLabelHandler = handler
}

func (PluginCommonHost) SetServerMotdHandler(handler ServerMotdHandler) {
	Value.ServerMotdHandler = handler
}

func init() {
	Value.UserLabelHandler = func(username string, id network.NetworkIdentifier) string { return username }
	Value.ServerMotdHandler = func(username string, id network.NetworkIdentifier) string { return "EndCrystal" }
}
