// SPDX-License-Identifier: GPL-2.0

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>


char LICENSE[] SEC("license") = "GPL";


// Socket map used for redirecting TCP streams.
struct {
	__uint(type, BPF_MAP_TYPE_SOCKMAP);
	__uint(max_entries, 1024);
	__type(key, __u32);
	__type(value, __u32);
}
sock_map SEC(".maps");


// Attach to socket operations.
SEC("sockops")
int faultplane_sockops(
	struct bpf_sock_ops *ctx
)
{

	switch (ctx->op) {

	case BPF_SOCK_OPS_PASSIVE_ESTABLISHED_CB:

		// Future:
		// register healthy backend sockets

		break;


	case BPF_SOCK_OPS_ACTIVE_ESTABLISHED_CB:

		// Future:
		// track outbound connections

		break;

	default:

		break;
	}


	return 0;
}



// Redirect socket messages.
SEC("sk_msg")
int faultplane_redirect(
	struct sk_msg_md *msg
)
{

	int key = 0;


	return bpf_msg_redirect_map(
		msg,
		&sock_map,
		key,
		0
	);
}