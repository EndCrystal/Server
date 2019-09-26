package packet

import packed "github.com/EndCrystal/PackedIO"

type BatchPacket struct {
	ctx             *ParseContext
	ReceivedPackets []ReceiveOnlyPacket
	SendPackets     []SendOnlyPacket
}

func (pkt *BatchPacket) Load(in packed.Input) {
	in.IterateArray(func(length int) { pkt.ReceivedPackets = make([]ReceiveOnlyPacket, length) }, func(i int) {
		pkt.ReceivedPackets[i] = Parse(in, pkt.ctx)
	})
}

func (pkt BatchPacket) Save(out packed.Output) {
	out.WriteVarUint32(uint32(len(pkt.SendPackets)))
	for _, item := range pkt.SendPackets {
		item.Save(out)
	}
}

func (pkt BatchPacket) PacketId() PacketId            { return IdBatch }
func (pkt BatchPacket) Check(pctx *ParseContext) bool { return pctx.Check(256) }
