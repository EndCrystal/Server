package plug

import (
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/world/actor"
)

type PluginInterface struct {
	network.PluginNetworkHost
	actor.PluginActorHost
}
