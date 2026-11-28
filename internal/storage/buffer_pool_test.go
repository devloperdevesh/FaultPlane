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

// Issue #250: PCIe DMA Channel Binding & Intel DPDK Hardware Bypass Tests

func TestPCIeDMAChannelBinding(t *testing.T) {
	pool := NewBufferPool(32)

	channels := pool.GetDMAChannels()
	if len(channels) < 2 {
		t.Fatalf("expected at least 2 default PCIe DMA channels, got %d", len(channels))
	}

	if channels[0].BandwidthGBs != 64.0 {
		t.Errorf("expected 64.0 GB/s bandwidth, got %f", channels[0].BandwidthGBs)
	}

	newCh, err := pool.BindPCIeDMAChannel("0000:03:00.0 (NVIDIA H100 SXM5)", 128.0)
	if err != nil {
		t.Fatalf("failed to dynamically bind PCIe DMA channel: %v", err)
	}

	if newCh.ChannelID != 2 {
		t.Errorf("expected channel ID 2, got %d", newCh.ChannelID)
	}

	allChannels := pool.GetDMAChannels()
	if len(allChannels) != 3 {
		t.Errorf("expected 3 channels after binding, got %d", len(allChannels))
	}
}

func TestIntelDPDKHardwareBypass(t *testing.T) {
	pool := NewBufferPool(32)

	dpdk := pool.GetDPDKPort()
	if dpdk == nil {
		t.Fatalf("expected DPDK port to be initialized")
	}

	if dpdk.KernelBypassRate != 100.0 {
		t.Errorf("expected 100.0%% kernel bypass rate, got %f", dpdk.KernelBypassRate)
	}

	if !dpdk.HardwareRingActive {
		t.Errorf("expected DPDK hardware ring to be active")
	}

	// Test hot-swapping prefix context bits
	prefix := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04}
	err := pool.HotSwapPrefixContext(0, prefix)
	if err != nil {
		t.Fatalf("failed to hot-swap prefix context: %v", err)
	}

	slot := pool.GetSlot(0)
	if slot.ActiveTokens != 2 {
		t.Errorf("expected 2 active tokens, got %d", slot.ActiveTokens)
	}
}

func TestMoEExpertTensorSubTwoMsLatency(t *testing.T) {
	pool := NewBufferPool(64)

	slot, err := pool.AcquireSlot("org-nemotron-moe-550b", 3)
	if err != nil {
		t.Fatalf("failed to acquire slot: %v", err)
	}

	// Stream 4MB expert tensor chunk
	tensorBytes := uint32(4 * 1024 * 1024)
	res, err := pool.StreamMoEExpertTensor(slot.SlotID, 7, tensorBytes)
	if err != nil {
		t.Fatalf("failed to stream MoE expert tensor: %v", err)
	}

	if !res.SubTwoMsBoundMet {
		t.Fatalf("MoE expert tensor transfer must meet sub-2ms bounds, got %f ms", res.LatencyMs)
	}

	if res.LatencyMs >= 2.0 {
		t.Errorf("latency must be strictly < 2.0 ms, got %f ms", res.LatencyMs)
	}

	if !res.KernelBypassed {
		t.Errorf("expected kernel bypass flag to be true")
	}

	if slot.DMABound != 1 {
		t.Errorf("expected slot to be marked DMA bound, got %d", slot.DMABound)
	}

	if slot.ExpertID != 7 {
		t.Errorf("expected expert ID 7, got %d", slot.ExpertID)
	}

	snap := pool.Snapshot()
	if !snap.PCIeDMABound || !snap.DPDKBypassActive || !snap.SubTwoMsCongruent {
		t.Errorf("snapshot must reflect active PCIe DMA and sub-2ms congruence: %+v", snap)
	}
}

func TestMoEPrefixContextZeroAllocations(t *testing.T) {
	pool := NewBufferPool(64)
	slot, err := pool.AcquireSlot("org-zero-alloc-test", 1)
	if err != nil {
		t.Fatalf("failed to acquire slot: %v", err)
	}

	prefix := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = pool.HotSwapPrefixContext(slot.SlotID, prefix)
	})

	if allocs > 0 {
		t.Errorf("MoE prefix context hot-swap must have zero heap allocations, got %f", allocs)
	}
}

func TestConcurrentMoEDMAStreaming(t *testing.T) {
	pool := NewBufferPool(64)
	var wg sync.WaitGroup

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			slotID := uint32(workerID % 32)
			expertID := uint16(workerID % 8)
			tensorSize := uint32(1024 * 1024 * 2) // 2MB

			res, err := pool.StreamMoEExpertTensor(slotID, expertID, tensorSize)
			if err != nil {
				t.Errorf("worker %d failed DMA stream: %v", workerID, err)
				return
			}
			if !res.SubTwoMsBoundMet {
				t.Errorf("worker %d exceeded sub-2ms bounds: %f ms", workerID, res.LatencyMs)
			}
		}(i)
	}

	wg.Wait()

	snap := pool.Snapshot()
	if snap.TotalSlots != 64 {
		t.Errorf("expected 64 total slots, got %d", snap.TotalSlots)
	}
	if !snap.CacheAlignmentCongruent {
		t.Errorf("cache alignment corrupted under concurrent DMA streams")
	}
}
