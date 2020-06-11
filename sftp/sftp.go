package sftp

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

type sftpBacked struct {
	sftp *sftp.Client
	// Filesystem path where the resources are stored.
	path string
}

// NewSftpBacked initialises and returns a new sftpBacked instance.
func NewSftpBacked(cli *sftp.Client, path string) *sftpBacked {
	return &sftpBacked{
		sftp: cli,
		path: path,
	}
}

// getResourcePath returns the full filesystem path for where the data for resourceKey will be stored.
func (s *sftpBacked) getResourcePath(resourceKey string) string {
	return filepath.Join(s.path, resourceKey)
}

// Get attemps to retrieve file buffer from sftp path.
func (s *sftpBacked) Get(ctx context.Context, resourceKey string) (io.ReadCloser, error) {
	file, err := s.sftp.Open(s.getResourcePath(resourceKey))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Put attemps to save file buffer in sftp path.
func (s *sftpBacked) Put(ctx context.Context, resourceKey string, body io.Reader) error {
	f, err := s.sftp.Create(s.getResourcePath(resourceKey))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// Delete removes file from sftp path.
func (s *sftpBacked) Delete(ctx context.Context, resourceKey string) error {
	return s.sftp.Remove(s.getResourcePath(resourceKey))
}

// List reads the directory named by dir and returns a list of directory entries.
func (s *sftpBacked) List(ctx context.Context, dir string) ([]os.FileInfo, error) {
	return s.sftp.ReadDir(s.getResourcePath(dir))
}
