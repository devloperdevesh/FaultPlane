package pool

import (
	"errors"
	"sync"
)

var (
	ErrBufferSaturated = errors.New("faultplane: ring buffer saturated")
	ErrBufferExhausted = errors.New("faultplane: ring buffer empty")
)

type ZeroAllocRingBuffer struct {
	mu           sync.Mutex
	storage      []uint64
	writePointer uint64
	readPointer  uint64
}

func NewRingBuffer(slotPower uint32) *ZeroAllocRingBuffer {
	if slotPower == 0 {
		slotPower = 1
	}

	size := uint64(1) << slotPower

	return &ZeroAllocRingBuffer{
		storage: make([]uint64, size),
	}
}

func (b *ZeroAllocRingBuffer) Capacity() uint64 {
	return uint64(len(b.storage))
}

func (b *ZeroAllocRingBuffer) Len() uint64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.writePointer - b.readPointer
}

func (b *ZeroAllocRingBuffer) EnqueueTokenEnvelope(
	tokenIndex uint64,
) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.writePointer-b.readPointer >= uint64(len(b.storage)) {
		return ErrBufferSaturated
	}

	index := b.writePointer % uint64(len(b.storage))
	b.storage[index] = tokenIndex
	b.writePointer++

	return nil
}

func (b *ZeroAllocRingBuffer) DequeueTokenEnvelope() (
	uint64,
	error,
) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.writePointer == b.readPointer {
		return 0, ErrBufferExhausted
	}

	index := b.readPointer % uint64(len(b.storage))
	value := b.storage[index]

	b.storage[index] = 0
	b.readPointer++

	return value, nil
}
