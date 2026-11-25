package storage

import (
	"sync"
	"testing"
	"unsafe"
)

func TestCacheLineSlotSizeIsExactly64Bytes(t *testing.T) {
	size := unsafe.Sizeof(CacheLineSlot{})
	if size != 64 {
		t.Fatalf("CacheLineSlot must be exactly 64 bytes to match CPU cache lines, got %d", size)
	}
}

func TestBufferPoolCacheAlignment(t *testing.T) {
	pool := NewBufferPool(64)

	if !pool.VerifyCacheAlignment() {
		t.Fatalf("BufferPool slots must all be separated across 64-byte structural boundaries")
	}

	baseAddr := uintptr(unsafe.Pointer(&pool.slots[0]))
	for i := 0; i < len(pool.slots); i++ {
		addr := uintptr(unsafe.Pointer(&pool.slots[i]))
		diff := addr - baseAddr
		if diff%64 != 0 {
			t.Errorf("Slot %d offset %d is not a multiple of 64 bytes", i, diff)
		}
	}
}

func TestBufferPoolAcquireAndRelease(t *testing.T) {
	pool := NewBufferPool(16)

	slot, err := pool.AcquireSlot("org-finance-01", 1)
	if err != nil {
		t.Fatalf("failed to acquire slot: %v", err)
	}

	if slot.Flags != SlotStatusActive {
		t.Errorf("expected slot to be active, got %d", slot.Flags)
	}

	if slot.DomainKey != 1 {
		t.Errorf("expected domain key 1, got %d", slot.DomainKey)
	}

	ok := pool.ReleaseSlot(slot.SlotID)
	if !ok {
		t.Fatalf("failed to release slot")
	}

	if slot.Flags != SlotStatusFree {
		t.Errorf("expected slot to be free, got %d", slot.Flags)
	}
}

func TestBufferPoolFencing(t *testing.T) {
	pool := NewBufferPool(8)

	slot, err := pool.AcquireSlot("org-defense-02", 2)
	if err != nil {
		t.Fatalf("failed to acquire slot: %v", err)
	}

	pool.FenceSlot(slot.SlotID)
	if slot.Flags != SlotStatusFenced {
		t.Errorf("expected slot to be fenced, got %d", slot.Flags)
	}

	snap := pool.Snapshot()
	if snap.FencedSlots != 1 {
		t.Errorf("expected 1 fenced slot in snapshot, got %d", snap.FencedSlots)
	}
}

func TestBufferPoolZeroAllocationHotPath(t *testing.T) {
	pool := NewBufferPool(64)

	// Pre-acquire
	slot, err := pool.AcquireSlot("org-benchmark-01", 1)
	if err != nil {
		t.Fatalf("failed to acquire: %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = pool.GetSlot(slot.SlotID)
	})

	if allocs > 0 {
		t.Errorf("hot-path slot retrieval must have 0 heap allocations, got %f", allocs)
	}
}

func TestBufferPoolConcurrentMultiTenant(t *testing.T) {
	pool := NewBufferPool(128)
	var wg sync.WaitGroup

	tenants := []string{"tenant-A", "tenant-B", "tenant-C", "tenant-D"}

	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			tenant := tenants[workerID%len(tenants)]
			domain := uint8((workerID % 4) + 1)

			slot, err := pool.AcquireSlot(tenant, domain)
			if err == nil {
				_ = pool.GetSlot(slot.SlotID)
				pool.ReleaseSlot(slot.SlotID)
			}
		}(i)
	}

	wg.Wait()

	snap := pool.Snapshot()
	if !snap.CacheAlignmentCongruent {
		t.Errorf("cache alignment must remain congruent under concurrency")
	}
}
