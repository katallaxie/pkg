package envx_test

import (
	"testing"

	"github.com/katallaxie/pkg/envx"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	c := envx.NewChecker()
	assert.NotNil(t, c)
}
