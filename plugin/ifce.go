package plug

import (
	"github.com/EndCrystal/Server/common"
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/world/block"
	"github.com/EndCrystal/Server/world/components"
	"github.com/EndCrystal/Server/world/dim"
)

type PluginInterface struct {
	common.PluginCommonHost
	network.PluginNetworkHost
	dim.PluginDimensionHost
	block.PluginBlockHost
	components.PluginComponentsHost
}
