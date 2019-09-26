package packet

import packed "github.com/EndCrystal/PackedIO"

type DisconnectPacket struct {
	Message string
}

func (pkt DisconnectPacket) Save(out packed.Output) {
	out.WriteString(pkt.Message)
}

func (pkt DisconnectPacket) PacketId() PacketId { return IdDisconnect }
