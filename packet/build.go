package packet

import packed "github.com/EndCrystal/PackedIO"

// WritePacket serialize packet
func WritePacket(pkt SendOnlyPacket, out packed.Output) {
	out.WriteUint8(uint8(pkt.PacketID()))
	pkt.Save(out)
}
