package cache_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/treaster/codepatterns/cache"
)

func TestFiniteStore(t *testing.T) {
	fs := cache.NewMemStore(2)

	// Get a nonexistent key
	_, err := fs.Get("key1")
	require.Error(t, err)

	// Set a key, then check that it is retrievable
	err = fs.Set("key1", "val1")
	require.NoError(t, err)

	val, err := fs.Get("key1")
	require.NoError(t, err)
	require.Equal(t, "val1", val)

	// Set another key, then check that all keys are retrievable
	err = fs.Set("key2", "val2")
	require.NoError(t, err)

	val, err = fs.Get("key1")
	require.NoError(t, err)
	require.Equal(t, "val1", val)

	val, err = fs.Get("key2")
	require.NoError(t, err)
	require.Equal(t, "val2", val)

	// Set a third key. This will evict the first key created (fifo).
	err = fs.Set("key3", "val3")
	require.NoError(t, err)

	_, err = fs.Get("key1")
	require.Error(t, err)

	val, err = fs.Get("key2")
	require.NoError(t, err)
	require.Equal(t, "val2", val)

	val, err = fs.Get("key3")
	require.NoError(t, err)
	require.Equal(t, "val3", val)
}
