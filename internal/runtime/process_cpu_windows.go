package runtime

import (
	"sync"
	"time"

	"golang.org/x/sys/windows"
)

type ProcessCPUSampler struct {
	mu sync.Mutex

	process  windows.Handle
	lastWall time.Time
	lastCPU  time.Duration
}

func NewProcessCPUSampler() *ProcessCPUSampler {
	process := windows.CurrentProcess()

	cpu, _ := processCPUTime(process)

	return &ProcessCPUSampler{
		process:  process,
		lastWall: time.Now(),
		lastCPU:  cpu,
	}
}

func (s *ProcessCPUSampler) Sample() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	cpu, ok := processCPUTime(s.process)
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

func processCPUTime(
	process windows.Handle,
) (time.Duration, bool) {
	var creation windows.Filetime
	var exit windows.Filetime
	var kernel windows.Filetime
	var user windows.Filetime

	if err := windows.GetProcessTimes(
		process,
		&creation,
		&exit,
		&kernel,
		&user,
	); err != nil {
		return 0, false
	}

	kernel100ns := uint64(kernel.HighDateTime)<<32 |
		uint64(kernel.LowDateTime)

	user100ns := uint64(user.HighDateTime)<<32 |
		uint64(user.LowDateTime)

	total100ns := kernel100ns + user100ns

	return time.Duration(total100ns) * 100 * time.Nanosecond, true
}
