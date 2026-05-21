package handler

import (
	"management-backend/internal/auth"
	"management-backend/internal/middleware"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/errcode"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ProfileHandler 个人中心处理器
type ProfileHandler struct {
	userSvc   service.UserService
	blacklist *auth.Blacklist
}

func NewProfileHandler(userSvc service.UserService, blacklist *auth.Blacklist) *ProfileHandler {
	return &ProfileHandler{userSvc: userSvc, blacklist: blacklist}
}

// GetProfile 获取个人信息 GET /system/user/profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	data, err := h.userSvc.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, data)
}

type updateProfileReq struct {
	Nickname string `json:"nickname" binding:"omitempty,max=30"`
	Email    string `json:"email" binding:"omitempty,email,max=50"`
	Mobile   string `json:"mobile" binding:"omitempty,max=20"`
	Sex      int    `json:"sex" binding:"omitempty,oneof=0 1 2"`
}

// UpdateProfile 更新个人信息 PUT /system/user/profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	updateReq := &service.UserProfileUpdateReq{
		ID:       userID,
		Nickname: req.Nickname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Sex:      req.Sex,
	}

	if err := h.userSvc.UpdateProfile(c.Request.Context(), updateReq); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "更新成功")
}

type updatePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6,max=50"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=50"`
}

// UpdatePassword 修改密码 PUT /system/user/password
func (h *ProfileHandler) UpdatePassword(c *gin.Context) {
	var req updatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userSvc.GetRawByID(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, 500, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		response.Fail(c, errcode.UserPasswordErr, "旧密码错误")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, 500, "密码加密失败")
		return
	}

	if err := h.userSvc.UpdatePassword(c.Request.Context(), userID, string(hashed)); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	// 吊销该用户所有 Token
	_ = h.blacklist.RevokeByUser(c.Request.Context(), userID)

	response.OkMsg(c, "密码修改成功，请重新登录")
}
