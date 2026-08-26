package crypto

import (
	"crypto/ed25519"
	"testing"
)

func TestVerifyLineRateEnvelope(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	message := []byte("faultplane-envelope")
	signature := ed25519.Sign(privateKey, message)

	engine := NewHardwareCryptoEngine()

	err = engine.VerifyLineRateEnvelope(
		&VerifiedEnvelopeResult{
			Message:   message,
			Signature: signature,
			PublicKey: publicKey,
		},
	)

	if err != nil {
		t.Fatalf("valid signature rejected: %v", err)
	}
}

func TestRejectInvalidSignature(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	engine := NewHardwareCryptoEngine()

	err = engine.VerifyLineRateEnvelope(
		&VerifiedEnvelopeResult{
			Message:   []byte("tampered"),
			Signature: make([]byte, ed25519.SignatureSize),
			PublicKey: publicKey,
		},
	)

	if err != ErrCryptoVerificationFailed {
		t.Fatalf("expected verification failure, got %v", err)
	}
}
