// SPDX-License-Identifier: Apache-2.0
//go:build linux

package tensorplane

// UltraBypassDMAEngine is reserved for the Linux shared-memory transport path.
// Implementation is added only after the platform and memory-ordering contract
// has been verified against the existing tensor-plane interfaces.
type UltraBypassDMAEngine struct{}
