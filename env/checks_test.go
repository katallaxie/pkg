package env_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/katallaxie/pkg/env"

	"github.com/stretchr/testify/require"
)

func TestHasUser(t *testing.T) {
	t.Parallel()

	require.NoError(t, env.HasUser()(context.Background()))
}

func TestIsDirEmpty(t *testing.T) {
	t.Parallel()

	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	require.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	require.NoError(t, env.IsDirEmpty(tempDir)(context.Background()))

	f, err := os.Create(path.Join(tempDir, "test.txt"))
	require.NoError(t, err)

	f.Close()

	require.Error(t, env.IsDirEmpty(tempDir)(context.Background()))
}
