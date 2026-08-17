package failover

import "testing"

func TestConnectionStateValidation(t *testing.T) {

	tests := []struct {
		name  string
		state ConnectionState
		want  bool
	}{
		{
			name:  "healthy",
			state: StateHealthy,
			want:  true,
		},
		{
			name:  "failed",
			state: StateFailed,
			want:  true,
		},
		{
			name:  "invalid",
			state: ConnectionState("unknown"),
			want:  false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			if got := tt.state.Valid(); got != tt.want {

				t.Fatalf(
					"expected %v got %v",
					tt.want,
					got,
				)
			}
		})
	}
}
