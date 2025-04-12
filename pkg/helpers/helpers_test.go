package helpers

import (
	"testing"
)

func TestFindMin(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		want    int
		wantErr bool
	}{
		{
			name:    "empty slice",
			input:   []int{},
			want:    0,
			wantErr: true,
		},
		{
			name:    "single element",
			input:   []int{5},
			want:    5,
			wantErr: false,
		},
		{
			name:    "multiple elements",
			input:   []int{3, 1, 4, 1, 5, 9},
			want:    1,
			wantErr: false,
		},
		{
			name:    "negative numbers",
			input:   []int{-3, -1, -4, -1, -5, -9},
			want:    -9,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindMin(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindMin() = %v, want %v", got, tt.want)
			}
		})
	}
}
