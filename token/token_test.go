package token

import (
	"testing"
)

var privkey PrivKey
var pubkey PubKey

func init() {
	pubkey, privkey = GenerateKeys()
}

func TestGetTokenGenerator(t *testing.T) {
	gen := GetTokenGenerator(privkey)
	if gen == nil {
		t.Fatal("Cannot get generator")
	}
}

func TestGetTokenVerifier(t *testing.T) {
	ver := GetTokenVerifier(pubkey)
	if ver == nil {
		t.Fatal("Cannot get verifier")
	}
}

func TestGenerateKeys(t *testing.T) {
	gen := GetTokenGenerator(privkey)
	ver := GetTokenVerifier(pubkey)
	data := []byte("test")
	var token Token
	gen(data, &token)
	if !ver(data, token) {
		t.Fatal("Cannot generate correct key pair")
	}
}
