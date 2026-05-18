package model

import "management-backend/internal/model"

// Client 客户端管理实体
type Client struct {
	model.BaseModel
	model.StatusModel
	model.RemarkModel
	ClientID             string `gorm:"type:varchar(50);uniqueIndex" json:"clientId"`
	ClientKey            string `gorm:"type:varchar(100);not null" json:"clientKey"`
	ClientName           string `gorm:"type:varchar(100)" json:"clientName"`
	GrantType            string `gorm:"type:varchar(50);default:authorization_code" json:"grantType"`
	RedirectURI          string `gorm:"type:varchar(500)" json:"redirectUri"`
	AccessTokenValidity  int64  `gorm:"default:1800" json:"accessTokenValidity"`
	RefreshTokenValidity int64  `gorm:"default:604800" json:"refreshTokenValidity"`
}

func (Client) TableName() string {
	return "sys_client"
}
