package models

type Telemetry struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	NetworkRX   uint64  `json:"network_rx"`
	NetworkTX   uint64  `json:"network_tx"`
}
