package storage

import (
	"errors"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// CacheLineSize defines the standard L1 cache line size (64 bytes for x86-64/ARM64).
const CacheLineSize = 64

// Buffer represents a fixed-capacity contiguous byte slice wrapper used in storage engine IO.
type Buffer struct {
	ID   uint64
	Data []byte
}

// node represents a lock-free stack entry wrapping a Buffer allocation.
type node struct {
	buf  *Buffer
	next unsafe.Pointer // *node
}

// cachePad provides explicit 64-byte memory fence alignment to prevent false sharing.
type cachePad struct {
	_ [CacheLineSize]byte
}

// LockFreeBufferPool implements a high-throughput, zero-lock storage buffer pool.
type LockFreeBufferPool struct {
	_    cachePad
	head unsafe.Pointer // *node atomic head
	_    cachePad
	tail unsafe.Pointer // *node atomic tail (for ring/queue semantics)
	_    cachePad
	
	capacity   uint64
	bufferSize int
	allocated  atomic.Uint64
	
	_ cachePad
}

var (
	ErrPoolEmpty = errors.New("buffer pool exhausted: no free buffers available")
	ErrPoolFull  = errors.New("buffer pool is full")
)

// NewLockFreeBufferPool initializes a cache-aligned lock-free buffer pool.
func NewLockFreeBufferPool(capacity int, bufferSize int) *LockFreeBufferPool {
	pool := &LockFreeBufferPool{
		capacity:   uint64(capacity),
		bufferSize: bufferSize,
	}

	// Pre-allocate buffer capacity to eliminate runtime heap allocation during CAS acquire
	for i := 0; i < capacity; i++ {
		buf := &Buffer{
			ID:   uint64(i),
			Data: make([]byte, bufferSize),
		}
		pool.pushNode(&node{buf: buf})
		pool.allocated.Add(1)
	}

	return pool
}

// Acquire retrieves a storage buffer using hardware Atomic CAS loops without mutexes.
func (p *LockFreeBufferPool) Acquire() (*Buffer, error) {
	for {
		headPtr := atomic.LoadPointer(&p.head)
		if headPtr == nil {
			return nil, ErrPoolEmpty
		}

		headNode := (*node)(headPtr)
		nextPtr := atomic.LoadPointer(&headNode.next)

		// Hardware CAS loop to safely decouple head node
		if atomic.CompareAndSwapPointer(&p.head, headPtr, nextPtr) {
			// Isolate node references
			buf := headNode.buf
			headNode.next = nil
			return buf, nil
		}

		// Yield processor hints under extreme contention to reduce CPU bus locking
		runtime.Gosched()
	}
}

// Release returns a buffer back to the lock-free pool using Atomic CAS.
func (p *LockFreeBufferPool) Release(buf *Buffer) bool {
	if buf == nil || len(buf.Data) != p.bufferSize {
		return false
	}

	// Clear buffer slice contents for safety without reallocating
	clear(buf.Data)

	newNode := &node{buf: buf}

	for {
		headPtr := atomic.LoadPointer(&p.head)
		newNode.next = headPtr

		// Hardware CAS push onto free-list
		if atomic.CompareAndSwapPointer(&p.head, headPtr, unsafe.Pointer(newNode)) {
			return true
		}

		runtime.Gosched()
	}
}

// pushNode is an internal helper for pool initialization.
func (p *LockFreeBufferPool) pushNode(n *node) {
	for {
		headPtr := atomic.LoadPointer(&p.head)
		n.next = headPtr
		if atomic.CompareAndSwapPointer(&p.head, headPtr, unsafe.Pointer(n)) {
			break
		}
	}
}

// Capacity returns the maximum buffer capacity.
func (p *LockFreeBufferPool) Capacity() uint64 {
	return p.capacity
}

// Available returns the active count of available buffers via atomic load.
func (p *LockFreeBufferPool) Available() uint64 {
	var count uint64
	currPtr := atomic.LoadPointer(&p.head)
	for currPtr != nil {
		count++
		currNode := (*node)(currPtr)
		currPtr = atomic.LoadPointer(&currNode.next)
	}
	return count
}