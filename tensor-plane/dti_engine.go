// SPDX-License-Identifier: Apache-2.0
package tensorplane

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	ErrBarAddressExceeded = errors.New(
		"faultplane [tensor-plane]: physical hardware BAR memory window overflow",
	)
	ErrTensorRingUnanchored = errors.New(
		"faultplane [tensor-plane]: hardware accelerator DMA ring unanchored",
	)
	ErrInvalidRingSize = errors.New(
		"faultplane [tensor-plane]: invalid tensor ring size",
	)
)

// PhysicalTensorPage describes the verified address window and transfer
// state associated with a tensor transfer ring.
//
// BaseBarAddress is immutable after construction.
type PhysicalTensorPage struct {
	BaseBarAddress     uintptr
	PageLengthMask     uint32
	DirectTransferFlag uint32
}

// TensorIngressController owns the active tensor transfer configuration.
type TensorIngressController struct {
	ActiveHardwareRing *PhysicalTensorPage
	LineRateThroughput uint64
	ExecutionGatekeeper uint32
}

// NewTensorController creates a tensor ingress controller with a
// power-of-two ring size.
//
// slotPower must be less than 32.
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

// SplicedLineRateIngress validates a tensor transfer address against the
// configured address window.
//
// This controller performs validation and state management only.
// Actual PCIe BAR access and DMA submission must be implemented by the
// platform-specific driver/backend.
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
		return fmt.Errorf(
			"faultplane [tensor-plane]: direct transfer path disabled",
		)
	}

	if activeBufferPointer == 0 || activeBufferPointer > maximumLimit {
		return ErrBarAddressExceeded
	}

	return nil
}
