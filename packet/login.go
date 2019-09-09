package packet

import (
	"errors"
	"time"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/token"
	"github.com/EndCrystal/Server/world/actor/builtin"
)

type LoginPayload struct {
	ServerId string
	Username string
	Time     time.Time
	builtin.Position
	builtin.Rotation
}

type LoginPacket struct {
	token token.Token
	raw   []byte
}

var ENetworkVersionMismatched = errors.New("Network version mismatched")

func (pkt *LoginPacket) Load(in packed.Input) {
	nv := in.ReadVarUint32()
	if nv != NetworkVersion {
		panic(ENetworkVersionMismatched)
	}
	in.ReadFixedBytes(pkt.token[:])
	pkt.raw = in.ReadBytes()
}

func (pkt LoginPacket) Save(out packed.Output) {
	out.WriteVarUint32(NetworkVersion)
	out.WriteFixedBytes(pkt.token[:])
	out.WriteBytes(pkt.raw)
}

func (pkt *LoginPacket) Write(payload LoginPayload, gen token.TokenGenerator) {
	out, buf := packed.NewOutput()
	out.WriteString(payload.ServerId)
	out.WriteString(payload.Username)
	out.WriteInt64(payload.Time.Unix())
	payload.SavePosition(out)
	payload.SaveRotation(out)
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
	payload.ServerId = in.ReadString()
	payload.Username = in.ReadString()
	payload.Time = time.Unix(in.ReadInt64(), 0)
	payload.LoadPosition(in)
	payload.LoadRotation(in)
	return
}

func (pkt LoginPacket) PacketId() PacketId { return IdLogin }

func (pkt LoginPacket) Check(pctx *ParseContext) bool { return pctx.Check(ClientSide, 1024) }
