package maps_test

import (
	"testing"

	"github.com/katallaxie/pkg/utils/maps"

	"github.com/stretchr/testify/assert"
)

func TestFromSlice(t *testing.T) {
	m := maps.FromSlice([]string{"a=b", "c=d"})
	assert.Equal(t, map[string]string{"a": "b", "c": "d"}, m)
}
