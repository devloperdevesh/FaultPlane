package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
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
	capacities := []int{1, 2, 3, 7, 16, 63, 64, 128}
	for _, cap := range capacities {
		for trial := 0; trial < 10; trial++ {
			pool := NewBufferPool(cap)

			if !pool.VerifyCacheAlignment() {
				t.Fatalf("cap=%d trial=%d: VerifyCacheAlignment returned false", cap, trial)
			}

			baseAddr := uintptr(unsafe.Pointer(&pool.slots[0]))
			if baseAddr%CacheLineSize != 0 {
				t.Fatalf("cap=%d trial=%d: base address 0x%X is not 64-byte aligned", cap, trial, baseAddr)
			}

			for i := 0; i < len(pool.slots); i++ {
				addr := uintptr(unsafe.Pointer(&pool.slots[i]))
				if addr%CacheLineSize != 0 {
					t.Errorf("cap=%d slot=%d address 0x%X is not 64-byte aligned", cap, i, addr)
				}
				if diff := addr - baseAddr; diff != uintptr(i*CacheLineSize) {
					t.Errorf("cap=%d slot=%d offset %d != expected %d", cap, i, diff, i*CacheLineSize)
				}
			}

			snap := pool.Snapshot()
			if !snap.CacheAlignmentCongruent {
				t.Errorf("cap=%d: snapshot CacheAlignmentCongruent must be true", cap)
			}
			for _, s := range snap.Slots {
				if !s.CacheLineAligned {
					t.Errorf("cap=%d slot=%d: snapshot CacheLineAligned must be true", cap, s.Index)
				}
			}
		}
	}
}

func TestVerifyCacheAlignmentRejectsUnalignedBaseAndEmpty(t *testing.T) {
	empty := &BufferPool{}
	if empty.VerifyCacheAlignment() {
		t.Fatalf("empty BufferPool must not report aligned")
	}

	// Construct a backing buffer with guaranteed misaligned 64B start (+8 bytes)
	// even though consecutive slots remain 64 bytes apart.
	raw := make([]byte, 4*CacheLineSize+16)
	ptr := unsafe.Pointer(&raw[0])
	if uintptr(ptr)%CacheLineSize == 0 {
		ptr = unsafe.Add(ptr, 8)
	}

	unalignedSlots := unsafe.Slice((*CacheLineSlot)(ptr), 2)
	for i := range unalignedSlots {
		unalignedSlots[i].SlotID = uint32(i)
	}
	unaligned := &BufferPool{
		backing:   raw,
		slots:     unalignedSlots,
		slotCount: 2,
	}

	if unaligned.VerifyCacheAlignment() {
		t.Fatalf("VerifyCacheAlignment must return false when base address 0x%X is not 64-byte aligned", uintptr(ptr))
	}

	snap := unaligned.Snapshot()
	if snap.CacheAlignmentCongruent {
		t.Fatalf("Snapshot.CacheAlignmentCongruent must be false for misaligned base")
	}
	if snap.Slots[0].CacheLineAligned {
		t.Fatalf("Snapshot.Slots[0].CacheLineAligned must be false when address 0x%X is unaligned", uintptr(ptr))
	}
}

func TestBufferPoolAcquireAndRelease(t *testing.T) {
	pool := NewBufferPool(16)

	slot, err := pool.AcquireSlot("org-finance-01", 1)
	if err != nil {
		t.Fatalf("failed to acquire slot: %v", err)
	}

	if slot.Status() != SlotStatusActive {
		t.Errorf("expected slot to be active, got %d", slot.Status())
	}

	if slot.GetDomainKey() != 1 {
		t.Errorf("expected domain key 1, got %d", slot.GetDomainKey())
	}

	expectedHash := sha256.Sum256([]byte("org-finance-01"))
	if slot.TenantHash() != expectedHash {
		t.Errorf("expected tenant hash %x, got %x", expectedHash, slot.TenantHash())
	}

	ok := pool.ReleaseSlot(slot.SlotID)
	if !ok {
		t.Fatalf("failed to release slot")
	}

	if slot.Status() != SlotStatusFree {
		t.Errorf("expected slot to be free, got %d", slot.Status())
	}
	if slot.GetDomainKey() != 0 {
		t.Errorf("expected cleared domain key 0 after release, got %d", slot.GetDomainKey())
	}
	if slot.TenantHash() != ([32]byte{}) {
		t.Errorf("expected cleared tenant hash after release, got %x", slot.TenantHash())
	}
	if slot.ActiveTokenCount() != 0 {
		t.Errorf("expected cleared active tokens after release, got %d", slot.ActiveTokenCount())
	}
	if pool.ReleaseSlot(slot.SlotID) {
		t.Errorf("releasing an already free slot should return false")
	}
}

func TestBufferPoolFencingStateTransitions(t *testing.T) {
	pool := NewBufferPool(8)

	// Cannot fence out-of-range or free slots.
	if pool.FenceSlot(999) {
		t.Fatalf("expected FenceSlot on invalid ID to fail")
	}
	if pool.FenceSlot(0) {
		t.Fatalf("expected FenceSlot on free slot to fail")
	}
	if pool.GetSlot(0).Status() != SlotStatusFree {
		t.Fatalf("expected slot 0 to remain Free after failed fence")
	}

	slot, err := pool.AcquireSlot("org-defense-02", 2)
	if err != nil {
		t.Fatalf("failed to acquire slot: %v", err)
	}

	if !pool.FenceSlot(slot.SlotID) {
		t.Fatalf("expected FenceSlot on active slot to succeed")
	}
	if slot.Status() != SlotStatusFenced {
		t.Errorf("expected slot to be fenced, got %d", slot.Status())
	}
	if pool.FenceSlot(slot.SlotID) {
		t.Errorf("expected FenceSlot on already-fenced slot to return false")
	}
	if pool.ReleaseSlot(slot.SlotID) {
		t.Errorf("expected ReleaseSlot on fenced slot to return false until unfenced")
	}

	snap := pool.Snapshot()
	if snap.FencedSlots != 1 {
		t.Errorf("expected 1 fenced slot in snapshot, got %d", snap.FencedSlots)
	}

	if !pool.UnfenceSlot(slot.SlotID) {
		t.Fatalf("expected UnfenceSlot on fenced slot to succeed")
	}
	if pool.UnfenceSlot(slot.SlotID) {
		t.Errorf("expected second UnfenceSlot on active slot to return false")
	}
	if !pool.ReleaseSlot(slot.SlotID) {
		t.Fatalf("expected ReleaseSlot after UnfenceSlot to succeed")
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

func TestBufferPoolConcurrentAcquireReleaseFenceSnapshot(t *testing.T) {
	pool := NewBufferPool(16)
	var wg sync.WaitGroup

	tenants := []struct {
		id     string
		domain uint8
	}{
		{"tenant-A", 1},
		{"tenant-B", 2},
		{"tenant-C", 3},
		{"tenant-D", 4},
	}

	expectedPrefixByDomain := make(map[uint8]string, len(tenants))
	expectedHashByDomain := make(map[uint8][32]byte, len(tenants))
	for _, tc := range tenants {
		sum := sha256.Sum256([]byte(tc.id))
		expectedHashByDomain[tc.domain] = sum
		expectedPrefixByDomain[tc.domain] = hex.EncodeToString(sum[:])[:12]
	}

	const workers = 16
	const iters = 200
	var done atomic.Bool

	// Continuous Snapshot readers verifying consistency under concurrent churn.
	var readerWG sync.WaitGroup
	for r := 0; r < 4; r++ {
		readerWG.Add(1)
		go func() {
			defer readerWG.Done()
			for !done.Load() {
				snap := pool.Snapshot()
				if !snap.CacheAlignmentCongruent {
					t.Errorf("cache alignment must remain congruent under concurrency")
				}
				seenActive := 0
				seenFenced := 0
				for _, s := range snap.Slots {
					switch s.Status {
					case "ACTIVE", "FENCED":
						if s.Status == "ACTIVE" {
							seenActive++
						} else {
							seenFenced++
						}
						wantPrefix, ok := expectedPrefixByDomain[s.DomainKey]
						if !ok {
							t.Errorf("slot %d observed unexpected domain key %d with status %s", s.SlotID, s.DomainKey, s.Status)
							continue
						}
						if s.TenantID != wantPrefix {
							t.Errorf("torn snapshot read on slot %d: status=%s domain=%d tenant=%s want=%s", s.SlotID, s.Status, s.DomainKey, s.TenantID, wantPrefix)
						}
					case "FREE":
						if s.DomainKey != 0 || s.TenantID != "000000000000" || s.ActiveTokens != 0 {
							t.Errorf("unclean FREE slot %d in snapshot: domain=%d tenant=%s tokens=%d", s.SlotID, s.DomainKey, s.TenantID, s.ActiveTokens)
						}
					default:
						t.Errorf("transient status %q leaked into snapshot on slot %d", s.Status, s.SlotID)
					}
				}
				if seenActive != snap.ActiveSlots || seenFenced != snap.FencedSlots {
					t.Errorf("snapshot counts mismatch: got active=%d fenced=%d, counted active=%d fenced=%d", snap.ActiveSlots, snap.FencedSlots, seenActive, seenFenced)
				}
			}
		}()
	}

	// Concurrent writers acquiring, checking, fencing, unfencing, and releasing slots.
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			tc := tenants[workerID%len(tenants)]
			wantHash := expectedHashByDomain[tc.domain]

			for iter := 0; iter < iters; iter++ {
				slot, err := pool.AcquireSlot(tc.id, tc.domain)
				if err != nil {
					runtime.Gosched()
					continue
				}

				if got := slot.GetDomainKey(); got != tc.domain {
					t.Errorf("worker %d slot %d domain mismatch: got %d want %d", workerID, slot.SlotID, got, tc.domain)
				}
				if gotHash := slot.TenantHash(); gotHash != wantHash {
					t.Errorf("worker %d slot %d tenant hash mismatch", workerID, slot.SlotID)
				}

				if (workerID+iter)%5 == 0 {
					if pool.FenceSlot(slot.SlotID) {
						if slot.Status() != SlotStatusFenced {
							t.Errorf("expected slot %d fenced", slot.SlotID)
						}
						pool.UnfenceSlot(slot.SlotID)
					}
				}

				if !pool.ReleaseSlot(slot.SlotID) {
					t.Errorf("worker %d failed to release owned slot %d", workerID, slot.SlotID)
				}
			}
		}(w)
	}

	wg.Wait()
	done.Store(true)
	readerWG.Wait()

	finalSnap := pool.Snapshot()
	if !finalSnap.CacheAlignmentCongruent {
		t.Errorf("cache alignment must remain congruent after concurrency")
	}
	if finalSnap.ActiveSlots != 0 || finalSnap.FencedSlots != 0 {
		t.Errorf("expected all slots released cleanly, got active=%d fenced=%d", finalSnap.ActiveSlots, finalSnap.FencedSlots)
	}
}

func TestBufferPoolConcurrentFenceVsReleaseMutualExclusion(t *testing.T) {
	pool := NewBufferPool(4)

	for iter := 0; iter < 500; iter++ {
		slot, err := pool.AcquireSlot("tenant-race", 1)
		if err != nil {
			t.Fatalf("iter %d: acquire failed: %v", iter, err)
		}

		var subWG sync.WaitGroup
		var fenced, released atomic.Bool
		subWG.Add(2)
		go func() {
			defer subWG.Done()
			fenced.Store(pool.FenceSlot(slot.SlotID))
		}()
		go func() {
			defer subWG.Done()
			released.Store(pool.ReleaseSlot(slot.SlotID))
		}()
		subWG.Wait()

		if fenced.Load() && released.Load() {
			t.Fatalf("iter %d: slot %d was simultaneously fenced and released", iter, slot.SlotID)
		}
		if !fenced.Load() && !released.Load() {
			t.Fatalf("iter %d: neither FenceSlot nor ReleaseSlot succeeded on active slot %d", iter, slot.SlotID)
		}
		if fenced.Load() {
			if !pool.UnfenceSlot(slot.SlotID) {
				t.Fatalf("iter %d: failed to unfence slot %d", iter, slot.SlotID)
			}
			if !pool.ReleaseSlot(slot.SlotID) {
				t.Fatalf("iter %d: failed to release unfenced slot %d", iter, slot.SlotID)
			}
		}
	}
}

func TestBufferPoolTelemetryAccuracy(t *testing.T) {
	pool := NewBufferPool(8)

	var before debug.GCStats
	debug.ReadGCStats(&before)

	// Ensure at least one runtime GC cycle is recorded cleanly.
	sink := make([][]byte, 256)
	for i := range sink {
		sink[i] = make([]byte, 8192)
	}
	runtime.KeepAlive(sink)
	runtime.GC()

	var after debug.GCStats
	debug.ReadGCStats(&after)

	snap := pool.Snapshot()
	if !snap.GCPausesAvailable {
		t.Errorf("expected GCPausesAvailable to be true")
	}
	if snap.GCPauseCount <= before.NumGC || snap.GCPauseCount < after.NumGC {
		t.Errorf("expected GCPauseCount to advance with runtime GCs (before=%d, after=%d, snap=%d)", before.NumGC, after.NumGC, snap.GCPauseCount)
	}
	expectedMinMs := float64(after.PauseTotal.Nanoseconds()) / 1e6
	if snap.GCPausesMs < 0 || snap.GCPausesMs < expectedMinMs {
		t.Errorf("expected GCPausesMs >= %f (measured total pause), got %f", expectedMinMs, snap.GCPausesMs)
	}

	if snap.L1L2ThrashingAvailable {
		t.Errorf("expected L1L2ThrashingAvailable to be false without PMU counters")
	}
	if snap.L1L2ThrashingRate != nil {
		t.Errorf("expected nil L1L2ThrashingRate when unavailable, got %v", *snap.L1L2ThrashingRate)
	}
	if snap.L1L2ThrashingNote == "" {
		t.Errorf("expected explanatory L1L2ThrashingNote when unavailable")
	}
	if snap.HardwareMPKEnforced {
		t.Errorf("expected HardwareMPKEnforced to be false in userspace")
	}
	if snap.IsolationMode != IsolationModeLogicalTag {
		t.Errorf("expected IsolationMode %q, got %q", IsolationModeLogicalTag, snap.IsolationMode)
	}
}

func TestSlotSnapshotFormatStable(t *testing.T) {
	pool := NewBufferPool(4)
	slot, err := pool.AcquireSlot("tenant-format", 3)
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	snap := pool.Snapshot()
	var found bool
	for _, s := range snap.Slots {
		if s.SlotID == slot.SlotID {
			found = true
			if s.Status != "ACTIVE" || s.DomainKey != 3 {
				t.Errorf("unexpected slot state: %+v", s)
			}
			expectedOffset := fmt.Sprintf("0x%04X", int(slot.SlotID)*CacheLineSize)
			if s.MemoryOffsetHex != expectedOffset {
				t.Errorf("expected offset %s, got %s", expectedOffset, s.MemoryOffsetHex)
			}
		}
	}
	if !found {
		t.Fatalf("acquired slot %d missing from snapshot", slot.SlotID)
	}
}
