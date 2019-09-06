package packet

import packed "github.com/EndCrystal/PackedIO"

type BatchPacket struct {
	ctx     *ParseContext
	Packets []Packet
}

func (pkt *BatchPacket) Load(in packed.Input) {
	in.IterateArray(func(length int) { pkt.Packets = make([]Packet, length) }, func(i int) {
		pkt.Packets[i] = Parse(in, pkt.ctx)
	})
}

func (pkt BatchPacket) Save(out packed.Output) {
	out.WriteVarUint32(uint32(len(pkt.Packets)))
	for _, item := range pkt.Packets {
		item.Save(out)
	}
}

func (pkt BatchPacket) PacketId() PacketId           { return IdBatch }
func (pkt BatchPacket) Check(ctx *ParseContext) bool { return true }
