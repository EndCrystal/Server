package packet

import packed "github.com/EndCrystal/PackedIO"

// DisconnectPacket server-side disconnect notify
type DisconnectPacket struct {
	Message string
}

// Save save to data
func (pkt DisconnectPacket) Save(out packed.Output) {
	out.WriteString(pkt.Message)
}

// PacketID id
func (pkt DisconnectPacket) PacketID() PID { return IDDisconnect }
