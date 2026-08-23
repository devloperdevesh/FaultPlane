package models

type Worker struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	CPU    float64 `json:"cpu"`
	Memory int64   `json:"memory"`
}
