# FaultPlane eBPF Architecture


## Overview

FaultPlane uses eBPF sockmap programs
to perform kernel-level socket routing.


Architecture:


Application
    |
    v
Go Control Plane
    |
    v
BPF Loader
    |
    v
eBPF Program
    |
    v
Linux Kernel


## Components


### sockmap.bpf.c

Responsible for:

- socket lifecycle hooks
- sockmap registration
- kernel redirect path


### loader.go

Responsible for:

- loading BPF object
- managing lifecycle
- cleanup


## Requirements

Linux kernel with eBPF support.

Required tools:

- clang
- llvm
- libbpf
- bpftool


## Verification

Before production deployment:

- compile BPF object
- validate verifier acceptance
- test unsupported kernel behaviour