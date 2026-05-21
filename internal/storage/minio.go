package storage

import (
	"context"
	"fmt"
	"io"

	"management-backend/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage MinIO 存储实现
type MinioStorage struct {
	client *minio.Client
	bucket string
}

// NewMinioStorage 创建 MinIO 存储
func NewMinioStorage() (*MinioStorage, error) {
	cfg := config.C.OSS.Minio
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("MinIO 客户端创建失败: %w", err)
	}

	ctx := context.Background()
	if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
		if exists, _ := client.BucketExists(ctx, cfg.Bucket); !exists {
			return nil, fmt.Errorf("MinIO bucket 创建失败: %w", err)
		}
	}

	return &MinioStorage{client: client, bucket: cfg.Bucket}, nil
}

func (s *MinioStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (s *MinioStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
}

func (s *MinioStorage) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

func (s *MinioStorage) GetURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", config.C.OSS.Minio.Endpoint, s.bucket, key)
}
