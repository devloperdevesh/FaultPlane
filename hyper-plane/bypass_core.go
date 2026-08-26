// SPDX-License-Identifier: Apache-2.0
package hyperplane

import (
	"errors"
	"sync/atomic"
)

var (
	ErrContextContention = errors.New(
		"faultplane [hyper-plane]: target buffer outside configured execution window",
	)

	ErrHardwareFabricCold = errors.New(
		"faultplane [hyper-plane]: asynchronous execution fabric is inactive",
	)

	ErrInvalidMemoryWindow = errors.New(
		"faultplane [hyper-plane]: invalid execution memory window",
	)

	ErrExecutionPathDisabled = errors.New(
		"faultplane [hyper-plane]: asynchronous execution path is disabled",
	)
)

// SovereignExecutionDescriptor describes the configured execution window
// and transfer state.
//
// ActivePCIeBaseAddr is immutable after construction. Actual PCIe register
// access and DMA submission must be performed by a platform-specific driver.
type SovereignExecutionDescriptor struct {
	ActivePCIeBaseAddr uintptr
	DirectMemoryWindow uint64
	AsynchronousShunt  uint32
}

// HyperPlaneController owns the active execution descriptor and runtime
// safety state.
type HyperPlaneController struct {
	SovereignDescriptor  *SovereignExecutionDescriptor
	VolumetricTokenCount uint64
	SafetyEnforcement    uint32
}

// NewHyperPlaneController creates a HyperPlane controller with a validated
// execution memory window.
func NewHyperPlaneController(
	physicalBAR uintptr,
	memoryLimit uint64,
) (*HyperPlaneController, error) {
	if physicalBAR == 0 || memoryLimit == 0 {
		return nil, ErrInvalidMemoryWindow
	}

	descriptor := &SovereignExecutionDescriptor{
		ActivePCIeBaseAddr: physicalBAR,
		DirectMemoryWindow: memoryLimit,
		AsynchronousShunt:  1,
	}

	return &HyperPlaneController{
		SovereignDescriptor:  descriptor,
		VolumetricTokenCount: 0,
		SafetyEnforcement:    1,
	}, nil
}

// ExecuteAsynchronousShunt validates a target buffer against the configured
// execution window.
//
// This method does not directly bypass the operating system scheduler,
// access PCIe registers, or submit DMA. Those operations belong to the
// platform-specific driver/backend.
func (c *HyperPlaneController) ExecuteAsynchronousShunt(
	targetBufferPointer uintptr,
	maximumBound uintptr,
) error {
	if c == nil || c.SovereignDescriptor == nil {
		return ErrHardwareFabricCold
	}

	atomic.AddUint64(&c.VolumetricTokenCount, 1)

	if atomic.LoadUint32(&c.SafetyEnforcement) == 0 {
		return ErrHardwareFabricCold
	}

	desc := c.SovereignDescriptor

	if atomic.LoadUint32(&desc.AsynchronousShunt) == 0 {
		return ErrExecutionPathDisabled
	}

	if targetBufferPointer == 0 || targetBufferPointer > maximumBound {
		return ErrContextContention
	}

	if desc.DirectMemoryWindow != 0 &&
		uint64(targetBufferPointer) > desc.DirectMemoryWindow {
		return ErrContextContention
	}

	return nil
}
