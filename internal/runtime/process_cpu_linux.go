//go:build linux

package runtime

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ProcessCPUSampler struct {
	mu sync.Mutex

	lastWall time.Time
	lastCPU  time.Duration
}

func NewProcessCPUSampler() *ProcessCPUSampler {
	cpu, _ := processCPUTime()

	return &ProcessCPUSampler{
		lastWall: time.Now(),
		lastCPU:  cpu,
	}
}

func (s *ProcessCPUSampler) Sample() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	cpu, ok := processCPUTime()
	if !ok {
		return 0
	}

	wallDelta := now.Sub(s.lastWall)
	cpuDelta := cpu - s.lastCPU

	s.lastWall = now
	s.lastCPU = cpu

	if wallDelta <= 0 || cpuDelta <= 0 {
		return 0
	}

	usage := float64(cpuDelta) / float64(wallDelta) * 100.0

	if usage < 0 {
		return 0
	}

	if usage > 100 {
		return 100
	}

	return usage
}

func processCPUTime() (time.Duration, bool) {
	data, err := os.ReadFile("/proc/self/stat")
	if err != nil {
		return 0, false
	}

	fields := strings.Fields(string(data))
	if len(fields) < 17 {
		return 0, false
	}

	utime, err := strconv.ParseUint(fields[13], 10, 64)
	if err != nil {
		return 0, false
	}

	stime, err := strconv.ParseUint(fields[14], 10, 64)
	if err != nil {
		return 0, false
	}

	total := utime + stime

	const clockTicksPerSecond = 100

	return time.Duration(
		float64(total) / clockTicksPerSecond * float64(time.Second),
	), true
}
