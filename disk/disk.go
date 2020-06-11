package disk

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// diskBacked implements of the Filestore interface backed by a disk filesystem.
type diskBacked struct {
	// Filesystem path where the resources are stored.
	path string
}

// NewDiskStore initialises and returns a new diskBacked instance.
func NewDiskStore(path string) (*diskBacked, error) {
	// create upload directory (if it doesn't already exist)
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	return &diskBacked{
		path: path,
	}, nil
}

// getResourcePath returns the full filesystem path for where the data for resourceKey will be stored.
func (s *diskBacked) getResourcePath(resourceKey string) string {
	return filepath.Join(s.path, resourceKey)
}

// Get attemps to retrieve file buffer from disk path.
func (s *diskBacked) Get(ctx context.Context, resourceKey string) (io.ReadCloser, error) {
	file, err := os.Open(s.getResourcePath(resourceKey))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Put attempts to save file buffer in disk path.
func (s *diskBacked) Put(ctx context.Context, resourceKey string, body io.Reader) error {
	r := s.getResourcePath(resourceKey)
	dir, _ := filepath.Split(r)
	// create child directory (if it doesn't already exist)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write to a temp file then rename in order for the write operation to be atomic.
	f, err := ioutil.TempFile(s.path, "tmp")
	if err != nil {
		return err
	}
	defer f.Close()

	// Required cleanup, if we happen to fail before we can rename the temp file
	defer os.Remove(f.Name())

	// io.Copy will copy from r.Body to s3Writer until either EOF is reached on
	// src or an error occurs. If EOF is reached then a nil error is returned
	if _, err := io.Copy(f, body); err != nil {
		return err
	}

	if err := os.Rename(f.Name(), r); err != nil {
		return err
	}

	return nil
}

// Delete removes file from disk path.
func (s *diskBacked) Delete(ctx context.Context, resourceKey string) error {
	err := os.Remove(s.getResourcePath(resourceKey))
	if err != nil {
		return err
	}

	return nil
}

// List returns a list of files
func (s *diskBacked) List(ctx context.Context, dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(s.getResourcePath(dir))
	if err != nil {
		return nil, err
	}
	return files, nil
}
