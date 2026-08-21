// SPDX-License-Identifier: Apache-2.0

#ifndef FAULTPLANE_VECTOR_JIT_H
#define FAULTPLANE_VECTOR_JIT_H

#include <cstdint>

#ifdef __cplusplus
extern "C" {
#endif

std::uint32_t
faultplane_vector_match_512(const std::uint32_t *active_cluster_ips,
                            std::uint32_t target_packet_ip);

#ifdef __cplusplus
}
#endif

#endif // FAULTPLANE_VECTOR_JIT_H