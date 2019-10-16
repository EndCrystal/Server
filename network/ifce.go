package network

import (
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/EndCrystal/Server/packet"
)

// PacketSender send packet
type PacketSender interface {
	SendPacket(pkt packet.SendOnlyPacket) error
}

// PacketBroadcaster broadcast packet
type PacketBroadcaster interface {
	BroadcastPacket(pkt packet.SendOnlyPacket)
}

// Identifier represents network identifier
type Identifier interface {
	String() string
	GetIP() (net.IP, bool)
	GetPort() (uint16, bool)
}

// CommonIdentifier common ip identifier
type CommonIdentifier struct {
	IP   net.IP
	Port uint16
}

func (id CommonIdentifier) String() string {
	return fmt.Sprintf("%v:%d", id.IP, id.Port)
}

// GetIP get ip of identifier
func (id CommonIdentifier) GetIP() (net.IP, bool) { return id.IP, true }

// GetPort get port of identifir
func (id CommonIdentifier) GetPort() (uint16, bool) { return id.Port, true }

// ClientInstance represent Client instnace
type ClientInstance interface {
	PacketSender
	GetIdentifier() Identifier
	GetFetcher() <-chan packet.ReceiveOnlyPacket
	Disconnect()
}

// Server represent Server instance
type Server interface {
	GetFetcher() <-chan ClientInstance
	Stop()
}

// ErrInvalidScheme Unregistered scheme
var ErrInvalidScheme = errors.New("Invalid scheme")

// CreateServer create server for url
func CreateServer(u *url.URL) (Server, error) {
	if creator, ok := registry[u.Scheme]; ok {
		return creator(u)
	}
	return nil, ErrInvalidScheme
}

var registry = make(map[string]ServerCreator)

// ServerCreator create server function
type ServerCreator func(*url.URL) (Server, error)

// PluginNetworkHost plugin host
type PluginNetworkHost struct{}

// RegisterNetworkProtocol register network protocol
func (PluginNetworkHost) RegisterNetworkProtocol(name string, fn ServerCreator) {
	registry[name] = fn
}
