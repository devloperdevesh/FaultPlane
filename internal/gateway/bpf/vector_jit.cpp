//go:build cgo && linux

// SPDX-License-Identifier: Apache-2.0
#include "vector_jit.h"

#if defined(__x86_64__) || defined(_M_X64) || defined(__i386__) || defined(_M_IX86)
#if defined(__AVX512F__)
#include <immintrin.h>
#define FAULTPLANE_HAS_AVX512_INTRINSICS 1
#endif
#endif

#if defined(_MSC_VER)
#include <intrin.h>
#elif defined(__GNUC__) || defined(__clang__)
#include <cpuid.h>
#endif

static bool faultplane_avx512_runtime_supported() {
#if !defined(FAULTPLANE_HAS_AVX512_INTRINSICS)
    return false;
#else
#if defined(_MSC_VER)
    int regs[4] = {};
    __cpuidex(regs, 1, 0);
    if ((regs[2] & (1 << 27)) == 0) return false;
    unsigned long long xcr0 = _xgetbv(0);
    if ((xcr0 & 0xE6) != 0xE6) return false;
    __cpuidex(regs, 7, 0);
    return (regs[1] & (1 << 16)) != 0;
#elif defined(__GNUC__) || defined(__clang__)
    unsigned int eax = 0, ebx = 0, ecx = 0, edx = 0;
    if (!__get_cpuid(1, &eax, &ebx, &ecx, &edx)) return false;
    if ((ecx & (1u << 27)) == 0) return false;
    unsigned int xcr0_lo = 0, xcr0_hi = 0;
    __asm__ volatile(".byte 0x0f, 0x01, 0xd0" : "=a"(xcr0_lo), "=d"(xcr0_hi) : "c"(0));
    unsigned long long xcr0 = (((unsigned long long)xcr0_hi << 32) | xcr0_lo);
    if ((xcr0 & 0xE6ULL) != 0xE6ULL) return false;
    if (!__get_cpuid_count(7, 0, &eax, &ebx, &ecx, &edx)) return false;
    return (ebx & (1u << 16)) != 0;
#else
    return false;
#endif
#endif
}

extern "C" {

uint16_t match_16_ips_avx512(const uint32_t* ips_16, uint32_t target_ip) {
    if (!ips_16) {
        return 0;
    }

#if defined(FAULTPLANE_HAS_AVX512_INTRINSICS)
    if (!faultplane_avx512_runtime_supported()) {
        uint16_t mask = 0;
        for (size_t i = 0; i < 16; ++i) {
            if (ips_16[i] == target_ip) mask |= static_cast<uint16_t>(1 << i);
        }
        return mask;
    }
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
    if (!faultplane_avx512_runtime_supported()) {
        for (size_t i = 0; i < count; ++i) {
            if (ips[i] == target_ip) return static_cast<int>(i);
        }
        return -1;
    }
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
    return faultplane_avx512_runtime_supported() ? 1 : 0;
}

} // extern "C"
