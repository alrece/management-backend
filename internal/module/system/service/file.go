package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"management-backend/internal/model"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/storage"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"

	"github.com/gin-gonic/gin"
)

// FileService 文件管理接口
type FileService interface {
	Upload(c *gin.Context) (*smodel.File, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*smodel.FileResp, error)
	Page(ctx context.Context, req *smodel.FilePageReq) ([]smodel.FileResp, int64, error)
	Download(ctx context.Context, id int64) (io.ReadCloser, *smodel.File, error)
}

type fileRepo interface {
	Create(ctx context.Context, file *smodel.File) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.File, error)
	Page(ctx context.Context, req *smodel.FilePageReq) ([]smodel.File, int64, error)
}

type fileSvc struct {
	repo fileRepo
}

// NewFileService 创建文件服务
func NewFileService(repo fileRepo) FileService {
	return &fileSvc{repo: repo}
}

// fileTypeFromMime 从 content-type 推断文件类型
func fileTypeFromMime(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "image/"):
		return "image"
	case strings.HasPrefix(contentType, "video/"):
		return "video"
	case strings.HasPrefix(contentType, "audio/"):
		return "audio"
	case strings.Contains(contentType, "pdf"), strings.Contains(contentType, "document"),
		strings.Contains(contentType, "spreadsheet"), strings.Contains(contentType, "text/"):
		return "doc"
	default:
		return "other"
	}
}

func (s *fileSvc) Upload(c *gin.Context) (*smodel.File, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %w", err)
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 构建存储路径: yyyy-MM-dd/雪花ID.ext
	ext := filepath.Ext(header.Filename)
	date := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s/%d%s", date, snowflake.NextID(), ext)

	p := storage.GetProvider()
	if p == nil {
		return nil, fmt.Errorf("存储服务未初始化")
	}

	if err := p.Upload(c.Request.Context(), key, file, header.Size, contentType); err != nil {
		return nil, fmt.Errorf("文件上传失败: %w", err)
	}

	record := &smodel.File{
		BaseEntity:  model.BaseEntity{ID: snowflake.NextID()},
		FileName:    header.Filename,
		FilePath:    key,
		FileSize:    header.Size,
		ContentType: contentType,
		FileType:    fileTypeFromMime(contentType),
		StorageType: "minio",
	}
	record.Status = 0

	if err := s.repo.Create(c.Request.Context(), record); err != nil {
		return nil, fmt.Errorf("文件记录保存失败: %w", err)
	}

	return record, nil
}

func (s *fileSvc) Delete(ctx context.Context, id int64) error {
	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errcode.Err(errcode.FileNotFound)
	}

	p := storage.GetProvider()
	if p != nil {
		_ = p.Delete(ctx, file.FilePath)
	}

	return s.repo.Delete(ctx, id)
}

func (s *fileSvc) Get(ctx context.Context, id int64) (*smodel.FileResp, error) {
	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.Err(errcode.FileNotFound)
	}
	return toFileResp(file), nil
}

func (s *fileSvc) Page(ctx context.Context, req *smodel.FilePageReq) ([]smodel.FileResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]smodel.FileResp, 0, len(list))
	for i := range list {
		resps = append(resps, *toFileResp(&list[i]))
	}
	return resps, total, nil
}

func (s *fileSvc) Download(ctx context.Context, id int64) (io.ReadCloser, *smodel.File, error) {
	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, errcode.Err(errcode.FileNotFound)
	}

	p := storage.GetProvider()
	if p == nil {
		return nil, nil, fmt.Errorf("存储服务未初始化")
	}

	reader, err := p.Download(ctx, file.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("文件下载失败: %w", err)
	}
	return reader, file, nil
}

func toFileResp(f *smodel.File) *smodel.FileResp {
	return &smodel.FileResp{
		ID:          f.ID,
		FileName:    f.FileName,
		FilePath:    f.FilePath,
		FileSize:    f.FileSize,
		ContentType: f.ContentType,
		FileType:    f.FileType,
		BucketName:  f.BucketName,
		StorageType: f.StorageType,
		Status:      f.Status,
		Remark:      f.Remark,
		CreateTime:  f.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
