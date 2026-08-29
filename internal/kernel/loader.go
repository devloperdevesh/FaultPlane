package kernel

import (
	"fmt"
	"sync"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

const defaultCgroupPath = "/sys/fs/cgroup"

type Loader struct {
	mu sync.Mutex

	collection *ebpf.Collection
	sockops    link.Link

	skmsgTarget   int
	skmsgProgram  *ebpf.Program
	skmsgAttached bool

	loaded bool
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(objectPath string) error {
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
		SockOps  *ebpf.Program `ebpf:"faultplane_sockops"`
		Redirect *ebpf.Program `ebpf:"faultplane_redirect"`
		SockMap  *ebpf.Map     `ebpf:"faultplane_sockmap"`
	}

	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		return fmt.Errorf("load eBPF objects: %w", err)
	}

	closeObjects := func() {
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

	if err := link.RawAttachProgram(link.RawAttachProgramOptions{
		Target:  objs.SockMap.FD(),
		Program: objs.Redirect,
		Attach:  ebpf.AttachSkMsgVerdict,
	}); err != nil {
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
			"faultplane_sockmap": objs.SockMap,
		},
	}

	l.sockops = sockops
	l.skmsgTarget = objs.SockMap.FD()
	l.skmsgProgram = objs.Redirect
	l.skmsgAttached = true
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

	if l.skmsgAttached && l.skmsgProgram != nil {
		if err := link.RawDetachProgram(link.RawDetachProgramOptions{
			Target:  l.skmsgTarget,
			Program: l.skmsgProgram,
			Attach:  ebpf.AttachSkMsgVerdict,
		}); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("detach sk_msg: %w", err)
		}

		l.skmsgAttached = false
		l.skmsgTarget = 0
		l.skmsgProgram = nil
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

	l.loaded = false

	return firstErr
}
