package packet

import (
	"errors"

	packed "github.com/EndCrystal/PackedIO"
)

// TextPacket text content
type TextPacket struct {
	Flags   TextPacketFlag
	Sender  string
	Payload TextPacketPayload
}

// TextPacketFlag flag fot text content
type TextPacketFlag uint8

const (
	// TextPacketNormal normal text packet
	TextPacketNormal TextPacketFlag = 0x00
	// TextPacketFromSystem text packet from system
	TextPacketFromSystem = 0x01
	// TextPacketShowSender text packet has sender
	TextPacketShowSender = 0x02
)

// IsFromSystem check if sent from system
func (flag TextPacketFlag) IsFromSystem() bool { return flag&TextPacketFromSystem != 0 }

// IsShowSender check if show sender
func (flag TextPacketFlag) IsShowSender() bool { return flag&TextPacketShowSender != 0 }

// TextPacketPayload payload for text packet
type TextPacketPayload interface {
	packed.Serializable
	PayloadType() TextPacketPayloadType
}

// TextPacketPayloadType payload type for text packet
type TextPacketPayloadType uint8

const (
	// TextPacketPlainText plain text packet
	TextPacketPlainText TextPacketPayloadType = 0
)

// TextPacketPlainTextPayload plain text payload for text packet
type TextPacketPlainTextPayload struct{ Content string }

// PayloadType type
func (payload TextPacketPlainTextPayload) PayloadType() TextPacketPayloadType {
	return TextPacketPlainText
}

// Load load payload from data
func (payload *TextPacketPlainTextPayload) Load(in packed.Input) {
	payload.Content = in.ReadString()
}

// Save save payload to data
func (payload TextPacketPlainTextPayload) Save(out packed.Output) {
	out.WriteString(payload.Content)
}

// ErrInvalidTextPacketPayload Invalid payload of text packet
var ErrInvalidTextPacketPayload = errors.New("Invalid payload of text packet")

// Load load from data
func (pkt *TextPacket) Load(in packed.Input) {
	pkt.Flags = TextPacketFlag(in.ReadUint8())
	pkt.Sender = in.ReadString()
	switch TextPacketPayloadType(in.ReadUint8()) {
	case TextPacketPlainText:
		pkt.Payload = new(TextPacketPlainTextPayload)
	default:
		panic(ErrInvalidTextPacketPayload)
	}
	pkt.Payload.Load(in)
}

// Save save to data
func (pkt TextPacket) Save(out packed.Output) {
	out.WriteUint8(uint8(pkt.Flags))
	out.WriteString(pkt.Sender)
	out.WriteUint8(uint8(pkt.Payload.PayloadType()))
	pkt.Payload.Save(out)
}

// PacketID id
func (pkt TextPacket) PacketID() PID { return IDText }
