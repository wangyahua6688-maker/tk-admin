package constants

// tk-admin 业务状态码定义（仅后台管理服务使用）。
//
// 约定：
// 1) 0 代表成功；
// 2) 41xxx 代表后台业务参数/流程错误；
// 3) 51xxx 代表后台系统错误；
// 4) 42xxx 代表后台管理权限/认证错误。

// 通用成功码。
const (
	// CodeOK 表示请求成功。
	CodeOK = 0
)

// 业务码（后台业务逻辑）。
const (
	// AdminBizInvalidRequest 请求参数不合法。
	AdminBizInvalidRequest = 41001
	// AdminBizEmptyUpdate 更新内容为空。
	AdminBizEmptyUpdate = 41002
	// AdminBizResourceNotFound 业务资源不存在。
	AdminBizResourceNotFound = 41004
)

// 系统码（后台基础设施/依赖错误）。
const (
	// AdminSysInternalError 服务内部错误。
	AdminSysInternalError = 51001
	// AdminSysDatabaseError 数据库错误。
	AdminSysDatabaseError = 51002
	// AdminSysRedisError Redis 错误。
	AdminSysRedisError = 51003
)

// 管理后台认证与权限码。
const (
	// AdminAuthUnauthorized 未登录或会话无效。
	AdminAuthUnauthorized = 42001
	// AdminAuthTokenInvalid Token 无效或已过期。
	AdminAuthTokenInvalid = 42002
	// AdminAuthForbidden 权限不足。
	AdminAuthForbidden = 42003
	// AdminAuthUserDisabled 账号被禁用或不存在。
	AdminAuthUserDisabled = 42004
	// AdminAuthRateLimited 认证相关请求触发限流。
	AdminAuthRateLimited = 42005
)
