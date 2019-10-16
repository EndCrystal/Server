package main

import (
	"github.com/EndCrystal/Server/packet"
)

// BroadcastPacket Broadcast packet globally
func (g serverGlobal) BroadcastPacket(pkt packet.SendOnlyPacket) {
	g.users.Range(func(key interface{}, value interface{}) bool {
		value.(userInfoWithClient).SendPacket(pkt)
		return true
	})
}
