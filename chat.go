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
				Flags:   packet.TextPacketNormal,
				Sender:  msg.Sender,
				Payload: &packet.TextPacketPlainTextPayload{Content: msg.Message},
			})
		}
	}()
	return ret
}
