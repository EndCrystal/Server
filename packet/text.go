package packet

import (
	"errors"

	packed "github.com/EndCrystal/PackedIO"
)

type TextPacket struct {
	Flags   TextPacketFlag
	Sender  string
	Payload TextPacketPayload
}

type TextPacketFlag uint8

const (
	TextPacketNormal     TextPacketFlag = 0x00
	TextPacketFromSystem                = 0x01
	TextPacketShowSender                = 0x02
)

func (flag TextPacketFlag) IsFromSystem() bool { return flag&TextPacketFromSystem != 0 }
func (flag TextPacketFlag) IsShowSender() bool { return flag&TextPacketShowSender != 0 }

type TextPacketPayload interface {
	packed.Serializable
	PayloadType() TextPacketPayloadType
}

type TextPacketPayloadType uint8

const (
	TextPacketPlainText TextPacketPayloadType = 0
)

type TextPacketPlainTextPayload struct{ Content string }

func (payload TextPacketPlainTextPayload) PayloadType() TextPacketPayloadType {
	return TextPacketPlainText
}

func (payload *TextPacketPlainTextPayload) Load(in packed.Input) {
	payload.Content = in.ReadString()
}

func (payload TextPacketPlainTextPayload) Save(out packed.Output) {
	out.WriteString(payload.Content)
}

var EInvalidTextPacketPayload = errors.New("Invalid payload of text packet")

func (pkt *TextPacket) Load(in packed.Input) {
	pkt.Flags = TextPacketFlag(in.ReadUint8())
	pkt.Sender = in.ReadString()
	switch TextPacketPayloadType(in.ReadUint8()) {
	case TextPacketPlainText:
		pkt.Payload = new(TextPacketPlainTextPayload)
	default:
		panic(EInvalidTextPacketPayload)
	}
	pkt.Payload.Load(in)
}

func (pkt TextPacket) Save(out packed.Output) {
	out.WriteUint8(uint8(pkt.Flags))
	out.WriteString(pkt.Sender)
	out.WriteUint8(uint8(pkt.Payload.PayloadType()))
	pkt.Payload.Save(out)
}

func (pkt TextPacket) PacketId() PacketId { return IdText }
