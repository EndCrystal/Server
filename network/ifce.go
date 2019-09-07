package network

import "github.com/EndCrystal/Server/packet"

type PacketSender interface {
	SendPacket(pkt packet.Packet)
}

type PacketBroadcaster interface {
	BroadcastPacket(pkt packet.Packet)
}

type ClientInstance interface {
	PacketSender
	GetFetcher() <-chan packet.Packet
}

type Server interface {
	PacketBroadcaster
	GetFetcher() <-chan ClientInstance
}

var Registry map[string]func() Server
