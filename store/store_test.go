package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/katallaxie/pkg/store"
)

func TestStore(t *testing.T) {
	// Open a new store
	s, err := store.Open("__deleteme.db", nil)
	assert.NoError(t, err)
	defer s.Close()

	// Put a value
	err = s.Put(store.Byte("key"), []byte("value"))
	assert.NoError(t, err)

	// Get the value
	var val []byte
	err = s.Get(store.Byte("key"), &val)
	assert.NoError(t, err)
	assert.Equal(t, []byte("value"), val)

	// Delete the value
	err = s.Delete(store.Byte("key"))
	assert.NoError(t, err)

	// Try to get the deleted value
	err = s.Get(store.Byte("key"), nil)
	assert.Error(t, err)
	assert.Equal(t, store.ErrKeyNotExist, err)
}
