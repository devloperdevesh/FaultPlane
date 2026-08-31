// SPDX-License-Identifier: Apache-2.0
//go:build linux

package tensorplane

import (
	"errors"
	"fmt"
	"sync/atomic"

	"golang.org/x/sys/unix"
)

var (
	ErrBusInactive = errors.New("PCIe bus splicer is inactive")
	ErrInvalidBAR  = errors.New("invalid PCIe BAR mapping")
)

// PcieBAR abstracts one PCIe BAR exposed by the operating system.
//
// Mapping must come from the actual PCI device resource file, for example
// /sys/bus/pci/devices/<BDF>/resourceN, after the device has been validated.
type PcieBAR interface {
	Mapping() []byte
	Close() error
}

// PcieBusSplicer coordinates validated BAR mappings.
//
// It intentionally does not assume that a NIC BAR has NVMe registers at
// hard-coded offsets. Device-specific register layouts belong in adapters.
type PcieBusSplicer struct {
	active uint32

	network PcieBAR
	storage PcieBAR
}

// NewPcieBusSplicer validates both BAR mappings.
func NewPcieBusSplicer(networkBAR, storageBAR PcieBAR) (*PcieBusSplicer, error) {
	if networkBAR == nil || storageBAR == nil {
		return nil, errors.New("network and storage BAR providers are required")
	}

	if len(networkBAR.Mapping()) == 0 {
		return nil, ErrInvalidBAR
	}

	if len(storageBAR.Mapping()) == 0 {
		_ = networkBAR.Close()
		return nil, ErrInvalidBAR
	}

	return &PcieBusSplicer{
		active:  1,
		network: networkBAR,
		storage: storageBAR,
	}, nil
}

// NetworkBAR returns the validated network BAR mapping.
func (s *PcieBusSplicer) NetworkBAR() ([]byte, error) {
	if atomic.LoadUint32(&s.active) == 0 {
		return nil, ErrBusInactive
	}

	return s.network.Mapping(), nil
}

// StorageBAR returns the validated storage BAR mapping.
func (s *PcieBusSplicer) StorageBAR() ([]byte, error) {
	if atomic.LoadUint32(&s.active) == 0 {
		return nil, ErrBusInactive
	}

	return s.storage.Mapping(), nil
}

// SyncDeviceMapping synchronizes the mapped BAR range.
//
// Actual device register writes must be implemented by a device-specific
// adapter after validating the PCI device ID, BAR size, register offset and
// register semantics.
func (s *PcieBusSplicer) SyncDeviceMapping(bar PcieBAR) error {
	if atomic.LoadUint32(&s.active) == 0 {
		return ErrBusInactive
	}

	if bar == nil || len(bar.Mapping()) == 0 {
		return ErrInvalidBAR
	}

	// Prevent the compiler/runtime from treating the mapped range as unused.
	// This does not perform an unsafe hardware register write.
	if err := unix.Msync(bar.Mapping(), unix.MS_SYNC); err != nil {
		return fmt.Errorf("sync PCIe BAR mapping: %w", err)
	}

	return nil
}

// Close disables the bus splicer and releases both BAR mappings.
func (s *PcieBusSplicer) Close() error {
	if !atomic.CompareAndSwapUint32(&s.active, 1, 0) {
		return nil
	}

	var firstErr error

	if s.network != nil {
		if err := s.network.Close(); err != nil {
			firstErr = fmt.Errorf("close network BAR: %w", err)
		}
	}

	if s.storage != nil {
		if err := s.storage.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close storage BAR: %w", err)
		}
	}

	return firstErr
}

// MappedBAR is a concrete Linux BAR mapping helper.
type MappedBAR struct {
	data []byte
}

// MapBAR maps a validated PCI resource file.
func MapBAR(fd int, length int) (*MappedBAR, error) {
	if fd < 0 {
		return nil, errors.New("invalid BAR file descriptor")
	}

	if length <= 0 {
		return nil, errors.New("invalid BAR length")
	}

	data, err := unix.Mmap(
		fd,
		0,
		length,
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_SHARED,
	)
	if err != nil {
		return nil, fmt.Errorf("mmap PCI BAR: %w", err)
	}

	return &MappedBAR{data: data}, nil
}

func (b *MappedBAR) Mapping() []byte {
	if b == nil {
		return nil
	}

	return b.data
}

func (b *MappedBAR) Close() error {
	if b == nil || b.data == nil {
		return nil
	}

	err := unix.Munmap(b.data)
	b.data = nil
	return err
}
