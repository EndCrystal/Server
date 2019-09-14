package token

import "crypto/ed25519"

var (
	TokenLen   = ed25519.SignatureSize
	PubKeyLen  = ed25519.PublicKeySize
	PrivKeyLen = ed25519.PrivateKeySize
)

type (
	Token   [ed25519.SignatureSize]byte
	PubKey  [ed25519.PublicKeySize]byte
	PrivKey [ed25519.PrivateKeySize]byte
)

type (
	TokenGenerator func(data []byte, tok Token)
	TokenVerifier  func(data []byte, tok Token) bool
)

func GetTokenGenerator(priv PrivKey) TokenGenerator {
	return func(data []byte, tok Token) {
		s := ed25519.Sign(priv[:], data)
		copy(tok[:], s)
		return
	}
}

func GetTokenVerifier(pub PubKey) TokenVerifier {
	return func(data []byte, tok Token) bool { return ed25519.Verify(pub[:], data, tok[:]) }
}

func GenerateKeys() (pub PubKey, priv PrivKey) {
	pubs, privs, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	copy(pub[:], pubs)
	copy(priv[:], privs)
	return
}
