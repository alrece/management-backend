package handler

import (
	"context"
	"net/http"

	"management-backend/internal/auth"
	"management-backend/internal/middleware"
	"management-backend/internal/module/system/model"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler 认证接口处理器
type AuthHandler struct {
	svc       userQuerier
	loginLock *auth.LoginLock
	logger    *zap.Logger
}

// userQuerier 用户查询接口（避免循环依赖）
type userQuerier interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

// NewAuthHandler 构造函数
func NewAuthHandler(svc userQuerier, loginLock *auth.LoginLock, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{svc: svc, loginLock: loginLock, logger: logger}
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Login 登录 POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}

	// 检查账号锁定
	locked, _ := h.loginLock.IsLocked(c.Request.Context(), req.Username)
	if locked {
		remaining := h.loginLock.GetLockRemaining(c.Request.Context(), req.Username)
		response.Fail(c, http.StatusLocked, "账号已锁定，请 "+remaining.String()+" 后重试")
		return
	}

	// 查找用户
	user, err := h.svc.GetByUsername(c.Request.Context(), req.Username)
	if err != nil || user == nil {
		_ = h.loginLock.RecordFail(c.Request.Context(), req.Username)
		response.Fail(c, 401, "用户名或密码错误")
		return
	}

	// 校验密码
	if !middleware.CheckPassword(req.Password, user.Password) {
		_ = h.loginLock.RecordFail(c.Request.Context(), req.Username)
		response.Fail(c, 401, "用户名或密码错误")
		return
	}

	// 登录成功，重置失败计数
	_ = h.loginLock.ResetFail(c.Request.Context(), req.Username)

	// 生成双 Token
	pair, err := auth.GenerateTokenPair(user.ID, user.Username, user.TenantID)
	if err != nil {
		h.logger.Error("生成 Token 失败", zap.Error(err))
		response.Fail(c, 500, "登录失败")
		return
	}

	response.Ok(c, pair)
}

// RefreshToken 刷新 Token POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数校验失败")
		return
	}

	claims, err := auth.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		response.Fail(c, 401, "Refresh Token 无效或已过期")
		return
	}

	// 生成新的双 Token
	pair, err := auth.GenerateTokenPair(claims.UserID, claims.Username, claims.TenantID)
	if err != nil {
		response.Fail(c, 500, "刷新失败")
		return
	}

	response.Ok(c, pair)
}

// Logout 登出 POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	claims, err := auth.ParseAccessToken(tokenStr)
	if err != nil {
		response.Fail(c, 401, "Token 无效")
		return
	}

	// 将 JTI 加入黑名单
	if err := middleware.RevokeToken(c.Request.Context(), claims.JTI); err != nil {
		h.logger.Error("吊销 Token 失败", zap.Error(err))
		response.Fail(c, 500, "登出失败")
		return
	}

	response.OkMsg(c, "登出成功")
}
