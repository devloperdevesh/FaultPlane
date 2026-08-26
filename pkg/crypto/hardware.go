package crypto

import (
	"crypto/ed25519"
	"errors"
	"sync/atomic"
)

var ErrCryptoVerificationFailed = errors.New(
	"faultplane: cryptographic signature verification failed",
)

type HardwareCryptoEngine struct {
	TotalVerifiedPackets uint64
	ActiveHardwareMode   uint32
}

func NewHardwareCryptoEngine() *HardwareCryptoEngine {
	return &HardwareCryptoEngine{
		ActiveHardwareMode: 0,
	}
}

type VerifiedEnvelopeResult struct {
	Message   []byte
	Signature []byte
	PublicKey ed25519.PublicKey
}

func (e *HardwareCryptoEngine) VerifyLineRateEnvelope(
	result *VerifiedEnvelopeResult,
) error {
	if result == nil ||
		len(result.PublicKey) != ed25519.PublicKeySize ||
		len(result.Signature) != ed25519.SignatureSize {
		return ErrCryptoVerificationFailed
	}

	if !ed25519.Verify(
		result.PublicKey,
		result.Message,
		result.Signature,
	) {
		return ErrCryptoVerificationFailed
	}

	atomic.AddUint64(
		&e.TotalVerifiedPackets,
		1,
	)

	return nil
}
