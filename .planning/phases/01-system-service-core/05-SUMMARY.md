---
plan: 05
phase: 01-system-service-core
status: complete
started: 2026-05-18
completed: 2026-05-18
self_check: passed
---

# Plan 05: System CRUD 模块

## Summary

完成系统管理全套 CRUD：用户管理完善（多租户数据权限注入）、个人中心（改密吊销 Token）、用户导出（excelize StreamWriter 10万行限制）、字典管理（Redis 缓存 dict:{type}）、参数管理（Redis 缓存 param:{key}）、通知公告 CRUD、操作日志中间件 + 登录日志 + 4 个日志 API 端点。

## Tasks Completed

| Task | Description | Status |
|------|-------------|--------|
| 5.1 | 用户管理 CRUD（多租户 + 数据权限） | ✓ |
| 5.2 | 个人中心 + 用户导入导出 | ✓ |
| 5.3 | 字典管理（Redis 缓存） | ✓ |
| 5.4 | 参数管理（Redis 缓存） | ✓ |
| 5.5 | 通知公告 | ✓ |
| 5.6 | 操作日志 + 登录日志 | ✓ |

## Key Files

### Created
- internal/module/system/model/dict_entity.go — DictType/DictData 实体
- internal/module/system/model/param_entity.go — SysParam 实体
- internal/module/system/model/notice_entity.go — Notice + SysOperLog + SysLoginLog 实体
- internal/module/system/model/dict_dto.go — 字典 DTO
- internal/module/system/model/param_dto.go — 参数 DTO
- internal/module/system/model/notice_dto.go — 通知 DTO
- internal/module/system/repository/dict_type.go — DictTypeRepo
- internal/module/system/repository/dict_data.go — DictDataRepo 含 ListByDictType
- internal/module/system/repository/param.go — ParamRepo 含 GetByKey
- internal/module/system/repository/notice.go — NoticeRepo
- internal/module/system/repository/log.go — OperLogRepo + LoginLogRepo
- internal/module/system/service/dict.go — DictService 含 Redis 缓存
- internal/module/system/service/param.go — ParamService 含 Redis 缓存
- internal/module/system/service/notice.go — NoticeService
- internal/module/system/service/log.go — LogService
- internal/module/system/service/user_export.go — UserExportService StreamWriter 10万行
- internal/module/system/handler/dict.go — DictHandler
- internal/module/system/handler/param.go — ParamHandler
- internal/module/system/handler/notice.go — NoticeHandler
- internal/module/system/handler/log.go — LogHandler
- internal/module/system/handler/profile.go — ProfileHandler 含改密+RevokeByUser
- internal/middleware/operation_log.go — OperationLog 中间件

## Self-Check

- [x] 用户 CRUD（跨租户越权修复）
- [x] 个人中心（改密吊销 Token）
- [x] 用户导出（StreamWriter, 10万行限制）
- [x] 字典管理（Redis 缓存）
- [x] 参数管理（Redis 缓存）
- [x] 操作日志 + 登录日志
