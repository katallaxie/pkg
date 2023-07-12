package env_test

import (
	"context"
	"testing"

	"github.com/katallaxie/pkg/env"
	"github.com/stretchr/testify/assert"
)

func TestHasUser(t *testing.T) {
	t.Parallel()

	assert.NoError(t, env.HasUser()(context.Background()))
}
