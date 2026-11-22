//go:build cgo && linux

// SPDX-License-Identifier: Apache-2.0
#include "vector_jit.h"

#if defined(__x86_64__) || defined(_M_X64) || defined(__i386__) || defined(_M_IX86)
#if defined(__AVX512F__)
#include <immintrin.h>
#define FAULTPLANE_HAS_AVX512_INTRINSICS 1
#endif
#endif

extern "C" {

uint16_t match_16_ips_avx512(const uint32_t* ips_16, uint32_t target_ip) {
    if (!ips_16) {
        return 0;
    }

#if defined(FAULTPLANE_HAS_AVX512_INTRINSICS)
    __m512i ip_vec = _mm512_loadu_si512(reinterpret_cast<const void*>(ips_16));
    __m512i target_vec = _mm512_set1_epi32(target_ip);
    __mmask16 mask = _mm512_cmpeq_epi32_mask(ip_vec, target_vec);
    return static_cast<uint16_t>(mask);
#else
    uint16_t mask = 0;
    for (size_t i = 0; i < 16; ++i) {
        if (ips_16[i] == target_ip) {
            mask |= static_cast<uint16_t>(1 << i);
        }
    }
    return mask;
#endif
}

int evaluate_routing_ip_avx512(const uint32_t* ips, size_t count, uint32_t target_ip) {
    if (!ips || count == 0) {
        return -1;
    }

#if defined(FAULTPLANE_HAS_AVX512_INTRINSICS)
    __m512i target_vec = _mm512_set1_epi32(target_ip);
    size_t i = 0;

    for (; i + 16 <= count; i += 16) {
        __m512i chunk = _mm512_loadu_si512(reinterpret_cast<const void*>(ips + i));
        __mmask16 mask = _mm512_cmpeq_epi32_mask(chunk, target_vec);
        if (mask != 0) {
#if defined(__GNUC__) || defined(__clang__)
            return static_cast<int>(i + __builtin_ctz(mask));
#else
            for (int bit = 0; bit < 16; ++bit) {
                if (mask & (1 << bit)) return static_cast<int>(i + bit);
            }
#endif
        }
    }

    for (; i < count; ++i) {
        if (ips[i] == target_ip) {
            return static_cast<int>(i);
        }
    }

    return -1;
#else
    for (size_t i = 0; i < count; ++i) {
        if (ips[i] == target_ip) {
            return static_cast<int>(i);
        }
    }
    return -1;
#endif
}

int is_avx512_runtime_supported(void) {
#if defined(FAULTPLANE_HAS_AVX512_INTRINSICS)
    return 1;
#else
    return 0;
#endif
}

} // extern "C"
