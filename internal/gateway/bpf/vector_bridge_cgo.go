//go:build cgo && linux && (amd64 || 386)

// SPDX-License-Identifier: Apache-2.0
package bpf

/*
#cgo CXXFLAGS: -O3 -mavx512f -std=c++17
#cgo LDFLAGS: -lstdc++
#include "vector_jit.h"
*/
import "C"
import (
	"unsafe"
)

// CgoEnabled denotes native Cgo compilation profile availability.
const CgoEnabled = true

// NativeMatch16IPs bridges Go userspace slice to C++ AVX-512 register instructions.
func NativeMatch16IPs(ips []uint32, targetIP uint32) uint16 {
	if len(ips) < 16 {
		return 0
	}
	ptr := (*C.uint32_t)(unsafe.Pointer(&ips[0]))
	return uint16(C.match_16_ips_avx512(ptr, C.uint32_t(targetIP)))
}

// NativeEvaluateIPs scans routing allocation IPs using 512-bit vector registers.
func NativeEvaluateIPs(ips []uint32, targetIP uint32) int {
	if len(ips) == 0 {
		return -1
	}
	ptr := (*C.uint32_t)(unsafe.Pointer(&ips[0]))
	return int(C.evaluate_routing_ip_avx512(ptr, C.size_t(len(ips)), C.uint32_t(targetIP)))
}

// IsAVX512Compiled checks if AVX-512 extensions were enabled during native compilation.
func IsAVX512Compiled() bool {
	return C.is_avx512_runtime_supported() == 1
}
