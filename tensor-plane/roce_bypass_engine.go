// SPDX-License-Identifier: Apache-2.0
//go:build linux

package tensorplane

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const (
	RDMAWriteWithImmediate uint32 = 1
)

var (
	ErrRoCEInactive = errors.New("RoCE engine is inactive")
	ErrInvalidMR    = errors.New("invalid memory region")
)

// RoceMemoryRegion represents a memory region registered with the RDMA
// provider. The provider, not FaultPlane, owns the actual lkey/rkey values.
type RoceMemoryRegion struct {
	Address uint64
	Length  uint32
	LKey    uint32
	RKey    uint32
}

// RoceWorkRequest is the provider-independent representation of an RDMA write.
type RoceWorkRequest struct {
	ID            uint64
	Opcode        uint32
	LocalAddress  uint64
	LocalLength   uint32
	LocalKey      uint32
	RemoteAddress uint64
	RemoteKey     uint32
	Immediate     uint32
}

// RoceProvider is the real hardware boundary.
//
// A Linux RDMA provider should implement this using the RDMA userspace verbs
// API for the selected HCA/QP. No fake device FD or fake keys are accepted.
type RoceProvider interface {
	RegisterMemory(addr uintptr, length uint32, access uint32) (RoceMemoryRegion, error)
	PostWrite(wr RoceWorkRequest) error
	DeregisterMemory(mr RoceMemoryRegion) error
	Close() error
}

const (
	RoceAccessLocalWrite  uint32 = 1 << 0
	RoceAccessRemoteWrite uint32 = 1 << 1
	RoceAccessRemoteRead  uint32 = 1 << 2
)

type RoceBypassEngine struct {
	active   uint32
	provider RoceProvider
	localMR  RoceMemoryRegion
}

// NewRoceBypassEngine registers an existing memory region with the real RDMA
// provider.
func NewRoceBypassEngine(
	provider RoceProvider,
	buffer uintptr,
	length uint32,
) (*RoceBypassEngine, error) {
	if provider == nil {
		return nil, errors.New("nil RoCE provider")
	}

	if buffer == 0 || length == 0 {
		return nil, ErrInvalidMR
	}

	mr, err := provider.RegisterMemory(
		buffer,
		length,
		RoceAccessLocalWrite|RoceAccessRemoteWrite,
	)
	if err != nil {
		return nil, fmt.Errorf("register RoCE memory: %w", err)
	}

	return &RoceBypassEngine{
		active:   1,
		provider: provider,
		localMR:  mr,
	}, nil
}

// InjectRemoteServerRAM posts an RDMA write using provider-issued addressing
// and authorization keys.
func (e *RoceBypassEngine) InjectRemoteServerRAM(
	remoteAddress uint64,
	remoteKey uint32,
	localAddress uint64,
	length uint32,
) error {
	if atomic.LoadUint32(&e.active) == 0 {
		return ErrRoCEInactive
	}

	if length == 0 || length > e.localMR.Length {
		return fmt.Errorf("invalid transfer length %d", length)
	}

	wr := RoceWorkRequest{
		ID:            uint64(localAddress),
		Opcode:        RDMAWriteWithImmediate,
		LocalAddress:  localAddress,
		LocalLength:   length,
		LocalKey:      e.localMR.LKey,
		RemoteAddress: remoteAddress,
		RemoteKey:     remoteKey,
	}

	if err := e.provider.PostWrite(wr); err != nil {
		return fmt.Errorf("post RoCE RDMA write: %w", err)
	}

	return nil
}

// ShutdownRoCEEngine deregisters memory and closes the provider.
func (e *RoceBypassEngine) ShutdownRoCEEngine() error {
	if !atomic.CompareAndSwapUint32(&e.active, 1, 0) {
		return nil
	}

	var firstErr error

	if err := e.provider.DeregisterMemory(e.localMR); err != nil {
		firstErr = fmt.Errorf("deregister RoCE memory: %w", err)
	}

	if err := e.provider.Close(); err != nil && firstErr == nil {
		firstErr = fmt.Errorf("close RoCE provider: %w", err)
	}

	return firstErr
}
