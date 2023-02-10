package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempDir(t *testing.T) {
	f, fn, err := TempDir()
	assert.NoError(t, err)
	defer fn()

	_, err = os.Stat(f.Name())
	os.IsNotExist(err)
	assert.NoError(t, err)
}
