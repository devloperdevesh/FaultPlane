// SPDX-License-Identifier: Apache-2.0
//go:build linux

package tensorplane

import (
	"errors"
	"fmt"
	"sync/atomic"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	HugePageSize        = 2 * 1024 * 1024
	RingDescriptorCount = 4096
	CacheLineSize       = 64
	DescriptorSize      = int(unsafe.Sizeof(HardwareDescriptor{}))
)

var (
	ErrSplicerClosed = errors.New("pcie splicer is closed")
	ErrRingFull      = errors.New("pcie descriptor ring is full")
)

// HardwareDescriptor is a software-owned descriptor.
// PhysicalAddress must be supplied by a trusted DMA/IOMMU allocator;
// this type never converts an arbitrary userspace pointer into a DMA address.
type HardwareDescriptor struct {
	PhysicalAddress uint64
	BufferLength    uint32
	OwnershipFlag   uint32 // 0 = software, 1 = hardware
}

type PcieDMAProvider interface {
	Allocate(length uint32) (dmaAddress uint64, backing []byte, err error)
	Release(dmaAddress uint64) error
}

type PcieSplicerEngine struct {
	state uint32

	ring []HardwareDescriptor

	head uint64
	_    [CacheLineSize - 8]byte
	tail uint64

	dma PcieDMAProvider
}

// NewPcieSplicerEngine creates a descriptor ring backed by anonymous memory.
//
// This is the software half of the PCIe data path. Actual PCIe DMA addresses
// must come from the platform's DMA/IOMMU layer; arbitrary virtual addresses
// are deliberately not treated as physical addresses.
func NewPcieSplicerEngine(dma PcieDMAProvider) (*PcieSplicerEngine, error) {
	if dma == nil {
		return nil, errors.New("nil PCIe DMA provider")
	}

	ringBytes := RingDescriptorCount * DescriptorSize

	// Try huge pages first. Fall back to normal anonymous memory if the host
	// has no preallocated huge pages.
	flags := unix.MAP_PRIVATE | unix.MAP_ANONYMOUS
	mem, err := unix.Mmap(
		-1,
		0,
		ringBytes,
		unix.PROT_READ|unix.PROT_WRITE,
		flags|unix.MAP_HUGETLB,
	)
	if err != nil {
		mem, err = unix.Mmap(
			-1,
			0,
			ringBytes,
			unix.PROT_READ|unix.PROT_WRITE,
			flags,
		)
		if err != nil {
			return nil, fmt.Errorf("map PCIe descriptor ring: %w", err)
		}
	}

	ring := unsafe.Slice(
		(*HardwareDescriptor)(unsafe.Pointer(&mem[0])),
		RingDescriptorCount,
	)

	for i := range ring {
		ring[i] = HardwareDescriptor{}
	}

	return &PcieSplicerEngine{
		state: 1,
		ring:  ring,
		dma:   dma,
	}, nil
}

// QueueBuffer allocates a DMA-capable buffer through the platform provider
// and publishes its DMA address into the descriptor ring.
func (e *PcieSplicerEngine) QueueBuffer(payload []byte) (uint64, error) {
	if atomic.LoadUint32(&e.state) == 0 {
		return 0, ErrSplicerClosed
	}

	if len(payload) == 0 {
		return 0, errors.New("empty payload")
	}

	if len(payload) > int(^uint32(0)) {
		return 0, errors.New("payload too large")
	}

	head := atomic.LoadUint64(&e.head)
	tail := atomic.LoadUint64(&e.tail)

	if tail-head >= RingDescriptorCount {
		return 0, ErrRingFull
	}

	dmaAddress, backing, err := e.dma.Allocate(uint32(len(payload)))
	if err != nil {
		return 0, fmt.Errorf("allocate DMA buffer: %w", err)
	}

	copy(backing, payload)

	index := tail & (RingDescriptorCount - 1)
	desc := &e.ring[index]

	desc.PhysicalAddress = dmaAddress
	desc.BufferLength = uint32(len(payload))

	// Publish ownership only after descriptor contents are complete.
	atomic.StoreUint32(&desc.OwnershipFlag, 1)
	atomic.StoreUint64(&e.tail, tail+1)

	return dmaAddress, nil
}

// Complete marks a descriptor as consumed by the device and releases its DMA
// allocation through the provider.
func (e *PcieSplicerEngine) Complete(index uint64) error {
	if atomic.LoadUint32(&e.state) == 0 {
		return ErrSplicerClosed
	}

	slot := index & (RingDescriptorCount - 1)
	desc := &e.ring[slot]

	if atomic.LoadUint32(&desc.OwnershipFlag) == 0 {
		return nil
	}

	dmaAddress := desc.PhysicalAddress

	atomic.StoreUint32(&desc.OwnershipFlag, 0)
	desc.PhysicalAddress = 0
	desc.BufferLength = 0

	if err := e.dma.Release(dmaAddress); err != nil {
		return fmt.Errorf("release DMA buffer 0x%x: %w", dmaAddress, err)
	}

	atomic.StoreUint64(&e.head, index+1)
	return nil
}

// Close releases the descriptor mapping.
func (e *PcieSplicerEngine) Close() error {
	if !atomic.CompareAndSwapUint32(&e.state, 1, 0) {
		return nil
	}

	if len(e.ring) == 0 {
		return nil
	}

	ptr := unsafe.Pointer(&e.ring[0])
	size := len(e.ring) * DescriptorSize
	mapped := unsafe.Slice((*byte)(ptr), size)

	if err := unix.Munmap(mapped); err != nil {
		return fmt.Errorf("unmap PCIe descriptor ring: %w", err)
	}

	e.ring = nil
	return nil
}
