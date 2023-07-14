package env_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/katallaxie/pkg/env"
	"github.com/stretchr/testify/assert"
)

func TestHasUser(t *testing.T) {
	t.Parallel()

	assert.NoError(t, env.HasUser()(context.Background()))
}

func TestIsDirEmpty(t *testing.T) {
	t.Parallel()

	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	assert.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	assert.NoError(t, env.IsDirEmpty(tempDir)(context.Background()))

	f, err := os.Create(path.Join(tempDir, "test.txt"))
	assert.NoError(t, err)

	f.Close()

	assert.Error(t, env.IsDirEmpty(tempDir)(context.Background()))
}
