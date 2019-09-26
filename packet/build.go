package packet

import packed "github.com/EndCrystal/PackedIO"

func WritePacket(pkt SendOnlyPacket, out packed.Output) {
	out.WriteUint8(uint8(pkt.PacketId()))
	pkt.Save(out)
}
