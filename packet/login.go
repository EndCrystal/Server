package packet

import (
	"time"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/token"
)

type LoginPayload struct {
	Username string
	Time     time.Time
}

type LoginPacket struct {
	token token.Token
	raw   []byte
}

func (pkt *LoginPacket) Load(in packed.Input) {
	in.ReadFixedBytes(pkt.token[:])
	pkt.raw = in.ReadBytes()
}

func (pkt LoginPacket) Save(out packed.Output) {
	out.WriteFixedBytes(pkt.token[:])
	out.WriteBytes(pkt.raw)
}

func (pkt *LoginPacket) Write(payload LoginPayload, gen token.TokenGenerator) {
	out, buf := packed.NewOutput()
	out.WriteString(payload.Username)
	out.WriteInt64(payload.Time.Unix())
	pkt.raw = buf.Bytes()
	gen(pkt.raw, pkt.token)
}

func (pkt LoginPacket) Verify(verifier token.TokenVerifier) bool {
	return verifier(pkt.raw, pkt.token)
}

func (pkt LoginPacket) Read() (payload LoginPayload, ok bool) {
	ok = true
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	in := packed.InputFromBuffer(pkt.raw)
	payload.Username = in.ReadString()
	payload.Time = time.Unix(in.ReadInt64(), 0)
	return
}

func (pkt LoginPacket) PacketId() PacketId { return IdLogin }

func (pkt LoginPacket) Check(pctx *ParseContext) bool { return pctx.Check(ClientSide, 1024) }
