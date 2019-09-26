package packet

import (
	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/components"
)

type GameStartPacket struct {
	Username string
	Label    string
	Motd     string
}

func (pkt GameStartPacket) Save(out packed.Output) {
	out.WriteString(pkt.Username)
	out.WriteString(pkt.Label)
	out.WriteString(pkt.Motd)
	components.DescribeComponents(out)
}

func (pkt GameStartPacket) PacketId() PacketId { return IdGameStart }
