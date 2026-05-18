package storage

import (
	"fmt"
	"sync"

	"management-backend/internal/config"
)

var (
	provider Provider
	once     sync.Once
)

// Init 根据配置初始化存储后端
func Init() error {
	var err error
	once.Do(func() {
		switch config.C.OSS.Type {
		case "minio":
			provider, err = NewMinioStorage()
		case "local":
			provider = NewLocalStorage("./uploads", "/uploads")
		default:
			err = fmt.Errorf("不支持的存储类型: %s", config.C.OSS.Type)
		}
	})
	return err
}

// GetProvider 获取当前存储实例
func GetProvider() Provider {
	return provider
}

// SwitchProvider 动态切换存储后端（sys_param 触发）
func SwitchProvider(p Provider) {
	provider = p
}

// ProviderType 返回当前存储类型
func ProviderType() string {
	if provider == nil {
		return ""
	}
	return config.C.OSS.Type
}
