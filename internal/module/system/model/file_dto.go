package model

import "management-backend/pkg/page"

// FileResp 文件响应
type FileResp struct {
	ID          int64  `json:"id"`
	FileName    string `json:"fileName"`
	FilePath    string `json:"filePath"`
	FileSize    int64  `json:"fileSize"`
	ContentType string `json:"contentType"`
	FileType    string `json:"fileType"`
	BucketName  string `json:"bucketName"`
	StorageType string `json:"storageType"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
	CreateTime  string `json:"createTime"`
}

// FilePageReq 文件分页查询
type FilePageReq struct {
	page.Req
	FileName    string `form:"fileName"`
	FileType    string `form:"fileType"`
	StorageType string `form:"storageType"`
	CreateTime  string `form:"createTime"`
}
