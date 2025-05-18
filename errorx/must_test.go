package errorx_test

import (
	"fmt"
	"testing"

	"github.com/katallaxie/pkg/errorx"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: true,
		},
		{
			name:     "non-nil error",
			err:      fmt.Errorf("error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := errorx.Empty(tt.err)
			require.Equal(t, tt.expected, actual)
		})
	}
}
