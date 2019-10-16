package packet

import (
	packed "github.com/EndCrystal/PackedIO"
)

// SendOnlyPacket server-side packet
type SendOnlyPacket interface {
	PacketID() PID
	Save(packed.Output)
}

// ReceiveOnlyPacket client-side packet
type ReceiveOnlyPacket interface {
	PacketID() PID
	Load(packed.Input)
	Check(ctx *ParseContext) bool
}

// Packet packet interface
type Packet interface {
	packed.Serializable
	PacketID() PID
	Check(ctx *ParseContext) bool
}

// NetworkVersion version code
const NetworkVersion uint32 = 0x01

// PID unique id for packet
type PID uint8

const (
	// IDBatch batch
	IDBatch PID = iota
	// IDLogin login
	IDLogin
	// IDDisconnect disconnect
	IDDisconnect
	// IDGameStart game start
	IDGameStart
	// IDChat chat message
	IDChat
	// IDText general text packet
	IDText
	// IDChunkRequest chunk request
	IDChunkRequest
	// IDChunkData chunk data
	IDChunkData
)

// ParseContext context for parsing
type ParseContext struct {
	Quota uint16
}

// Check check context for parsing
func (ctx *ParseContext) Check(eat uint16) bool {
	if eat > ctx.Quota {
		return false
	}
	ctx.Quota -= eat
	return true
}
