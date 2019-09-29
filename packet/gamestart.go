package packet

import (
	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/block"
	"github.com/EndCrystal/Server/world/chunk"
	"github.com/EndCrystal/Server/world/components"
)

type GameStartPacket struct {
	Username        string
	Label           string
	Motd            string
	InitialPosition chunk.ChunkPos
}

func (pkt GameStartPacket) Save(out packed.Output) {
	out.WriteString(pkt.Username)
	out.WriteString(pkt.Label)
	out.WriteString(pkt.Motd)
	pkt.InitialPosition.Save(out)
	block.DescribeBlocks(out)
	components.DescribeComponents(out)
}

func (pkt GameStartPacket) PacketId() PacketId { return IdGameStart }
