package storage

import (
	"context"
	"io"
)

// Provider 文件存储接口
type Provider interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}
