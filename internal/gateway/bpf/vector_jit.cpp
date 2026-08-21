// SPDX-License-Identifier: Apache-2.0

#include <cstdint>
#include <immintrin.h>

extern "C" {

/*
 * Compares 16 IPv4/route identifiers against one target value
 * using AVX-512 integer comparison.
 *
 * Returns a 16-bit match mask:
 * bit N is set when active_cluster_ips[N] == target_packet_ip.
 *
 * The caller must provide at least 16 uint32_t elements.
 */
std::uint32_t
faultplane_vector_match_512(const std::uint32_t *active_cluster_ips,
                            std::uint32_t target_packet_ip) {
  if (active_cluster_ips == nullptr) {
    return 0;
  }

  const __m512i target = _mm512_set1_epi32(static_cast<int>(target_packet_ip));

  const __m512i clusters = _mm512_loadu_si512(active_cluster_ips);

  const __mmask16 matches = _mm512_cmpeq_epi32_mask(clusters, target);

  return static_cast<std::uint32_t>(matches);
}

} // extern "C"