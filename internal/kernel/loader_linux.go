package kernel

import (
	"fmt"
	"sync"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

const defaultCgroupPath = "/sys/fs/cgroup"

type Loader struct {
	mu sync.Mutex

	collection *ebpf.Collection
	sockops    link.Link
	skmsg      *link.RawLink

	enforcement *ebpf.Map
	loaded      bool
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(objectPath string) error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("remove eBPF memlock: %w", err)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.loaded {
		return fmt.Errorf("eBPF collection already loaded")
	}

	spec, err := ebpf.LoadCollectionSpec(objectPath)
	if err != nil {
		return fmt.Errorf("load eBPF spec: %w", err)
	}

	var objs struct {
		SockOps     *ebpf.Program `ebpf:"faultplane_sockops"`
		Redirect    *ebpf.Program `ebpf:"faultplane_redirect"`
		SockMap     *ebpf.Map     `ebpf:"faultplane_sockmap"`
		Enforcement *ebpf.Map     `ebpf:"faultplane_enforcement"`
	}

	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		return fmt.Errorf("load eBPF objects: %w", err)
	}

	closeObjects := func() {
		if objs.Enforcement != nil {
			_ = objs.Enforcement.Close()
		}

		if objs.SockMap != nil {
			_ = objs.SockMap.Close()
		}
		if objs.Redirect != nil {
			_ = objs.Redirect.Close()
		}
		if objs.SockOps != nil {
			_ = objs.SockOps.Close()
		}
	}

	sockops, err := link.AttachCgroup(link.CgroupOptions{
		Path:    defaultCgroupPath,
		Program: objs.SockOps,
		Attach:  ebpf.AttachCGroupSockOps,
	})
	if err != nil {
		closeObjects()
		return fmt.Errorf("attach sockops: %w", err)
	}

	skmsg, err := link.AttachRawLink(link.RawLinkOptions{
		Target:  objs.SockMap.FD(),
		Program: objs.Redirect,
		Attach:  ebpf.AttachSkMsgVerdict,
	})
	if err != nil {
		_ = sockops.Close()
		closeObjects()
		return fmt.Errorf("attach sk_msg: %w", err)
	}

	l.collection = &ebpf.Collection{
		Programs: map[string]*ebpf.Program{
			"faultplane_sockops":  objs.SockOps,
			"faultplane_redirect": objs.Redirect,
		},
		Maps: map[string]*ebpf.Map{
			"faultplane_sockmap":     objs.SockMap,
			"faultplane_enforcement": objs.Enforcement,
		},
	}

	l.sockops = sockops
	l.skmsg = skmsg
	l.enforcement = objs.Enforcement
	l.loaded = true

	return nil
}

func (l *Loader) Loaded() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.loaded
}

func (l *Loader) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.loaded {
		return nil
	}

	var firstErr error

	if l.skmsg != nil {
		if err := l.skmsg.Close(); err != nil {
			firstErr = fmt.Errorf("detach sk_msg: %w", err)
		}
		l.skmsg = nil
	}

	if l.sockops != nil {
		if err := l.sockops.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("detach sockops: %w", err)
		}
		l.sockops = nil
	}

	if l.collection != nil {
		l.collection.Close()
		l.collection = nil
	}
	l.enforcement = nil

	l.loaded = false

	return firstErr
}

const (
	EnforcementAllow    uint32 = 0
	EnforcementRedirect uint32 = 1
	EnforcementDrop     uint32 = 2
)

// SetEnforcementMode updates the live eBPF enforcement control map.
func (l *Loader) SetEnforcementMode(mode uint32) error {
	if mode > EnforcementDrop {
		return fmt.Errorf("invalid enforcement mode: %d", mode)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.loaded {
		return fmt.Errorf("eBPF collection is not loaded")
	}

	if l.enforcement == nil {
		return fmt.Errorf("eBPF enforcement map is unavailable")
	}

	key := uint32(0)

	if err := l.enforcement.Update(&key, &mode, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("update enforcement mode: %w", err)
	}

	return nil
}

// EnforcementMode returns the current live eBPF enforcement mode.
func (l *Loader) EnforcementMode() (uint32, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.loaded {
		return 0, fmt.Errorf("eBPF collection is not loaded")
	}

	if l.enforcement == nil {
		return 0, fmt.Errorf("eBPF enforcement map is unavailable")
	}

	key := uint32(0)
	var mode uint32

	if err := l.enforcement.Lookup(&key, &mode); err != nil {
		return 0, fmt.Errorf("read enforcement mode: %w", err)
	}

	return mode, nil
}
