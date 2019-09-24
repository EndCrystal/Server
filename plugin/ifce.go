package plug

import (
	"github.com/EndCrystal/Server/common"
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/world/actor"
)

type PluginInterface struct {
	common.PluginCommonHost
	network.PluginNetworkHost
	actor.PluginActorHost
}
