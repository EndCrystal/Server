package main

import (
	"strings"

	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/packet"
)

func handleChat(broadcaster network.PacketBroadcaster) chan<- ChatMessage {
	log := logprefix.Get("[chat system] ")
	ret := make(chan ChatMessage)
	go func() {
		for msg := range ret {
			message := strings.Trim(msg.Message, " ")
			if len(message) == 0 {
				continue
			}
			log.Printf("<%s> %s", msg.Sender, msg.Message)
			broadcaster.BroadcastPacket(&packet.TextPacket{
				Flags:   packet.TextPacketNormal,
				Sender:  msg.Sender,
				Payload: &packet.TextPacketPlainTextPayload{Content: message},
			})
		}
	}()
	return ret
}
