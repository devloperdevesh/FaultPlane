package storage

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"unsafe"
)

const (
	// CacheLineSize defines the CPU L1/L2 cache line width (64 bytes).
	CacheLineSize = 64

	// DefaultSlotCount specifies the standard multi-tenant ring capacity.
	DefaultSlotCount = 64

	// SlotStatusFree indicates an unassigned slot whose metadata has been cleared.
	SlotStatusFree uint32 = 0
	// SlotStatusActive indicates an initialized slot assigned to a tenant.
	SlotStatusActive uint32 = 1
	// SlotStatusFenced indicates an active slot isolated under a logical software fence.
	SlotStatusFenced uint32 = 2
	// SlotStatusInitializing is a transitory state while AcquireSlot writes slot metadata.
	SlotStatusInitializing uint32 = 3
	// SlotStatusReleasing is a transitory state while ReleaseSlot clears slot metadata.
	SlotStatusReleasing uint32 = 4

	// IsolationModeLogicalTag indicates that tenant isolation and domain keys
	// represent logical software metadata tags rather than hardware Intel MPK/PKRU enforcement.
	IsolationModeLogicalTag = "logical_domain_tag"
)

var (
	ErrPoolExhausted = errors.New("buffer pool slots exhausted")
	ErrInvalidSlotID = errors.New("invalid slot id")
)

// CacheLineSlot is a multi-tenant queue slot padded to exactly 64 bytes so adjacent
// slots occupy separate CPU cache lines, mitigating false sharing between cores.
// All mutable fields are accessed exclusively through atomic operations and a seqlock
// sequence counter so concurrent writers and Snapshot readers remain race-free.
type CacheLineSlot struct {
	accessCount  uint64    // 8 bytes (offset 0..7): Atomic acquisition counter
	SlotID       uint32    // 4 bytes (offset 8..11): Immutable slot identifier
	activeTokens uint32    // 4 bytes (offset 12..15): Atomic in-flight token count
	tenantWords  [4]uint64 // 32 bytes (offset 16..47): Atomic SHA-256 tenant workspace digest
	flags        uint32    // 4 bytes (offset 48..51): Atomic SlotStatus flags
	seq          uint32    // 4 bytes (offset 52..55): Atomic seqlock counter (odd while mutating)
	domainKey    uint32    // 4 bytes (offset 56..59): Logical protection domain tag (low 8 bits)
	_            uint32    // 4 bytes (offset 60..63): Padding to reach exactly 64 bytes
}

// Compile-time assertion that CacheLineSlot is exactly 64 bytes.
var _ [CacheLineSize]byte = [unsafe.Sizeof(CacheLineSlot{})]byte{}

// Status atomically returns the slot's current lifecycle state.
func (s *CacheLineSlot) Status() uint32 {
	if s == nil {
		return SlotStatusFree
	}
	return atomic.LoadUint32(&s.flags)
}

// GetDomainKey returns a consistent snapshot of the slot's logical domain tag.
func (s *CacheLineSlot) GetDomainKey() uint8 {
	if s == nil {
		return 0
	}
	_, domainKey, _, _, _ := s.snapshotState()
	return domainKey
}

// TenantHash returns a consistent snapshot of the slot's SHA-256 tenant digest.
func (s *CacheLineSlot) TenantHash() [32]byte {
	if s == nil {
		return [32]byte{}
	}
	_, _, tenantHash, _, _ := s.snapshotState()
	return tenantHash
}

// ActiveTokenCount atomically returns the slot's active token count.
func (s *CacheLineSlot) ActiveTokenCount() uint32 {
	if s == nil {
		return 0
	}
	return atomic.LoadUint32(&s.activeTokens)
}

// TotalAccessCount atomically returns the slot's acquisition count.
func (s *CacheLineSlot) TotalAccessCount() uint64 {
	if s == nil {
		return 0
	}
	return atomic.LoadUint64(&s.accessCount)
}

// snapshotState reads all mutable slot fields atomically while verifying via the
// seqlock sequence counter and flags that no concurrent AcquireSlot, ReleaseSlot,
// or state transition occurred mid-read.
func (s *CacheLineSlot) snapshotState() (flags uint32, domainKey uint8, tenantHash [32]byte, activeTokens uint32, accessCount uint64) {
	for {
		seqBefore := atomic.LoadUint32(&s.seq)
		if seqBefore&1 != 0 {
			runtime.Gosched()
			continue
		}

		flagsBefore := atomic.LoadUint32(&s.flags)
		if flagsBefore == SlotStatusInitializing || flagsBefore == SlotStatusReleasing {
			runtime.Gosched()
			continue
		}

		w0 := atomic.LoadUint64(&s.tenantWords[0])
		w1 := atomic.LoadUint64(&s.tenantWords[1])
		w2 := atomic.LoadUint64(&s.tenantWords[2])
		w3 := atomic.LoadUint64(&s.tenantWords[3])
		dk := uint8(atomic.LoadUint32(&s.domainKey))
		tokens := atomic.LoadUint32(&s.activeTokens)
		accesses := atomic.LoadUint64(&s.accessCount)

		flagsAfter := atomic.LoadUint32(&s.flags)
		seqAfter := atomic.LoadUint32(&s.seq)
		if seqBefore == seqAfter && flagsBefore == flagsAfter {
			binary.BigEndian.PutUint64(tenantHash[0:8], w0)
			binary.BigEndian.PutUint64(tenantHash[8:16], w1)
			binary.BigEndian.PutUint64(tenantHash[16:24], w2)
			binary.BigEndian.PutUint64(tenantHash[24:32], w3)
			return flagsBefore, dk, tenantHash, tokens, accesses
		}

		runtime.Gosched()
	}
}

// BufferPool manages a pre-allocated, 64-byte aligned circular ring of cache-line slots.
type BufferPool struct {
	backing   []byte
	slots     []CacheLineSlot
	slotCount int
	head      uint64
	active    int64
}

// NewBufferPool allocates and initializes a buffer pool whose base address and
// individual slots are all aligned to 64-byte cache-line boundaries.
func NewBufferPool(capacity int) *BufferPool {
	if capacity <= 0 {
		capacity = DefaultSlotCount
	}

	// Over-allocate by (CacheLineSize - 1) bytes so the slot slice can start on
	// an exact 64-byte boundary. CacheLineSlot contains no Go pointers, so aliasing
	// a retained backing byte slice via unsafe.Slice is safe for the Go GC and checkptr.
	backing := make([]byte, capacity*CacheLineSize+CacheLineSize-1)
	rawAddr := uintptr(unsafe.Pointer(&backing[0]))
	alignOffset := uintptr(0)
	if rem := rawAddr % uintptr(CacheLineSize); rem != 0 {
		alignOffset = uintptr(CacheLineSize) - rem
	}

	alignedPtr := unsafe.Add(unsafe.Pointer(&backing[0]), alignOffset)
	slots := unsafe.Slice((*CacheLineSlot)(alignedPtr), capacity)

	pool := &BufferPool{
		backing:   backing,
		slots:     slots,
		slotCount: capacity,
	}

	for i := 0; i < capacity; i++ {
		pool.slots[i].SlotID = uint32(i)
		atomic.StoreUint32(&pool.slots[i].flags, SlotStatusFree)
	}

	return pool
}

// AcquireSlot reserves a free slot via CAS into SlotStatusInitializing, initializes all
// metadata before publishing SlotStatusActive, and returns a pointer to the reserved slot.
func (p *BufferPool) AcquireSlot(tenantID string, domainKey uint8) (*CacheLineSlot, error) {
	if p == nil || p.slotCount == 0 {
		return nil, ErrPoolExhausted
	}

	startIndex := atomic.AddUint64(&p.head, 1) % uint64(p.slotCount)

	// Hash tenant ID into fixed 32 bytes.
	h := sha256.Sum256([]byte(tenantID))
	w0 := binary.BigEndian.Uint64(h[0:8])
	w1 := binary.BigEndian.Uint64(h[8:16])
	w2 := binary.BigEndian.Uint64(h[16:24])
	w3 := binary.BigEndian.Uint64(h[24:32])

	for i := 0; i < p.slotCount; i++ {
		idx := (startIndex + uint64(i)) % uint64(p.slotCount)
		slot := &p.slots[idx]

		if atomic.CompareAndSwapUint32(&slot.flags, SlotStatusFree, SlotStatusInitializing) {
			atomic.AddUint32(&slot.seq, 1) // odd: write in progress

			atomic.StoreUint64(&slot.tenantWords[0], w0)
			atomic.StoreUint64(&slot.tenantWords[1], w1)
			atomic.StoreUint64(&slot.tenantWords[2], w2)
			atomic.StoreUint64(&slot.tenantWords[3], w3)
			atomic.StoreUint32(&slot.domainKey, uint32(domainKey))
			atomic.StoreUint32(&slot.activeTokens, 0)
			atomic.AddUint64(&slot.accessCount, 1)
			atomic.AddInt64(&p.active, 1)

			atomic.AddUint32(&slot.seq, 1) // even: payload committed
			atomic.StoreUint32(&slot.flags, SlotStatusActive)
			return slot, nil
		}
	}

	return nil, ErrPoolExhausted
}

// ReleaseSlot transitions an Active slot to Releasing, clears all tenant metadata before
// publishing SlotStatusFree, and decrements the active slot counter.
func (p *BufferPool) ReleaseSlot(slotID uint32) bool {
	if p == nil || int(slotID) >= p.slotCount {
		return false
	}

	slot := &p.slots[slotID]
	if !atomic.CompareAndSwapUint32(&slot.flags, SlotStatusActive, SlotStatusReleasing) {
		return false
	}

	atomic.AddUint32(&slot.seq, 1) // odd: cleanup in progress

	atomic.StoreUint64(&slot.tenantWords[0], 0)
	atomic.StoreUint64(&slot.tenantWords[1], 0)
	atomic.StoreUint64(&slot.tenantWords[2], 0)
	atomic.StoreUint64(&slot.tenantWords[3], 0)
	atomic.StoreUint32(&slot.domainKey, 0)
	atomic.StoreUint32(&slot.activeTokens, 0)
	atomic.AddInt64(&p.active, -1)

	atomic.AddUint32(&slot.seq, 1) // even: cleanup committed
	atomic.StoreUint32(&slot.flags, SlotStatusFree)
	return true
}

// FenceSlot transitions an Active slot into a logical software fence (SlotStatusFenced).
// It refuses to fence Free, Initializing, Releasing, or already Fenced slots, and does not
// configure hardware Intel MPK/PKRU protection keys in this userspace implementation.
func (p *BufferPool) FenceSlot(slotID uint32) bool {
	if p == nil || int(slotID) >= p.slotCount {
		return false
	}

	slot := &p.slots[slotID]
	return atomic.CompareAndSwapUint32(&slot.flags, SlotStatusActive, SlotStatusFenced)
}

// UnfenceSlot transitions a Fenced slot back to Active so it can resume normal use or release.
func (p *BufferPool) UnfenceSlot(slotID uint32) bool {
	if p == nil || int(slotID) >= p.slotCount {
		return false
	}

	slot := &p.slots[slotID]
	return atomic.CompareAndSwapUint32(&slot.flags, SlotStatusFenced, SlotStatusActive)
}

// GetSlot retrieves a direct pointer to a slot by ID.
func (p *BufferPool) GetSlot(slotID uint32) *CacheLineSlot {
	if p == nil || int(slotID) >= p.slotCount {
		return nil
	}
	return &p.slots[slotID]
}

// VerifyCacheAlignment verifies both 64-byte base-address alignment and exact 64-byte
// slot stride across every slot in the pool.
func (p *BufferPool) VerifyCacheAlignment() bool {
	if p == nil || len(p.slots) == 0 {
		return false
	}
	if unsafe.Sizeof(CacheLineSlot{}) != CacheLineSize {
		return false
	}

	baseAddr := uintptr(unsafe.Pointer(&p.slots[0]))
	if baseAddr%uintptr(CacheLineSize) != 0 {
		return false
	}

	for i := 0; i < len(p.slots); i++ {
		slotAddr := uintptr(unsafe.Pointer(&p.slots[i]))
		if slotAddr%uintptr(CacheLineSize) != 0 {
			return false
		}
		if slotAddr-baseAddr != uintptr(i*CacheLineSize) {
			return false
		}
	}
	return true
}

// SlotSnapshot holds telemetry data for a single slot in the matrix view.
type SlotSnapshot struct {
	Index            int    `json:"index"`
	SlotID           uint32 `json:"slot_id"`
	TenantID         string `json:"tenant_id"`
	DomainKey        uint8  `json:"domain_key"`
	Status           string `json:"status"`
	ActiveTokens     uint32 `json:"active_tokens"`
	AccessCount      uint64 `json:"access_count"`
	MemoryOffsetHex  string `json:"memory_offset_hex"`
	CacheLineAligned bool   `json:"cache_line_aligned"`
	ColorChannel     string `json:"color_channel"`
}

// BufferPoolSnapshot contains the telemetry state of the buffer pool matrix.
type BufferPoolSnapshot struct {
	TotalSlots              int            `json:"total_slots"`
	ActiveSlots             int            `json:"active_slots"`
	FencedSlots             int            `json:"fenced_slots"`
	CacheLineSize           int            `json:"cache_line_size"`
	CacheAlignmentCongruent bool           `json:"cache_alignment_congruent"`
	L1L2ThrashingRate       *float64       `json:"l1_l2_thrashing_rate"`
	L1L2ThrashingAvailable  bool           `json:"l1_l2_thrashing_available"`
	L1L2ThrashingNote       string         `json:"l1_l2_thrashing_note,omitempty"`
	GCPausesMs              float64        `json:"gc_pauses_ms"`
	GCPauseCount            int64          `json:"gc_pause_count"`
	GCPausesAvailable       bool           `json:"gc_pauses_available"`
	IsolationMode           string         `json:"isolation_mode"`
	HardwareMPKEnforced     bool           `json:"hardware_mpk_enforced"`
	ActiveTenantsCount      int            `json:"active_tenants_count"`
	Slots                   []SlotSnapshot `json:"slots"`
}

// GetTenantNeonColor returns the display color channel associated with a logical tenant domain.
func GetTenantNeonColor(domainKey uint8, flags uint32) string {
	if flags == SlotStatusFree {
		return "#1E293B" // Slate / guard band
	}
	if flags == SlotStatusFenced {
		return "#FF0055" // Neon Crimson Fence
	}
	switch domainKey % 4 {
	case 1:
		return "#00F0FF" // Neon Cyan: Tenant Alpha (Finance)
	case 2:
		return "#00FF88" // Neon Emerald: Tenant Beta (Defense)
	case 3:
		return "#BD00FF" // Neon Violet: Tenant Gamma (Health)
	default:
		return "#FFB800" // Neon Amber: Tenant Delta (Sovereign AI)
	}
}

// measureGCPausesMs reads cumulative GC stop-the-world pause duration from runtime/debug
// without triggering a stop-the-world heap scan.
func measureGCPausesMs() (float64, int64, bool) {
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	pauseMs := float64(stats.PauseTotal.Nanoseconds()) / 1e6
	return pauseMs, stats.NumGC, true
}

// Snapshot exports a race-safe point-in-time matrix view of the buffer pool.
func (p *BufferPool) Snapshot() BufferPoolSnapshot {
	if p == nil {
		gcMs, gcCount, gcAvail := measureGCPausesMs()
		return BufferPoolSnapshot{
			CacheLineSize:           CacheLineSize,
			CacheAlignmentCongruent: false,
			L1L2ThrashingRate:       nil,
			L1L2ThrashingAvailable:  false,
			L1L2ThrashingNote:       "Unavailable in portable userspace runtime; requires hardware PMU counters (e.g., perf_event_open).",
			GCPausesMs:              gcMs,
			GCPauseCount:            gcCount,
			GCPausesAvailable:       gcAvail,
			IsolationMode:           IsolationModeLogicalTag,
			HardwareMPKEnforced:     false,
		}
	}

	slotsSnapshot := make([]SlotSnapshot, p.slotCount)
	activeCount := 0
	fencedCount := 0
	tenantSet := make(map[string]struct{})

	baseAddr := uintptr(0)
	if len(p.slots) > 0 {
		baseAddr = uintptr(unsafe.Pointer(&p.slots[0]))
	}

	for i := 0; i < p.slotCount; i++ {
		slot := &p.slots[i]
		flags, domainKey, tenantHash, activeTokens, accessCount := slot.snapshotState()

		var statusStr string
		switch flags {
		case SlotStatusActive:
			statusStr = "ACTIVE"
			activeCount++
		case SlotStatusFenced:
			statusStr = "FENCED"
			fencedCount++
		default:
			statusStr = "FREE"
		}

		tenantHex := hex.EncodeToString(tenantHash[:])
		if flags != SlotStatusFree {
			tenantSet[tenantHex[:8]] = struct{}{}
		}

		slotAddr := uintptr(unsafe.Pointer(slot))
		offset := slotAddr - baseAddr

		slotsSnapshot[i] = SlotSnapshot{
			Index:            i,
			SlotID:           slot.SlotID,
			TenantID:         tenantHex[:12],
			DomainKey:        domainKey,
			Status:           statusStr,
			ActiveTokens:     activeTokens,
			AccessCount:      accessCount,
			MemoryOffsetHex:  fmt.Sprintf("0x%04X", offset),
			CacheLineAligned: (slotAddr % uintptr(CacheLineSize)) == 0,
			ColorChannel:     GetTenantNeonColor(domainKey, flags),
		}
	}

	gcPausesMs, gcPauseCount, gcAvail := measureGCPausesMs()

	return BufferPoolSnapshot{
		TotalSlots:              p.slotCount,
		ActiveSlots:             activeCount,
		FencedSlots:             fencedCount,
		CacheLineSize:           CacheLineSize,
		CacheAlignmentCongruent: p.VerifyCacheAlignment(),
		L1L2ThrashingRate:       nil,
		L1L2ThrashingAvailable:  false,
		L1L2ThrashingNote:       "Unavailable in portable userspace runtime; requires hardware PMU counters (e.g., perf_event_open).",
		GCPausesMs:              gcPausesMs,
		GCPauseCount:            gcPauseCount,
		GCPausesAvailable:       gcAvail,
		IsolationMode:           IsolationModeLogicalTag,
		HardwareMPKEnforced:     false,
		ActiveTenantsCount:      len(tenantSet),
		Slots:                   slotsSnapshot,
	}
}
