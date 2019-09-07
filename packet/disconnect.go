package packet

import packed "github.com/EndCrystal/PackedIO"

type DisconnectPacket struct {
	Message string
}

func (pkt *DisconnectPacket) Load(in packed.Input) {
	pkt.Message = in.ReadString()
}

func (pkt DisconnectPacket) Save(out packed.Output) {
	out.WriteString(pkt.Message)
}

func (pkt DisconnectPacket) PacketId() PacketId            { return IdDisconnect }
func (pkt DisconnectPacket) Check(pctx *ParseContext) bool { return pctx.Check(ServerSide, 0) }
