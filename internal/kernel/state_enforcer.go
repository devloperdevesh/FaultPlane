// SPDX-License-Identifier: Apache-2.0
//go:build linux

package kernel

// ProductionKernelEnforcer represents the Linux kernel-state enforcement layer.
// Concrete eBPF attachment is implemented behind the Linux build boundary.
type ProductionKernelEnforcer struct{}
