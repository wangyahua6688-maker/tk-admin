package constants

// RBAC 权限码常量。
// 说明：统一定义权限码，避免在路由层硬编码字符串导致拼写错误。
const (
	PermRoleList           = "sys:role:list"
	PermRoleCreate         = "sys:role:create"
	PermRoleView           = "sys:role:view"
	PermRoleUpdate         = "sys:role:update"
	PermRoleDelete         = "sys:role:delete"
	PermRolePermissionView = "sys:role_permission:view"
	PermRolePermissionBind = "sys:role_permission:bind"

	PermPermissionList   = "sys:permission:list"
	PermPermissionCreate = "sys:permission:create"
	PermPermissionView   = "sys:permission:view"
	PermPermissionUpdate = "sys:permission:update"
	PermPermissionDelete = "sys:permission:delete"

	PermMenuList           = "sys:menu:list"
	PermMenuCreate         = "sys:menu:create"
	PermMenuView           = "sys:menu:view"
	PermMenuUpdate         = "sys:menu:update"
	PermMenuDelete         = "sys:menu:delete"
	PermMenuFrontendTree   = "sys:menu:frontend"
	PermMenuPermissionView = "sys:menu_permission:view"
	PermMenuPermissionBind = "sys:menu_permission:bind"

	PermUserList    = "sys:user:list"
	PermUserCreate  = "sys:user:create"
	PermUserView    = "sys:user:view"
	PermUserUpdate  = "sys:user:update"
	PermUserDelete  = "sys:user:delete"
	PermUserProfile = "sys:user:profile"

	PermUserRoleView      = "sys:user_role:view"
	PermUserRoleBind      = "sys:user_role:bind"
	PermUserRoleAdd       = "sys:user_role:add"
	PermUserRoleRemove    = "sys:user_role:remove"
	PermAuditLoginLogList = "sys:audit:login_log:list"

	PermMessageList    = "sys:message:list"
	PermMessageRead    = "sys:message:read"
	PermMessageReadAll = "sys:message:read_all"
)
