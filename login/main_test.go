package main

import (
	"testing"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/packet"
	"github.com/EndCrystal/Server/token"
)

var verifier token.TokenVerifier

func TestVerify(t *testing.T) {
	pub, priv := token.GenerateKeys()
	generator = token.GetTokenGenerator(priv)
	verifier = token.GetTokenVerifier(pub)

	pkt := genpacket("default", "codehz")
	if !pkt.Verify(verifier) {
		t.Fatalf("Failed to verify in stage 1")
	}

	out, buf := packed.NewOutput()
	packet.WritePacket(&pkt, out)

	in := packed.InputFromBuffer(buf.Bytes())
	parsed, err := packet.ParsePacket(in, packet.ClientSide, ^uint16(0))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	ppkt, ok := parsed.(*packet.LoginPacket)
	if !ok {
		t.Fatalf("Failed to get LoginPacket")
	}
	if !ppkt.Verify(verifier) {
		t.Fatalf("Failed to verify in stage 2")
	}
}
