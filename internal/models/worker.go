package models

type Worker struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Region string  `json:"region,omitempty"`
	CPU    float64 `json:"cpu"`
	Memory int64   `json:"memory"`
}
