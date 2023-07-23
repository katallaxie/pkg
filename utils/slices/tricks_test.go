package slices_test

import (
	"testing"

	"github.com/katallaxie/pkg/utils/slices"

	"github.com/stretchr/testify/assert"
)

func TestPop(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name     string
		input    []int
		el       int
		expected []int
	}{
		{
			name:     "pop from slice",
			input:    []int{1, 2, 3},
			el:       3,
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			el, actual := slices.Pop(tt.input...)
			assert.Equal(t, tt.el, el)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPush(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name     string
		el       int
		input    []int
		expected []int
	}{
		{
			name:     "push to slice",
			input:    []int{1, 2},
			el:       3,
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Push(tt.el, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCut(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name     string
		start    int
		end      int
		input    []int
		expected []int
	}{
		{
			name:     "cut from slice",
			start:    1,
			end:      2,
			input:    []int{1, 2, 3},
			expected: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Cut(tt.start, tt.end, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name     string
		idx      int
		input    []int
		expected []int
	}{
		{
			name:     "delete from slice",
			idx:      1,
			input:    []int{1, 2, 3},
			expected: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Delete(tt.idx, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
