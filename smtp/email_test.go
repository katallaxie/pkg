package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMesage(t *testing.T) {
	t.Parallel()

	m, err := NewMessage()
	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.NotEmpty(t, m.ID)
	assert.Equal(t, map[Header][]string{}, m.Headers)
}
