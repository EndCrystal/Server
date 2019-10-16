package packet

import (
	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/block"
	"github.com/EndCrystal/Server/world/chunk"
	"github.com/EndCrystal/Server/world/components"
)

// GameStartPacket game start packet
type GameStartPacket struct {
	Username        string
	Label           string
	Motd            string
	MaxViewDistance uint32
	InitialPosition chunk.CPos
}

// Save save to data
func (pkt GameStartPacket) Save(out packed.Output) {
	out.WriteString(pkt.Username)
	out.WriteString(pkt.Label)
	out.WriteString(pkt.Motd)
	out.WriteVarUint32(pkt.MaxViewDistance)
	pkt.InitialPosition.Save(out)
	block.DescribeBlocks(out)
	components.DescribeComponents(out)
}

// PacketID id
func (pkt GameStartPacket) PacketID() PID { return IDGameStart }
