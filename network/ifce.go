package network

import (
	"errors"
	"net/url"

	"github.com/EndCrystal/Server/packet"
)

type PacketSender interface {
	SendPacket(pkt packet.Packet) error
}

type PacketBroadcaster interface {
	BroadcastPacket(pkt packet.Packet)
}

type ClientInstance interface {
	PacketSender
	GetFetcher() <-chan packet.Packet
	Disconnect()
}

type Server interface {
	GetFetcher() <-chan ClientInstance
	Stop()
}

var EInvalidScheme = errors.New("Invalid scheme")

func CreateServer(u *url.URL) (Server, error) {
	if creator, ok := registry[u.Scheme]; ok {
		return creator(u)
	}
	return nil, EInvalidScheme
}

var registry map[string]ServerCreator

type ServerCreator func(*url.URL) (Server, error)

type PluginHost struct{}

func (PluginHost) RegisterNetworkProtocol(name string, fn ServerCreator) {
	registry[name] = fn
}
