// SPDX-License-Identifier: GPL-2.0

#include <linux/bpf.h>
#include <linux/types.h>

#include <bpf/bpf_helpers.h>


char LICENSE[] SEC("license") = "GPL";


#define MAX_BACKENDS 1024
#define DEFAULT_BACKEND_KEY 0


/*
 * Sockmap used by FaultPlane transport layer.
 *
 * Flow:
 *
 * Application socket
 *        |
 *        v
 *     sock_map
 *        |
 *        v
 * Healthy backend socket
 *
 */
struct
{
	__uint(type, BPF_MAP_TYPE_SOCKMAP);
	__uint(max_entries, MAX_BACKENDS);

	__type(key, __u32);
	__type(value, __u32);

} faultplane_sockmap SEC(".maps");



/*
 * Socket lifecycle hook.
 *
 * Used for:
 * - observing TCP socket creation
 * - registering backend sockets
 * - attaching sockets into sockmap
 *
 * Actual backend selection happens
 * from userspace control-plane.
 */
SEC("sockops")
int faultplane_sockops(
	struct bpf_sock_ops *ctx
)
{

	switch (ctx->op) {


	case BPF_SOCK_OPS_ACTIVE_ESTABLISHED_CB:

		/*
		 * Active connection established.
		 *
		 * Future:
		 * - attach socket metadata
		 * - register connection state
		 */
		break;



	case BPF_SOCK_OPS_PASSIVE_ESTABLISHED_CB:

		/*
		 * Incoming connection established.
		 *
		 * Future:
		 * - validate backend health
		 * - update routing state
		 */
		break;



	default:

		break;
	}


	return 0;
}



/*
 * Kernel-level socket redirect.
 *
 * Traffic does not return to userspace.
 *
 * Kernel:
 *
 * socket
 *   |
 *   v
 * sockmap
 *   |
 *   v
 * backend socket
 *
 */
SEC("sk_msg")
int faultplane_redirect(
	struct sk_msg_md *msg
)
{

	__u32 backend_key = DEFAULT_BACKEND_KEY;


	return bpf_msg_redirect_map(
		msg,
		&faultplane_sockmap,
		backend_key,
		0
	);
}