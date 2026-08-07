package simulator

import "time"

type Workflow struct {
	ID string

	CurrentStep int

	Status string

	UpdatedAt time.Time
}
