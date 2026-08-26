package dynamic

import (
	"sync/atomic"
)

type InstructionPatchFrame struct {
	OpcodeOffset        uintptr
	NewInstructionValue uint64
	ExecutionSequence   uint64
}

type KernelHotPatcher struct {
	ActivePatchMap     atomic.Pointer[InstructionPatchFrame]
	VerifierStatusGate atomic.Uint32
}

func NewHotPatcher() *KernelHotPatcher {
	p := &KernelHotPatcher{}

	p.ActivePatchMap.Store(&InstructionPatchFrame{})
	p.VerifierStatusGate.Store(1)

	return p
}

func (p *KernelHotPatcher) InjectLiveInstructionPatch(
	targetAddr uintptr,
	byteStream uint64,
) bool {
	if p.VerifierStatusGate.Load() == 0 {
		return false
	}

	if targetAddr == 0 {
		return false
	}

	current := p.ActivePatchMap.Load()
	if current == nil {
		return false
	}

	next := &InstructionPatchFrame{
		OpcodeOffset:        targetAddr,
		NewInstructionValue: byteStream,
		ExecutionSequence:   current.ExecutionSequence + 1,
	}

	p.ActivePatchMap.Store(next)

	return true
}

func (p *KernelHotPatcher) Snapshot() InstructionPatchFrame {
	current := p.ActivePatchMap.Load()

	if current == nil {
		return InstructionPatchFrame{}
	}

	return *current
}
