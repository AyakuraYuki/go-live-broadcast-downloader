package sequence

import (
	"testing"
)

func TestID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for index := 0; index < 1; index++ {
				t.Log(ID())
			}
		})
	}
}
