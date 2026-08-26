// SPDX-License-Identifier: Apache-2.0
package controlplane

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	ErrTopologyMismatched = errors.New(
		"faultplane [control-plane]: CPU core allocation does not match policy",
	)
	ErrBandwidthSaturated = errors.New(
		"faultplane [control-plane]: configured bandwidth quota exceeded",
	)
	ErrInvalidPolicy = errors.New(
		"faultplane [control-plane]: invalid policy configuration",
	)
)

type TopologyPinningMatrix struct {
	DedicatedCoreID     uint32
	MaxBandwidthQuotaGB uint64
	PolicyEnforceFlag   uint32
}

type PolicyEngine struct {
	ActiveMatrix        *TopologyPinningMatrix
	MonitoredThroughput uint64
	EngineLockGate      uint32
}

func NewPolicyEngine(
	targetCore uint32,
	bandwidthLimit uint64,
) (*PolicyEngine, error) {
	if bandwidthLimit == 0 {
		return nil, ErrInvalidPolicy
	}

	matrix := &TopologyPinningMatrix{
		DedicatedCoreID:     targetCore,
		MaxBandwidthQuotaGB: bandwidthLimit,
		PolicyEnforceFlag:   1,
	}

	return &PolicyEngine{
		ActiveMatrix:        matrix,
		MonitoredThroughput: 0,
		EngineLockGate:      1,
	}, nil
}

func (e *PolicyEngine) EnforceSovereignTopology(
	activeCore uintptr,
	totalTransitBytes uint64,
) error {
	if e == nil || e.ActiveMatrix == nil {
		return ErrInvalidPolicy
	}

	atomic.AddUint64(&e.MonitoredThroughput, totalTransitBytes)

	if atomic.LoadUint32(&e.EngineLockGate) == 0 {
		return ErrTopologyMismatched
	}

	matrix := e.ActiveMatrix

	if atomic.LoadUint32(&matrix.PolicyEnforceFlag) == 0 {
		return fmt.Errorf(
			"faultplane [control-plane]: policy enforcement is disabled",
		)
	}

	if activeCore != uintptr(matrix.DedicatedCoreID) {
		return ErrTopologyMismatched
	}

	if totalTransitBytes > matrix.MaxBandwidthQuotaGB {
		return ErrBandwidthSaturated
	}

	return nil
}
