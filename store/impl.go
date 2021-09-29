package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type memImpl struct {
	data map[string]string
}

func (i *memImpl) Get(key string) (string, error) {
	if i.data == nil {
		return "", fmt.Errorf("key %q not found", key)
	}

	value, hasKey := i.data[key]
	if !hasKey {
		return "", fmt.Errorf("key %q not found", key)
	}

	return value, nil
}

func (i *memImpl) Set(key string, value string) error {
	if i.data == nil {
		i.data = map[string]string{}
	}

	i.data[key] = value
	return nil
}

func (i *memImpl) HasKey(key string) bool {
	_, hasKey := i.data[key]
	return hasKey
}

func (i *memImpl) Evict(key string) error {
	delete(i.data, key)
	return nil
}

type fileImpl struct {
	storeDir string
}

func (i *fileImpl) filenameFromKey(key string) string {
	return filepath.Join(i.storeDir, key)
}

func (i *fileImpl) Get(key string) (string, error) {
	filename := i.filenameFromKey(key)
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("file read error: %s", err.Error())
	}

	return string(contents), nil
}

func (i *fileImpl) Set(key string, value string) error {
	filename := i.filenameFromKey(key)
	return ioutil.WriteFile(filename, []byte(value), 0644)
}

func (i *fileImpl) HasKey(key string) bool {
	filename := i.filenameFromKey(key)
	_, err := os.Stat(filename)
	return err == nil
}

func (i *fileImpl) Evict(key string) error {
	filename := i.filenameFromKey(key)
	return os.Remove(filename)
}
