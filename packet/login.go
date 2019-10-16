package packet

import (
	"errors"
	"time"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/token"
)

// LoginPayload payload for login
type LoginPayload struct {
	ServerID string
	Username string
	Time     time.Time
}

// LoginPacket login packet
type LoginPacket struct {
	token token.Token
	raw   []byte
}

// ErrNetworkVersionMismatched Network version mismatched
var ErrNetworkVersionMismatched = errors.New("Network version mismatched")

// Load load from data
func (pkt *LoginPacket) Load(in packed.Input) {
	nv := in.ReadVarUint32()
	if nv != NetworkVersion {
		panic(ErrNetworkVersionMismatched)
	}
	in.ReadFixedBytes(pkt.token[:])
	pkt.raw = in.ReadBytes()
}

// Save save to data
func (pkt LoginPacket) Save(out packed.Output) {
	out.WriteVarUint32(NetworkVersion)
	out.WriteFixedBytes(pkt.token[:])
	out.WriteBytes(pkt.raw)
}

func (pkt *LoginPacket) Write(payload LoginPayload, gen token.Generator) {
	out, buf := packed.NewOutput()
	out.WriteString(payload.ServerID)
	out.WriteString(payload.Username)
	out.WriteInt64(payload.Time.Unix())
	pkt.raw = buf.Bytes()
	gen(pkt.raw, &pkt.token)
}

// Verify verify data
func (pkt LoginPacket) Verify(verifier token.Verifier) bool {
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
	payload.ServerID = in.ReadString()
	payload.Username = in.ReadString()
	payload.Time = time.Unix(in.ReadInt64(), 0)
	return
}

// PacketID id
func (pkt LoginPacket) PacketID() PID { return IDLogin }

// Check check for parsing
func (pkt LoginPacket) Check(pctx *ParseContext) bool {
	return pctx.Check(1024 + uint16(len(pkt.raw)))
}
