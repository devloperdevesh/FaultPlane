//go:build !cgo || !linux || (!amd64 && !386)

// SPDX-License-Identifier: Apache-2.0
package bpf

// CgoEnabled denotes native Cgo compilation profile availability.
const CgoEnabled = false

// NativeMatch16IPs is the zero-overhead fallback stub when Cgo/AVX-512 is uncompiled.
func NativeMatch16IPs(ips []uint32, targetIP uint32) uint16 {
	return 0
}

// NativeEvaluateIPs is the zero-overhead fallback stub when Cgo/AVX-512 is uncompiled.
func NativeEvaluateIPs(ips []uint32, targetIP uint32) int {
	return -1
}

// IsAVX512Compiled reports false when native SIMD extensions are uncompiled.
func IsAVX512Compiled() bool {
	return false
}
