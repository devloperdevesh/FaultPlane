//go:build cgo && linux

// SPDX-License-Identifier: Apache-2.0
#ifndef FAULTPLANE_VECTOR_JIT_H
#define FAULTPLANE_VECTOR_JIT_H

#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

/*
 * match_16_ips_avx512 evaluates 16 IPv4 addresses (32-bit uints) simultaneously
 * against target_ip using 512-bit vector registers (_mm512_cmpeq_epi32_mask).
 * Returns a 16-bit bitmask where bit i is 1 if ips_16[i] == target_ip.
 */
uint16_t match_16_ips_avx512(const uint32_t* ips_16, uint32_t target_ip);

/*
 * evaluate_routing_ip_avx512 searches for target_ip in ips array of length count
 * using 512-bit vector registers in 16-address batches.
 * Returns the 0-based index of the matching IP, or -1 if not found.
 */
int evaluate_routing_ip_avx512(const uint32_t* ips, size_t count, uint32_t target_ip);

/*
 * is_avx512_runtime_supported returns 1 if AVX-512 extensions were compiled into this binary.
 */
int is_avx512_runtime_supported(void);

#ifdef __cplusplus
}
#endif

#endif /* FAULTPLANE_VECTOR_JIT_H */
