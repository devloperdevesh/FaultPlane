package speculative

import (
	"errors"
	"sync"
)

var (
	ErrInvalidCandidate = errors.New("speculative candidate is invalid")
	ErrCandidateExists  = errors.New("speculative candidate already exists")
	ErrCandidateMissing = errors.New("speculative candidate not found")
	ErrAlreadyCommitted = errors.New("speculative candidate already committed")
	ErrAlreadyCancelled = errors.New("speculative candidate already cancelled")
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusCommitted Status = "committed"
	StatusCancelled Status = "cancelled"
)

type Candidate struct {
	ID      string
	Step    uint64
	Payload []byte
	Status  Status
}

type Coordinator struct {
	mu         sync.RWMutex
	candidates map[string]Candidate
	committed  string
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		candidates: make(map[string]Candidate),
	}
}

func (c *Coordinator) Register(candidate Candidate) error {
	if c == nil || candidate.ID == "" {
		return ErrInvalidCandidate
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.candidates[candidate.ID]; exists {
		return ErrCandidateExists
	}

	payload := append([]byte(nil), candidate.Payload...)

	c.candidates[candidate.ID] = Candidate{
		ID:      candidate.ID,
		Step:    candidate.Step,
		Payload: payload,
		Status:  StatusPending,
	}

	return nil
}

func (c *Coordinator) Commit(id string) (Candidate, error) {
	if c == nil || id == "" {
		return Candidate{}, ErrCandidateMissing
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	candidate, ok := c.candidates[id]
	if !ok {
		return Candidate{}, ErrCandidateMissing
	}

	switch candidate.Status {
	case StatusCommitted:
		return Candidate{}, ErrAlreadyCommitted
	case StatusCancelled:
		return Candidate{}, ErrAlreadyCancelled
	}

	if c.committed != "" && c.committed != id {
		return Candidate{}, errors.New("another speculative candidate is already committed")
	}

	candidate.Status = StatusCommitted
	candidate.Payload = append([]byte(nil), candidate.Payload...)

	c.candidates[id] = candidate
	c.committed = id

	return candidate, nil
}

func (c *Coordinator) Cancel(id string) error {
	if c == nil || id == "" {
		return ErrCandidateMissing
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	candidate, ok := c.candidates[id]
	if !ok {
		return ErrCandidateMissing
	}

	switch candidate.Status {
	case StatusCommitted:
		return ErrAlreadyCommitted
	case StatusCancelled:
		return ErrAlreadyCancelled
	}

	candidate.Status = StatusCancelled
	candidate.Payload = nil
	c.candidates[id] = candidate

	return nil
}

func (c *Coordinator) Get(id string) (Candidate, bool) {
	if c == nil {
		return Candidate{}, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	candidate, ok := c.candidates[id]
	if !ok {
		return Candidate{}, false
	}

	candidate.Payload = append([]byte(nil), candidate.Payload...)

	return candidate, true
}

func (c *Coordinator) Committed() (Candidate, bool) {
	if c == nil {
		return Candidate{}, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.committed == "" {
		return Candidate{}, false
	}

	candidate, ok := c.candidates[c.committed]
	if !ok {
		return Candidate{}, false
	}

	candidate.Payload = append([]byte(nil), candidate.Payload...)

	return candidate, true
}
