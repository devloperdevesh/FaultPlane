// SPDX-License-Identifier: Apache-2.0

package fabric

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

var (
	ErrDescriptorConflict = errors.New("fabric: descriptor ownership changed")
	ErrNilDescriptor      = errors.New("fabric: descriptor must not be nil")
)

// PCIEDmaController represents the control-plane boundary for a DMA-backed
// descriptor fabric.
//
// The controller does not access physical PCIe registers directly. Hardware
// specific DMA programming belongs in a platform-specific driver layer.
type PCIEDmaController struct {
	HardwareRingAddress atomic.Uint64
	ActiveDescriptor    atomic.Pointer[byte]
}

// NewDmaController initializes the descriptor control boundary.
func NewDmaController(
	ringBase uint64,
	initialDesc unsafe.Pointer,
) *PCIEDmaController {
	controller := &PCIEDmaController{}

	controller.HardwareRingAddress.Store(ringBase)

	if initialDesc != nil {
		controller.ActiveDescriptor.Store((*byte)(initialDesc))
	}

	return controller
}

// ExecuteDescriptorHandoff atomically publishes a replacement descriptor.
//
// This operation is limited to the userspace ownership/reference boundary.
// It does not perform raw physical-memory access or bypass the operating
// system's DMA/IOMMU protections.
func (d *PCIEDmaController) ExecuteDescriptorHandoff(
	oldDesc unsafe.Pointer,
	targetDesc unsafe.Pointer,
) (bool, error) {
	if oldDesc == nil || targetDesc == nil {
		return false, ErrNilDescriptor
	}

	swapped := d.ActiveDescriptor.CompareAndSwap(
		(*byte)(oldDesc),
		(*byte)(targetDesc),
	)
	if !swapped {
		return false, ErrDescriptorConflict
	}

	return true, nil
}

// ActiveDescriptorPointer returns the currently published descriptor.
func (d *PCIEDmaController) ActiveDescriptorPointer() unsafe.Pointer {
	return unsafe.Pointer(d.ActiveDescriptor.Load())
}

// RingAddress returns the configured hardware ring identifier.
//
// The value is treated as metadata by this package; no physical memory is
// dereferenced here.
func (d *PCIEDmaController) RingAddress() uint64 {
	return d.HardwareRingAddress.Load()
}
