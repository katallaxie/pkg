package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	assert.Equal(t, 0, Zero[int]())
	assert.Equal(t, "", Zero[string]())
	assert.Equal(t, false, Zero[bool]())
	assert.Equal(t, nil, Zero[interface{}]())
	assert.Equal(t, struct{}{}, Zero[struct{}]())
}

func TestIsZero(t *testing.T) {
	assert.True(t, IsZero(0))
	assert.True(t, IsZero(""))
	assert.True(t, IsZero(false))
	assert.True(t, IsZero(struct{}{}))
}
