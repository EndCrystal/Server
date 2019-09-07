package plug

import (
	"github.com/EndCrystal/Server/network"
)

type PluginInterface interface {
	RegisterNetworkProtocol(name string, fn network.ServerCreator)
}
