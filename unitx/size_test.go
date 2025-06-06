package unitx_test

import (
	"testing"

	"github.com/katallaxie/pkg/unitx"
	"github.com/stretchr/testify/require"
)

func TestHumanSize_ToInt64(t *testing.T) {
	tests := []struct {
		name    string
		h       unitx.HumanSize
		want    int64
		wantErr bool
	}{
		{
			name:    "zero",
			h:       unitx.HumanSize("0"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "binary",
			h:       unitx.HumanSize("1kB"),
			want:    1024,
			wantErr: false,
		},
		{
			name:    "binary",
			h:       unitx.HumanSize("1KiB"),
			want:    1024,
			wantErr: false,
		},
		{
			name:    "decimal",
			h:       unitx.HumanSize("1M"),
			want:    1000 * 1000,
			wantErr: false,
		},
		{
			name:    "decimal",
			h:       unitx.HumanSize("1MB"),
			want:    1024 * 1024,
			wantErr: false,
		},
		{
			name:    "invalid",
			h:       unitx.HumanSize("1MBa"),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.ToInt64()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestHumanSize_ToInt(t *testing.T) {
	tests := []struct {
		name    string
		h       unitx.HumanSize
		want    int
		wantErr bool
	}{
		{
			name:    "zero",
			h:       unitx.HumanSize("0"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "binary",
			h:       unitx.HumanSize("1kB"),
			want:    1024,
			wantErr: false,
		},
		{
			name:    "binary",
			h:       unitx.HumanSize("1KiB"),
			want:    1024,
			wantErr: false,
		},
		{
			name:    "decimal",
			h:       unitx.HumanSize("1M"),
			want:    1000 * 1000,
			wantErr: false,
		},
		{
			name:    "decimal",
			h:       unitx.HumanSize("1MB"),
			want:    1024 * 1024,
			wantErr: false,
		},
		{
			name:    "invalid",
			h:       unitx.HumanSize("1MBa"),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.ToInt()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestHumanSize_ToInt32(t *testing.T) {
	tests := []struct {
		name    string
		h       unitx.HumanSize
		want    int32
		wantErr bool
	}{
		{
			name:    "zero",
			h:       unitx.HumanSize("0"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "binary",
			h:       unitx.HumanSize("1kB"),
			want:    1024,
			wantErr: false,
		},
		{
			name:    "binary",
			h:       unitx.HumanSize("1KiB"),
			want:    1024,
			wantErr: false,
		},
		{
			name:    "decimal",
			h:       unitx.HumanSize("1M"),
			want:    1000 * 1000,
			wantErr: false,
		},
		{
			name:    "decimal",
			h:       unitx.HumanSize("1MB"),
			want:    1024 * 1024,
			wantErr: false,
		},
		{
			name:    "invalid",
			h:       unitx.HumanSize("1MBa"),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.ToInt32()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
