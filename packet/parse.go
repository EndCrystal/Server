package packet

import (
	"errors"
	"fmt"

	packed "github.com/EndCrystal/PackedIO"
)

var EUnknownPacket = errors.New("Unknown Packet ID")
var EDisallowedPacket = errors.New("Disallowed Packet")

func Parse(in packed.Input, ctx *ParseContext) (pkt ReceiveOnlyPacket) {
	id := PacketId(in.ReadUint8())
	switch id {
	case IdBatch:
		pkt = &BatchPacket{ctx, nil, nil}
	case IdLogin:
		pkt = new(LoginPacket)
	case IdChat:
		pkt = new(ChatPacket)
	default:
		panic(EUnknownPacket)
	}
	pkt.Load(in)
	if !pkt.Check(ctx) {
		panic(EDisallowedPacket)
	}
	return
}

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
