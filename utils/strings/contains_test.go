package strings_test

import (
	"testing"

	"github.com/katallaxie/pkg/utils/strings"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	var tests = []struct {
		in       []string
		test     string
		expected bool
	}{
		{
			in:       []string{"foo", "bar"},
			test:     "foo",
			expected: true,
		},
		{
			in:       []string{"bar"},
			test:     "foo",
			expected: false,
		},
	}

	for _, tc := range tests {
		out := strings.Contains(tc.in, tc.test)
		assert.Equal(t, tc.expected, out)
	}
}
