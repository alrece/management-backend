package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage 本地存储实现
type LocalStorage struct {
	baseDir string
	baseURL string
}

// NewLocalStorage 创建本地存储
func NewLocalStorage(baseDir, baseURL string) *LocalStorage {
	return &LocalStorage{baseDir: baseDir, baseURL: baseURL}
}

func (s *LocalStorage) Upload(_ context.Context, key string, reader io.Reader, _ int64, _ string) error {
	path := filepath.Join(s.baseDir, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}

func (s *LocalStorage) Download(_ context.Context, key string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.baseDir, key))
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
	return os.Remove(filepath.Join(s.baseDir, key))
}

func (s *LocalStorage) GetURL(key string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}
