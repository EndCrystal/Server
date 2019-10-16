package packet

import packed "github.com/EndCrystal/PackedIO"

// ChatPacket client-side chat packet
type ChatPacket struct {
	Message string
}

// Load load from data
func (pkt *ChatPacket) Load(in packed.Input) {
	pkt.Message = in.ReadString()
}

// Save save to binary data
func (pkt ChatPacket) Save(out packed.Output) {
	out.WriteString(pkt.Message)
}

// PacketID id
func (pkt ChatPacket) PacketID() PID { return IDChat }

// Check for parse
func (pkt ChatPacket) Check(pctx *ParseContext) bool {
	l := len(pkt.Message)
	if l > 1024 {
		return false
	}
	return pctx.Check(uint16(l * 2))
}
