package packet

import (
	packed "github.com/EndCrystal/PackedIO"
)

type Side uint8

const (
	InternalSide = 0
	ServerSide   = 1
	ClientSide   = 2
)

type Packet interface {
	packed.Serializable
	PacketId() PacketId
	Check(ctx *ParseContext) bool
}

type PacketId uint8

const (
	IdBatch = 0xFF
	IdLogin = 0x01
)

type ParseContext struct {
	Side  Side
	Quota uint16
}

func (ctx *ParseContext) Check(side Side, eat uint16) bool {
	if side != 0 && ctx.Side != 0 {
		if side != ctx.Side {
			return false
		}
	}
	if eat > ctx.Quota {
		return false
	}
	ctx.Quota -= eat
	return true
}