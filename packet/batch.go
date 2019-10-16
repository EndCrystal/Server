package packet

import packed "github.com/EndCrystal/PackedIO"

// BatchPacket represent packet batch
type BatchPacket struct {
	ctx             *ParseContext
	ReceivedPackets []ReceiveOnlyPacket
	SendPackets     []SendOnlyPacket
}

// Load load from data
func (pkt *BatchPacket) Load(in packed.Input) {
	in.IterateArray(func(length int) { pkt.ReceivedPackets = make([]ReceiveOnlyPacket, length) }, func(i int) {
		pkt.ReceivedPackets[i] = Parse(in, pkt.ctx)
	})
}

// Save save as binary data
func (pkt BatchPacket) Save(out packed.Output) {
	out.WriteVarUint32(uint32(len(pkt.SendPackets)))
	for _, item := range pkt.SendPackets {
		item.Save(out)
	}
}

// PacketID id
func (pkt BatchPacket) PacketID() PID { return IDBatch }

// Check check for parse
func (pkt BatchPacket) Check(pctx *ParseContext) bool { return pctx.Check(256) }
