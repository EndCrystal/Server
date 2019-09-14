package main

import (
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/packet"
)

func handleChat(broadcaster network.PacketBroadcaster) chan<- ChatMessage {
	ret := make(chan ChatMessage)
	go func() {
		for msg := range ret {
			broadcaster.BroadcastPacket(&packet.TextPacket{
				packet.TextPacketNormal,
				msg.Sender,
				&packet.TextPacketPlainTextPayload{msg.Message},
			})
		}
	}()
	return ret
}
