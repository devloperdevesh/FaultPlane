// SPDX-License-Identifier: Apache-2.0
package tensorplane

import (
	"errors"
	"sync/atomic"
)

var (
	ErrBarAddressExceeded = errors.New(
		"faultplane [tensor-plane]: buffer address outside configured BAR window",
	)

	ErrTensorRingUnanchored = errors.New(
		"faultplane [tensor-plane]: tensor ring is not active",
	)

	ErrInvalidRingSize = errors.New(
		"faultplane [tensor-plane]: invalid ring size",
	)

	ErrTransferDisabled = errors.New(
		"faultplane [tensor-plane]: transfer path is disabled",
	)
)

// PhysicalTensorPage describes the verified address window associated
// with a tensor transfer ring.
//
// BaseBarAddress is treated as immutable after construction.
type PhysicalTensorPage struct {
	BaseBarAddress     uintptr
	PageLengthMask     uint32
	DirectTransferFlag uint32
}

// TensorIngressController owns the active tensor transfer configuration.
type TensorIngressController struct {
	ActiveHardwareRing  *PhysicalTensorPage
	LineRateThroughput  uint64
	ExecutionGatekeeper uint32
}

// NewTensorController creates a controller for a power-of-two ring.
//
// slotPower must produce a ring size representable by uint32.
func NewTensorController(
	hardwareBAR uintptr,
	slotPower uint32,
) (*TensorIngressController, error) {
	if slotPower >= 32 {
		return nil, ErrInvalidRingSize
	}

	ringSize := uint32(1) << slotPower
	if ringSize == 0 {
		return nil, ErrInvalidRingSize
	}

	page := &PhysicalTensorPage{
		BaseBarAddress:     hardwareBAR,
		PageLengthMask:     ringSize - 1,
		DirectTransferFlag: 1,
	}

	return &TensorIngressController{
		ActiveHardwareRing:  page,
		LineRateThroughput:  0,
		ExecutionGatekeeper: 1,
	}, nil
}

// SplicedLineRateIngress validates a transfer address against the
// configured address window.
//
// This method does not perform DMA or directly access PCIe hardware;
// the actual hardware operation must be supplied by the platform
// driver/backend.
func (c *TensorIngressController) SplicedLineRateIngress(
	activeBufferPointer uintptr,
	maximumLimit uintptr,
) error {
	if c == nil || c.ActiveHardwareRing == nil {
		return ErrTensorRingUnanchored
	}

	atomic.AddUint64(&c.LineRateThroughput, 1)

	if atomic.LoadUint32(&c.ExecutionGatekeeper) == 0 {
		return ErrTensorRingUnanchored
	}

	page := c.ActiveHardwareRing

	if atomic.LoadUint32(&page.DirectTransferFlag) == 0 {
		return ErrTransferDisabled
	}

	if activeBufferPointer == 0 || activeBufferPointer > maximumLimit {
		return ErrBarAddressExceeded
	}

	return nil
}
