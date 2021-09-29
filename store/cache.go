package store

import (
	"fmt"
	"time"
)

// Return a FiniteStore that's backed with an in-memory storage.
func NewMemStore(maxSize int) FiniteStore {
	return &finiteStore{
		&memImpl{},
		map[string]time.Time{},
		maxSize,
	}
}

// Return a FinieStore that's backed with a filesystem-based storage.
func NewFileStore(storeDir string, maxSize int) FiniteStore {
	return &finiteStore{
		&fileImpl{},
		map[string]time.Time{},
		maxSize,
	}
}

type FiniteStore interface {
	Get(string) (string, error)
	Set(string, string) error
}

type storeImpl interface {
	Get(string) (string, error)
	Set(string, string) error
	HasKey(string) bool
	Evict(string) error
}

type finiteStore struct {
	impl         storeImpl
	fifoTracker  map[string]time.Time
	maxStoreSize int
}

func (c *finiteStore) Get(key string) (string, error) {
	value, err := fs.impl.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *finiteStore) Set(key string, value string) error {
	// If the store is full and the new key is not in the store, evict the oldest key
	if !fs.impl.HasKey(key) {
		if len(fs.fifoTracker) >= fs.maxStoreSize {
			oldestT := time.Now()
			var oldestKey string
			for k, t := range fs.fifoTracker {
				if t.Before(oldestT) {
					oldestT = t
					oldestKey = k
				}
			}

			err := fs.impl.Evict(oldestKey)
			if err != nil {
				return fmt.Errorf("unable to evict old key from store: %s", err.Error())
			}
			delete(fs.fifoTracker, oldestKey)
		}

		// Set the time on the new key, when it is created.
		fs.fifoTracker[key] = time.Now()
	}

	err := fs.impl.Set(key, value)
	if err != nil {
		return err
	}

	return nil
}
