package files

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDirEmpty(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	assert.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	isEmpty, err := IsDirEmpty(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, true, isEmpty)

	f, err := os.Create(path.Join(tempDir, "test.txt"))
	assert.NoError(t, err)

	f.Close()

	isEmpty, err = IsDirEmpty(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, false, isEmpty)
}
