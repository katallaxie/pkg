package files

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMkdirAll(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	assert.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	path := strings.Join([]string{tempDir, "example"}, "/")

	err = MkdirAll(path, 0o755)
	assert.NoError(t, err)

	_, err = os.Stat(path)
	assert.NoError(t, err)
}
