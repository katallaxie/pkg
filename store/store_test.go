package store_test

import (
	"testing"

	"github.com/katallaxie/pkg/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	// Open a new store
	s, err := store.Open("__deleteme.db", nil)
	require.NoError(t, err)
	defer s.Close()

	// Put a value
	err = s.Put(store.Byte("key"), []byte("value"))
	require.NoError(t, err)

	// Get the value
	var val []byte
	err = s.Get(store.Byte("key"), &val)
	require.NoError(t, err)
	assert.Equal(t, []byte("value"), val)

	// Delete the value
	err = s.Delete(store.Byte("key"))
	require.NoError(t, err)

	// Try to get the deleted value
	err = s.Get(store.Byte("key"), nil)
	require.Error(t, err)
	assert.Equal(t, store.ErrKeyNotExist, err)
}
