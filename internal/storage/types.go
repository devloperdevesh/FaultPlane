package storage

import "time"

// Record represents persisted workflow state.
//
// Storage layer owns persistence models.
// It must not import control package.
type Record struct {
	ID string `json:"id"`

	Step int `json:"step"`

	Status string `json:"status"`

	Payload string `json:"payload"`

	UpdatedAt time.Time `json:"updated_at"`
}
