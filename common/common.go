package common

import (
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/packet"
)

type PluginCommonHost struct{}

type GameStartHandler func(username string, identify network.NetworkIdentifier) packet.GameStartPacket

var Value struct {
	GameStartHandler GameStartHandler
}

func (PluginCommonHost) SetGameStartHandler(handler GameStartHandler) {
	Value.GameStartHandler = handler
}

func init() {
	Value.GameStartHandler = func(username string, identify network.NetworkIdentifier) packet.GameStartPacket {
		return packet.GameStartPacket{
			Username: username,
			Label:    username,
			Motd:     "EndCrystal",
		}
	}
}
