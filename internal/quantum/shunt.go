package quantum

import (
	"sync/atomic"
)

type SpeculativeDescriptorFrame struct {
	PrimarySocketFD     uintptr
	SpeculativeFD       uintptr
	ShuntExecutionCount uint64
	QuantumLockFlag     uint32
}

type QuantumShuntController struct {
	ActiveFrame         atomic.Pointer[SpeculativeDescriptorFrame]
	TelemetryDeviation  atomic.Int64
	EnforcementRegistry atomic.Uint32
}

func NewQuantumController(
	primary,
	standby uintptr,
) *QuantumShuntController {
	c := &QuantumShuntController{}

	c.ActiveFrame.Store(&SpeculativeDescriptorFrame{
		PrimarySocketFD: primary,
		SpeculativeFD:   standby,
	})

	c.EnforcementRegistry.Store(1)

	return c
}

func (c *QuantumShuntController) InterceptAndSpeculateShunt(
	currentJitterNS int64,
	thresholdNS int64,
) bool {
	c.TelemetryDeviation.Store(currentJitterNS)

	if c.EnforcementRegistry.Load() == 0 {
		return false
	}

	if currentJitterNS <= thresholdNS {
		return false
	}

	frame := c.ActiveFrame.Load()
	if frame == nil {
		return false
	}

	if !atomic.CompareAndSwapUint32(
		&frame.QuantumLockFlag,
		0,
		1,
	) {
		return false
	}

	atomic.AddUint64(
		&frame.ShuntExecutionCount,
		1,
	)

	return true
}

func (c *QuantumShuntController) Snapshot() SpeculativeDescriptorFrame {
	frame := c.ActiveFrame.Load()

	if frame == nil {
		return SpeculativeDescriptorFrame{}
	}

	return SpeculativeDescriptorFrame{
		PrimarySocketFD: frame.PrimarySocketFD,
		SpeculativeFD:   frame.SpeculativeFD,
		ShuntExecutionCount: atomic.LoadUint64(
			&frame.ShuntExecutionCount,
		),
		QuantumLockFlag: atomic.LoadUint32(
			&frame.QuantumLockFlag,
		),
	}
}
