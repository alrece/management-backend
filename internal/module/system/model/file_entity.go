package model

import "management-backend/internal/model"

// File 文件信息实体
type File struct {
	model.BaseEntity
	model.StatusModel
	model.RemarkModel
	FileName    string `gorm:"type:varchar(255)" json:"fileName"`
	FilePath    string `gorm:"type:varchar(500)" json:"filePath"`
	FileSize    int64  `json:"fileSize"`
	ContentType string `gorm:"type:varchar(100)" json:"contentType"`
	FileType    string `gorm:"type:varchar(20)" json:"fileType"`   // doc/image/video/audio/other
	BucketName  string `gorm:"type:varchar(50)" json:"bucketName"`
	StorageType string `gorm:"type:varchar(20)" json:"storageType"` // minio/local
}

func (File) TableName() string {
	return "sys_file"
}
