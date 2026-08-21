// SPDX-License-Identifier: GPL-2.0

#include <linux/types.h>
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

SEC("xdp")
int faultplane_ingress_filter(struct xdp_md *ctx)
{
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;

    if (data + sizeof(__u32) > data_end) {
        return XDP_PASS;
    }

    bpf_printk("FaultPlane XDP ingress filter active");

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";