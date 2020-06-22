/*
	Package memory implements the filesystem interface backed by a in memory file storage.
	This package is designed for testing purposes to act as a fake file storage based on a simple key value map
*/

package memory

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

// A fake implementation of the S3 interface based on a simple key value map.
type Item struct {
	key  string
	data []byte
}

type MemoryStore struct {
	// in memory key/value map to store data
	assets map[string][]Item
	// a prefix string of all keys. root by default
	root string
	// Mutex to synchronise all access
	mutex sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		assets: make(map[string][]Item),
		root:   "root",
	}
}

// Get retrieves files from file storage
func (ms *MemoryStore) Get(ctx context.Context, resourceKey string) (io.ReadCloser, error) {

	var err error

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	dir, key := ms.getResources(resourceKey)

	items, exists := ms.assets[dir]
	if !exists {
		err = errors.New(fmt.Sprintf("Key %s does not exist", resourceKey))
		return nil, err
	}

	for _, item := range items {
		if item.key == key {
			return ioutil.NopCloser(bytes.NewReader(item.data)), nil
		}
	}

	err = errors.New(fmt.Sprintf("Key %s does not exist", resourceKey))
	return nil, err
}

// ListFileNames returns a list of file names of files contained in file storage
func (ms *MemoryStore) ListFileNames(ctx context.Context, path string) ([]string, error) {

	var err error

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	path = filepath.Join(ms.root, path)
	if len(ms.assets[path]) < 1 {
		err = errors.New(fmt.Sprintf("Key %s does not exist", path))
		return nil, err
	}

	var payload []string

	for _, data := range ms.assets[path] {
		payload = append(payload, data.key)
	}

	// We're assuming that the fake service cannot fail for now.
	return payload, nil
}

// Put writes file to the  file storage
func (ms *MemoryStore) Put(ctx context.Context, resourceKey string, body io.Reader) error {

	// read in body into our own array
	buffer, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	dir, key := ms.getResources(resourceKey)

	item := Item{
		key:  key,
		data: buffer,
	}
	for i, val := range ms.assets[dir] {
		if val.key == key {
			ms.assets[dir][i] = item
			return nil
		}
	}
	ms.assets[dir] = append(ms.assets[dir], item)

	return nil
}

// Delete removed files from file storage
func (ms *MemoryStore) Delete(ctx context.Context, resourceKey string) error {

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	dir, key := ms.getResources(resourceKey)

	if len(ms.assets[dir]) < 1 {
		err := errors.New(fmt.Sprintf("Key %s does not exist", resourceKey))
		return err
	}

	for _, item := range ms.assets[dir] {
		if item.key == key {
			ms.assets[dir] = pop(ms.assets[dir], item)
		}
	}

	// We're assuming that the fake service cannot fail for now.
	return nil
}

// getResources takes a resourceKey and returns the full filepath and full filename.
func (ms *MemoryStore) getResources(resourceKey string) (string, string) {
	dir, key := filepath.Split(resourceKey)
	return filepath.Join(ms.root, dir), key
}

// pop removes an Item from an Item slice
func pop(items []Item, item Item) []Item {
	for i, v := range items {
		if v.key == item.key {
			return append(items[:i], items[i+1:]...)
		}
	}
	return items
}

// trimPath removes forward slash in the end of a given path string
// in memory storage uses path as the key, adding a slash in the end makes the keys messy
// as it can accept both key with or without forward slash in the end of a path string
func trimPath(path string) string {
	return strings.TrimSuffix(path, "/")
}
