package network

import (
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/EndCrystal/Server/packet"
)

type PacketSender interface {
	SendPacket(pkt packet.SendOnlyPacket) error
}

type PacketBroadcaster interface {
	BroadcastPacket(pkt packet.SendOnlyPacket)
}

type NetworkIdentifier interface {
	String() string
	GetIP() (net.IP, bool)
	GetPort() (uint16, bool)
}

type CommonNetworkIdentifier struct {
	IP   net.IP
	Port uint16
}

func (id CommonNetworkIdentifier) String() string {
	return fmt.Sprintf("%v:%d", id.IP, id.Port)
}

func (id CommonNetworkIdentifier) GetIP() (net.IP, bool)   { return id.IP, true }
func (id CommonNetworkIdentifier) GetPort() (uint16, bool) { return id.Port, true }

type ClientInstance interface {
	PacketSender
	GetIdentifier() NetworkIdentifier
	GetFetcher() <-chan packet.ReceiveOnlyPacket
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

var registry = make(map[string]ServerCreator)

type ServerCreator func(*url.URL) (Server, error)

type PluginNetworkHost struct{}

func (PluginNetworkHost) RegisterNetworkProtocol(name string, fn ServerCreator) {
	registry[name] = fn
}
