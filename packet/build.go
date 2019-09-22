package packet

import packed "github.com/EndCrystal/PackedIO"

func BuildPacket(pkt Packet, out packed.Output) {
	out.WriteUint8(uint8(pkt.PacketId()))
	pkt.Save(out)
}
