// SPDX-License-Identifier: GPL-2.0

#include <linux/bpf.h>
#include <linux/types.h>

#include <bpf/bpf_helpers.h>


char LICENSE[] SEC("license") = "GPL";


#define MAX_BACKENDS 1024
#define PRIMARY_BACKEND 0


/*
 * Active transport sockets.
 *
 * Userspace control-plane manages backend lifecycle.
 */
struct {
	__uint(type, BPF_MAP_TYPE_SOCKMAP);
	__uint(max_entries, MAX_BACKENDS);

	__type(key, __u32);
	__type(value, __u32);

} faultplane_sockmap SEC(".maps");



/*
 * Backend health state.
 *
 * Updated by Go control-plane.
 *
 * 1 = healthy
 * 0 = unavailable
 */
struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, MAX_BACKENDS);

	__type(key, __u32);
	__type(value, __u32);

} faultplane_health SEC(".maps");



/*
 * Register established sockets
 * into kernel transport table.
 */
SEC("sockops")
int faultplane_sockops(
	struct bpf_sock_ops *ctx
)
{
	__u32 key = PRIMARY_BACKEND;


	switch (ctx->op)
	{

	case BPF_SOCK_OPS_ACTIVE_ESTABLISHED_CB:
	case BPF_SOCK_OPS_PASSIVE_ESTABLISHED_CB:

		bpf_sock_map_update(
			ctx,
			&faultplane_sockmap,
			&key,
			BPF_ANY
		);

		break;


	default:
		break;
	}


	return 0;
}



/*
 * Zero-copy kernel socket redirect.
 *
 * Decision happens inside kernel
 * using backend health state.
 */
SEC("sk_msg")
int faultplane_redirect(
	struct sk_msg_md *msg
)
{
	__u32 backend = PRIMARY_BACKEND;


	__u32 *healthy =
		bpf_map_lookup_elem(
			&faultplane_health,
			&backend
		);


	if (!healthy)
	{
		return SK_PASS;
	}


	if (*healthy == 0)
	{
		return SK_PASS;
	}


	return bpf_msg_redirect_map(
		msg,
		&faultplane_sockmap,
		backend,
		0
	);
}