package store

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImpls(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "store_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	impls := []storeImpl{
		&memImpl{},
		&fileImpl{tmpDir},
	}

	// Test that all storeImpls provide the same basic behavior.
	for _, impl := range impls {
		// Get an unset key
		require.False(t, impl.HasKey("key1"))
		_, err := impl.Get("key1")
		require.Error(t, err)

		// Set the key, then Get again.
		err = impl.Set("key1", "val1")
		require.NoError(t, err)
		require.True(t, impl.HasKey("key1"))
		val, err := impl.Get("key1")
		require.NoError(t, err)
		require.Equal(t, "val1", val)

		// Get another unset key
		require.False(t, impl.HasKey("key2"))
		_, err = impl.Get("key2")
		require.Error(t, err)

		// Set the second key, then verify both the first and second keys.
		err = impl.Set("key2", "val2")
		require.NoError(t, err)

		require.True(t, impl.HasKey("key1"))
		val, err = impl.Get("key1")
		require.NoError(t, err)
		require.Equal(t, "val1", val)

		require.True(t, impl.HasKey("key2"))
		val, err = impl.Get("key2")
		require.NoError(t, err)
		require.Equal(t, "val2", val)

		// Overwrite a key. Make sure the new value sticks.
		err = impl.Set("key2", "val2-2")
		require.NoError(t, err)

		require.True(t, impl.HasKey("key2"))
		val, err = impl.Get("key2")
		require.NoError(t, err)
		require.Equal(t, "val2-2", val)

		// Evict a key. Make sure it becomes unavailable while other keys remain available.
		err = impl.Evict("key2")
		require.NoError(t, err)

		require.False(t, impl.HasKey("key2"))
		_, err = impl.Get("key2")
		require.Error(t, err)

		require.True(t, impl.HasKey("key1"))
		val, err = impl.Get("key1")
		require.NoError(t, err)
		require.Equal(t, "val1", val)
	}
}
