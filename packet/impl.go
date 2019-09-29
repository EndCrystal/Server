package packet

import (
	packed "github.com/EndCrystal/PackedIO"
)

type SendOnlyPacket interface {
	PacketId() PacketId
	Save(packed.Output)
}

type ReceiveOnlyPacket interface {
	PacketId() PacketId
	Load(packed.Input)
	Check(ctx *ParseContext) bool
}

type Packet interface {
	packed.Serializable
	PacketId() PacketId
	Check(ctx *ParseContext) bool
}

const NetworkVersion uint32 = 0x01

type PacketId uint8

const (
	IdBatch PacketId = iota
	IdLogin
	IdDisconnect
	IdGameStart
	IdChat
	IdText
	IdChunkRequest
	IdChunkData
)

type ParseContext struct {
	Quota uint16
}

func (ctx *ParseContext) Check(eat uint16) bool {
	if eat > ctx.Quota {
		return false
	}
	ctx.Quota -= eat
	return true
}
