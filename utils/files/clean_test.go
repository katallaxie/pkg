package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClean(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	assert.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	path := filepath.Join(tempDir, "example")
	err = MkdirAll(path, 0o755)
	assert.NoError(t, err)

	file, err := os.Create(filepath.Join(tempDir, "example", "test.txt"))
	assert.NoError(t, err)
	defer file.Close()

	err = Clean(path, 0o755)
	assert.NoError(t, err)

	_, err = os.Stat(path)
	assert.NoError(t, err)

	_, err = os.Stat(filepath.Join(tempDir, "example", "test.txt"))
	assert.Error(t, err)
}
