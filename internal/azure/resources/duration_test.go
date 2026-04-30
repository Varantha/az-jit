package azure

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "7 Hours",
			input:   "PT7H",
			want:    7 * time.Hour,
			wantErr: false,
		},
		{
			name:    "30 Minutes",
			input:   "PT30M",
			want:    30 * time.Minute,
			wantErr: false,
		},
		{
			name:    "1 Hour 30 Minutes",
			input:   "PT1H30M",
			want:    1*time.Hour + 30*time.Minute,
			wantErr: false,
		},
		{
			name:    "90 Days",
			input:   "P90D",
			want:    90 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "1 Day 5 hours",
			input:   "P1DT5H",
			want:    1*24*time.Hour + 5*time.Hour,
			wantErr: false,
		},
		{
			name:    "Incorrect String",
			input:   "garbage",
			want:    0,
			wantErr: true,
		},
		{
			name:    "1 Year",
			input:   "P1Y",
			want:    0,
			wantErr: true,
		},
		{
			name:    "1 Month",
			input:   "P1M",
			want:    0,
			wantErr: true,
		},
		{
			name:    "1 Second without T",
			input:   "P1S",
			want:    0,
			wantErr: true,
		},
		{
			name:    "1 Hour without T",
			input:   "P1H",
			want:    0,
			wantErr: true,
		},
		{
			name:    "day and minute without T",
			input:   "P1D1M",
			want:    0,
			wantErr: true,
		},
		{
			name:    "duplicate H",
			input:   "PT1H1H",
			want:    0,
			wantErr: true,
		},
		{
			name:    "0 values",
			input:   "PT0S",
			want:    0,
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseISO8601Duration(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
