package main

import (
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/packet"
)

func (g Global) BroadcastPacket(pkt packet.Packet) {
	g.users.Range(func(key interface{}, value interface{}) bool {
		value.(network.ClientInstance).SendPacket(pkt)
		return true
	})
}
