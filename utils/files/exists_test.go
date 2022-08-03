package files

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	assert.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	ok, err := FileExists(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)

	path := strings.Join([]string{tempDir, "example.txt"}, "/")
	f, err := os.Create(path)
	assert.NoError(t, err)
	f.Close()

	ok, err = FileExists(path)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)

	err = os.Remove(path)
	assert.NoError(t, err)

	ok, err = FileExists(path)
	assert.Error(t, err)
	assert.Equal(t, false, ok)
}

func TestFileNotExists(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	assert.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	path := strings.Join([]string{tempDir, "example.txt"}, "/")
	f, err := os.Create(path)
	assert.NoError(t, err)
	f.Close()

	ok, err := FileNotExists(path)
	assert.NoError(t, err)
	assert.Equal(t, false, ok)

	err = os.Remove(path)
	assert.NoError(t, err)

	ok, err = FileNotExists(path)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)

	path = strings.Join([]string{tempDir, "demo123"}, "/")
	ok, err = FileNotExists(path)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)
}
