package main

import (
	"github.com/EndCrystal/Server/packet"
)

func (g Global) BroadcastPacket(pkt packet.SendOnlyPacket) {
	g.users.Range(func(key interface{}, value interface{}) bool {
		value.(UserInfoWithClient).SendPacket(pkt)
		return true
	})
}
