package packet

import packed "github.com/EndCrystal/PackedIO"

type ChatPacket struct {
	Message string
}

func (pkt *ChatPacket) Load(in packed.Input) {
	pkt.Message = in.ReadString()
}

func (pkt ChatPacket) Save(out packed.Output) {
	out.WriteString(pkt.Message)
}

func (pkt ChatPacket) PacketId() PacketId { return IdChat }
func (pkt ChatPacket) Check(pctx *ParseContext) bool {
	l := len(pkt.Message)
	if l > 1024 {
		return false
	}
	return pctx.Check(uint16(l * 2))
}
