package packet

import (
	"errors"
	"fmt"

	packed "github.com/EndCrystal/PackedIO"
)

// ErrUnknownPacket Unknown Packet ID
var ErrUnknownPacket = errors.New("Unknown Packet ID")

// ErrDisallowedPacket Disallowed Packet
var ErrDisallowedPacket = errors.New("Disallowed Packet")

// Parse parse packet from data
func Parse(in packed.Input, ctx *ParseContext) (pkt ReceiveOnlyPacket) {
	id := PID(in.ReadUint8())
	switch id {
	case IDBatch:
		pkt = &BatchPacket{ctx, nil, nil}
	case IDLogin:
		pkt = new(LoginPacket)
	case IDChat:
		pkt = new(ChatPacket)
	case IDChunkRequest:
		pkt = new(ChunkRequestPacket)
	default:
		panic(ErrUnknownPacket)
	}
	pkt.Load(in)
	if !pkt.Check(ctx) {
		panic(ErrDisallowedPacket)
	}
	return
}

// ParsePacket safe parse
func ParsePacket(in packed.Input, quota uint16) (pkt ReceiveOnlyPacket, err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			if err, ok = e.(error); !ok {
				err = fmt.Errorf("%v", e)
			}
			pkt = nil
			return
		}
	}()
	pkt = Parse(in, &ParseContext{quota})
	return
}
