package verifier

import (
	"errors"
	"sync/atomic"
)

var ErrBoundaryViolation = errors.New(
	"faultplane: runtime state constraint out of bounds",
)

type RuntimeSafetyVerifier struct {
	TotalProcessedLoops uint64
	ViolationInvariants uint64
	ActiveEnforcement   uint32
}

func NewSafetyVerifier() *RuntimeSafetyVerifier {
	return &RuntimeSafetyVerifier{
		ActiveEnforcement: 1,
	}
}

func (v *RuntimeSafetyVerifier) AssertStateInvariant(
	allocatedDescriptor uintptr,
	upperLimit uintptr,
) error {
	atomic.AddUint64(&v.TotalProcessedLoops, 1)

	if atomic.LoadUint32(&v.ActiveEnforcement) == 0 {
		return nil
	}

	if allocatedDescriptor == 0 ||
		upperLimit == 0 ||
		allocatedDescriptor > upperLimit {
		atomic.AddUint64(&v.ViolationInvariants, 1)
		return ErrBoundaryViolation
	}

	return nil
}

func (v *RuntimeSafetyVerifier) Snapshot() (
	uint64,
	uint64,
	bool,
) {
	return atomic.LoadUint64(&v.TotalProcessedLoops),
		atomic.LoadUint64(&v.ViolationInvariants),
		atomic.LoadUint32(&v.ActiveEnforcement) == 1
}
