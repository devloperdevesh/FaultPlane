package telemetry

import "testing"

func BenchmarkRecordLatency(b *testing.B) {
	registry := NewRegistry()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		registry.RecordLatency(1)
	}
}

func BenchmarkAverageLatency(b *testing.B) {
	registry := NewRegistry()

	for i := 0; i < 1024; i++ {
		registry.RecordLatency(1)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = registry.Snapshot().AverageLatency()
	}
}

func BenchmarkSnapshot(b *testing.B) {
	registry := NewRegistry()

	registry.RecordLatency(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = registry.Snapshot()
	}
}
