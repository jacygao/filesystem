One interface, many filesystem implementations.

```go
type FileStore interface {
	Get(ctx context.Context, resourceKey string) (io.ReadCloser, error)
	Put(ctx context.Context, resourceKey string, body io.Reader) error
	Delete(ctx context.Context, resourceKey string) error
}
```