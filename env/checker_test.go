package env_test

import (
	"testing"

	"github.com/katallaxie/pkg/env"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	c := env.NewChecker()
	assert.NotNil(t, c)
}
