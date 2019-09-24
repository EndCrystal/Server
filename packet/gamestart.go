package packet

import packed "github.com/EndCrystal/PackedIO"

type GameStartPacket struct {
	Username string
	Label    string
	Motd     string
}

func (pkt *GameStartPacket) Load(in packed.Input) {
	pkt.Username = in.ReadString()
	pkt.Label = in.ReadString()
	pkt.Motd = in.ReadString()
}

func (pkt GameStartPacket) Save(out packed.Output) {
	out.WriteString(pkt.Username)
	out.WriteString(pkt.Label)
	out.WriteString(pkt.Motd)
}

func (pkt GameStartPacket) PacketId() PacketId           { return IdGameStart }
func (pkt GameStartPacket) Check(ctx *ParseContext) bool { return ctx.Check(ServerSide, 0) }
