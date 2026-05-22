package errcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsg_KnownCodes(t *testing.T) {
	tests := []struct {
		code int
		msg  string
	}{
		{Success, "操作成功"},
		{BadRequest, "参数错误"},
		{Unauthorized, "未授权，请先登录"},
		{Forbidden, "没有权限访问该资源"},
		{NotFound, "资源不存在"},
		{InternalError, "服务器内部错误"},
		{UserNotFound, "用户不存在"},
		{UserPasswordErr, "密码错误"},
		{UserDisabled, "用户已被禁用"},
		{UserExists, "用户名已存在"},
		{TokenExpired, "Token 已过期"},
		{TokenInvalid, "Token 无效"},
		{LoginLocked, "账号已被锁定"},
		{RoleNotFound, "角色不存在"},
		{MenuNotFound, "菜单不存在"},
		{DeptNotFound, "部门不存在"},
		{TenantNotFound, "租户不存在"},
		{TenantDisabled, "租户已被禁用"},
		{PermissionDenied, "权限不足"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.msg, Msg(tt.code), "code=%d", tt.code)
	}
}

func TestMsg_UnknownCode(t *testing.T) {
	assert.Equal(t, "未知错误", Msg(99999))
}
