package errcode

// 错误码分段定义
const (
	Success         = 0
	BadRequest      = 400
	Unauthorized    = 401
	Forbidden       = 403
	NotFound        = 404
	TooManyRequests = 429
	InternalError   = 500

	// 系统模块 1000-1999
	UserNotFound     = 1001
	UserPasswordErr  = 1002
	UserDisabled     = 1003
	UserExists       = 1004
	PasswordWeak     = 1005
	PasswordExpired  = 1006
	TokenInvalid     = 1041
	TokenExpired     = 1040
	LoginLocked      = 1007
	LoginFailLimit   = 1008
	RoleNotFound     = 1010
	RoleExists       = 1011
	MenuNotFound     = 1012
	DeptNotFound     = 1013
		PostNotFound     = 1014
		PostExists       = 1015
		MenuHasChildren  = 1020
		MenuHasRoleRef   = 1021
		DeptHasChildren  = 1022
		DeptHasUsers     = 1023
		RoleHasUsers     = 1024
		PermissionDenied = 1401

	// 基础设施 2000-2999
	DictNotFound    = 2001
	ParamNotFound   = 2002
	FileUploadFail  = 2003
	FileNotFound    = 2004
	ConfigKeyExists = 2010

	// 租户 3000-3999
	TenantNotFound   = 3001
	TenantExpired    = 3002
	TenantDisabled   = 3003
	TenantDBInitFail = 3004
)

var messages = map[int]string{
	Success:         "操作成功",
	BadRequest:      "参数错误",
	Unauthorized:    "未授权，请先登录",
	Forbidden:       "没有权限访问该资源",
	NotFound:        "资源不存在",
	TooManyRequests: "请求过于频繁",
	InternalError:   "服务器内部错误",

	UserNotFound:    "用户不存在",
	UserPasswordErr: "密码错误",
	UserDisabled:    "用户已被禁用",
	UserExists:      "用户名已存在",
	PasswordWeak:    "密码强度不足",
	TokenExpired:    "Token 已过期",
	TokenInvalid:    "Token 无效",
	LoginLocked:     "账号已被锁定",
	LoginFailLimit:  "登录失败次数超限",

	RoleNotFound:     "角色不存在",
		RoleExists:       "角色编码已存在",
		MenuNotFound:     "菜单不存在",
		DeptNotFound:     "部门不存在",
		PostNotFound:     "岗位不存在",
		PostExists:       "岗位编码已存在",
		MenuHasChildren:  "存在子菜单，无法删除",
		MenuHasRoleRef:   "菜单已分配给角色，无法删除",
		DeptHasChildren:  "存在子部门，无法删除",
		DeptHasUsers:     "部门下存在用户，无法删除",
		RoleHasUsers:     "角色下存在用户，无法删除",
		PermissionDenied: "权限不足",
		
		TenantNotFound: "租户不存在",
	TenantDisabled: "租户已被禁用",
	TenantExpired:  "租户已过期",
}

// Msg 获取错误码对应的中文消息
func Msg(code int) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}
