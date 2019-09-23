package filesystem

import (
	"context"
	"io"
)

// Filesystem defines common behaviours of filesystem.
// This interface can be implemented by different backends.
type Filesystem interface {
	Get(ctx context.Context, resourceKey string) (io.ReadCloser, error)
	Put(ctx context.Context, resourceKey string, body io.Reader) error
	Delete(ctx context.Context, resourceKey string) error
}
