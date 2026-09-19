package smart

import (
	"sync"
	"time"
)

type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Requests  float64   `json:"requests"`
	Latency   float64   `json:"latency"`
	CPU       float64   `json:"cpu"`
	Memory    float64   `json:"memory"`
}

type Signal struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Baseline  float64 `json:"baseline"`
	Deviation float64 `json:"deviation"`
	Threshold float64 `json:"threshold"`
	Triggered bool    `json:"triggered"`
}

type Result struct {
	Timestamp time.Time `json:"timestamp"`
	Anomaly   bool      `json:"anomaly"`
	Signals   []Signal  `json:"signals"`
}

type Detector struct {
	mu        sync.RWMutex
	window    []Snapshot
	maxWindow int
}

func NewDetector(maxWindow int) *Detector {
	if maxWindow < 2 {
		maxWindow = 30
	}

	return &Detector{
		window:    make([]Snapshot, 0, maxWindow),
		maxWindow: maxWindow,
	}
}

func (d *Detector) Observe(snapshot Snapshot) Result {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.window = append(d.window, snapshot)

	if len(d.window) > d.maxWindow {
		d.window = d.window[len(d.window)-d.maxWindow:]
	}

	signals := []Signal{
		d.evaluate("latency", snapshot.Latency, func(s Snapshot) float64 {
			return s.Latency
		}),
		d.evaluate("cpu", snapshot.CPU, func(s Snapshot) float64 {
			return s.CPU
		}),
		d.evaluate("memory", snapshot.Memory, func(s Snapshot) float64 {
			return s.Memory
		}),
	}

	anomaly := false
	for _, signal := range signals {
		if signal.Triggered {
			anomaly = true
			break
		}
	}

	return Result{
		Timestamp: snapshot.Timestamp,
		Anomaly:   anomaly,
		Signals:   signals,
	}
}

func (d *Detector) evaluate(
	name string,
	value float64,
	extract func(Snapshot) float64,
) Signal {
	if len(d.window) <= 1 {
		return Signal{
			Name:      name,
			Value:     value,
			Baseline:  value,
			Deviation: 0,
			Threshold: 0.50,
			Triggered: false,
		}
	}

	var sum float64
	for _, sample := range d.window[:len(d.window)-1] {
		sum += extract(sample)
	}

	baseline := sum / float64(len(d.window)-1)

	deviation := 0.0
	if baseline != 0 {
		deviation = (value - baseline) / baseline
	}

	const threshold = 0.50

	return Signal{
		Name:      name,
		Value:     value,
		Baseline:  baseline,
		Deviation: deviation,
		Threshold: threshold,
		Triggered: deviation >= threshold || deviation <= -threshold,
	}
}
