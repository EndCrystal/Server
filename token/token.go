package token

import "crypto/ed25519"

const (
	// TokenLen length for token
	TokenLen   = ed25519.SignatureSize
	// PubKeyLen length of public key
	PubKeyLen  = ed25519.PublicKeySize
	// PrivKeyLen length of private key
	PrivKeyLen = ed25519.PrivateKeySize
)

type (
	// Token token type
	Token   [ed25519.SignatureSize]byte
	// PubKey public key
	PubKey  [ed25519.PublicKeySize]byte
	// PrivKey private key
	PrivKey [ed25519.PrivateKeySize]byte
)

type (
	// Generator generate token
	Generator func(data []byte, tok *Token)

	// Verifier verify token
	Verifier  func(data []byte, tok Token) bool
)

// GetTokenGenerator create generator
func GetTokenGenerator(priv PrivKey) Generator {
	return func(data []byte, tok *Token) {
		s := ed25519.Sign(priv[:], data)
		copy(tok[:], s)
		return
	}
}

// GetTokenVerifier create verifier
func GetTokenVerifier(pub PubKey) Verifier {
	return func(data []byte, tok Token) bool { return ed25519.Verify(pub[:], data, tok[:]) }
}

// GenerateKeys generate key pair
func GenerateKeys() (pub PubKey, priv PrivKey) {
	pubs, privs, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	copy(pub[:], pubs)
	copy(priv[:], privs)
	return
}
