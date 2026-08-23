package models

type Metrics struct {
	Requests int64   `json:"requests"`
	Latency  float64 `json:"latency"`
	CPU      float64 `json:"cpu"`
	Memory   int64   `json:"memory"`
}
