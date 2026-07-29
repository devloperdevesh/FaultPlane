/*
 * Hardware-Accelerated AES-NI Stream Encryption with Memory Barrier Fencing
 * Target Architecture: x86_64 with AES-NI support
 * Compilation: gcc -O3 -maes -msse4.1 -shared -fPIC -o libaes_ni_fence.so aes_ni_fence.c
 */

#include <wmmintrin.h> // AES-NI intrinsics
#include <emmintrin.h> // SSE2 intrinsics
#include <smmintrin.h> // SSE4.1 intrinsics
#include <stdint.h>
#include <string.h>

#define BLOCK_SIZE 16
#define PARALLEL_BLOCKS 4

// Key expansion for AES-128
static inline void aes128_key_expansion(__m128i key, __m128i *round_keys) {
    round_keys[0] = key;
    
    #define EXPAND_KEY(i, rcon) \
        { \
            __m128i temp = _mm_aeskeygenassist_si128(round_keys[i-1], rcon); \
            temp = _mm_shuffle_epi32(temp, 0xff); \
            __m128i key_temp = round_keys[i-1]; \
            key_temp = _mm_xor_si128(key_temp, _mm_slli_si128(key_temp, 4)); \
            key_temp = _mm_xor_si128(key_temp, _mm_slli_si128(key_temp, 4)); \
            key_temp = _mm_xor_si128(key_temp, _mm_slli_si128(key_temp, 4)); \
            round_keys[i] = _mm_xor_si128(key_temp, temp); \
        }

    EXPAND_KEY(1, 0x01);
    EXPAND_KEY(2, 0x02);
    EXPAND_KEY(3, 0x04);
    EXPAND_KEY(4, 0x08);
    EXPAND_KEY(5, 0x10);
    EXPAND_KEY(6, 0x20);
    EXPAND_KEY(7, 0x40);
    EXPAND_KEY(8, 0x80);
    EXPAND_KEY(9, 0x1B);
    EXPAND_KEY(10, 0x36);
    #undef EXPAND_KEY
}

// Increment 128-bit big-endian counter vector
static inline __m128i inc_counter(__m128i ctr) {
    static const uint8_t ONE[16] __attribute__((aligned(16))) = {
        0,0,0,0, 0,0,0,0, 0,0,0,0, 0,0,0,1
    };
    __m128i one = _mm_load_si128((const __m128i*)ONE);
    return _mm_add_epi64(ctr, one);
}

/**
 * aes_ni_ctr_process_fenced
 * Encrypts/Decrypts buffer in CTR mode using AES-NI with strict memory bit fencing.
 */
void aes_ni_ctr_process_fenced(
    const uint8_t *in,
    uint8_t *out,
    size_t length,
    const uint8_t key[16],
    uint8_t iv[16]
) {
    if (!in || !out || length == 0) return;

    __m128i round_keys[11];
    __m128i raw_key = _mm_loadu_si128((const __m128i*)key);
    aes128_key_expansion(raw_key, round_keys);

    __m128i ctr0 = _mm_loadu_si128((const __m128i*)iv);

    size_t offset = 0;

    // Load fence before processing sensitive buffers
    _mm_lfence();

    // Process 4 blocks (64 bytes) in parallel
    while (offset + (PARALLEL_BLOCKS * BLOCK_SIZE) <= length) {
        __m128i c0 = ctr0;
        __m128i c1 = inc_counter(c0);
        __m128i c2 = inc_counter(c1);
        __m128i c3 = inc_counter(c2);
        ctr0 = inc_counter(c3);

        // Initial Round
        c0 = _mm_xor_si128(c0, round_keys[0]);
        c1 = _mm_xor_si128(c1, round_keys[0]);
        c2 = _mm_xor_si128(c2, round_keys[0]);
        c3 = _mm_xor_si128(c3, round_keys[0]);

        // Main Rounds (1-9)
        for (int r = 1; r < 10; r++) {
            c0 = _mm_aesenc_si128(c0, round_keys[r]);
            c1 = _mm_aesenc_si128(c1, round_keys[r]);
            c2 = _mm_aesenc_si128(c2, round_keys[r]);
            c3 = _mm_aesenc_si128(c3, round_keys[r]);
        }

        // Final Round
        c0 = _mm_aesenclast_si128(c0, round_keys[10]);
        c1 = _mm_aesenclast_si128(c1, round_keys[10]);
        c2 = _mm_aesenclast_si128(c2, round_keys[10]);
        c3 = _mm_aesenclast_si128(c3, round_keys[10]);

        // XOR with Plaintext
        __m128i p0 = _mm_loadu_si128((const __m128i*)(in + offset));
        __m128i p1 = _mm_loadu_si128((const __m128i*)(in + offset + 16));
        __m128i p2 = _mm_loadu_si128((const __m128i*)(in + offset + 32));
        __m128i p3 = _mm_loadu_si128((const __m128i*)(in + offset + 48));

        _mm_storeu_si128((__m128i*)(out + offset),      _mm_xor_si128(p0, c0));
        _mm_storeu_si128((__m128i*)(out + offset + 16), _mm_xor_si128(p1, c1));
        _mm_storeu_si128((__m128i*)(out + offset + 32), _mm_xor_si128(p2, c2));
        _mm_storeu_si128((__m128i*)(out + offset + 48), _mm_xor_si128(p3, c3));

        offset += PARALLEL_BLOCKS * BLOCK_SIZE;
    }

    // Remainder handling
    while (offset < length) {
        __m128i c0 = _mm_xor_si128(ctr0, round_keys[0]);
        for (int r = 1; r < 10; r++) {
            c0 = _mm_aesenc_si128(c0, round_keys[r]);
        }
        c0 = _mm_aesenclast_si128(c0, round_keys[10]);

        uint8_t ks[BLOCK_SIZE];
        _mm_storeu_si128((__m128i*)ks, c0);

        size_t chunk = (length - offset < BLOCK_SIZE) ? (length - offset) : BLOCK_SIZE;
        for (size_t i = 0; i < chunk; i++) {
            out[offset + i] = in[offset + i] ^ ks[i];
        }

        ctr0 = inc_counter(ctr0);
        offset += chunk;
    }

    // Store fence to flush pipeline writes and enforce memory isolation boundaries
    _mm_sfence();

    // Save back IV state
    _mm_storeu_si128((__m128i*)iv, ctr0);
}