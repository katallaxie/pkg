package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	assert.Equal(t, Zero[int](), 0)
	assert.Equal(t, Zero[string](), "")
	assert.Equal(t, Zero[bool](), false)
	assert.Equal(t, Zero[interface{}](), nil)
	assert.Equal(t, Zero[struct{}](), struct{}{})
}

func TestIsZero(t *testing.T) {
	assert.True(t, IsZero(0))
	assert.True(t, IsZero(""))
	assert.True(t, IsZero(false))
	assert.True(t, IsZero(struct{}{}))
}
