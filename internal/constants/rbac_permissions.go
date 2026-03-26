package constants

// RBAC 权限码常量。
// 说明：统一定义权限码，避免在路由层硬编码字符串导致拼写错误。
const (
	// 更新当前变量或字段值。
	PermRoleList = "sys:role:list"
	// 更新当前变量或字段值。
	PermRoleCreate = "sys:role:create"
	// 更新当前变量或字段值。
	PermRoleView = "sys:role:view"
	// 更新当前变量或字段值。
	PermRoleUpdate = "sys:role:update"
	// 更新当前变量或字段值。
	PermRoleDelete = "sys:role:delete"
	// 更新当前变量或字段值。
	PermRolePermissionView = "sys:role_permission:view"
	// 更新当前变量或字段值。
	PermRolePermissionBind = "sys:role_permission:bind"

	// 更新当前变量或字段值。
	PermPermissionList = "sys:permission:list"
	// 更新当前变量或字段值。
	PermPermissionCreate = "sys:permission:create"
	// 更新当前变量或字段值。
	PermPermissionView = "sys:permission:view"
	// 更新当前变量或字段值。
	PermPermissionUpdate = "sys:permission:update"
	// 更新当前变量或字段值。
	PermPermissionDelete = "sys:permission:delete"

	// 更新当前变量或字段值。
	PermMenuList = "sys:menu:list"
	// 更新当前变量或字段值。
	PermMenuCreate = "sys:menu:create"
	// 更新当前变量或字段值。
	PermMenuView = "sys:menu:view"
	// 更新当前变量或字段值。
	PermMenuUpdate = "sys:menu:update"
	// 更新当前变量或字段值。
	PermMenuDelete = "sys:menu:delete"
	// 更新当前变量或字段值。
	PermMenuFrontendTree = "sys:menu:frontend"
	// 更新当前变量或字段值。
	PermMenuPermissionView = "sys:menu_permission:view"
	// 更新当前变量或字段值。
	PermMenuPermissionBind = "sys:menu_permission:bind"

	// 更新当前变量或字段值。
	PermUserList = "sys:user:list"
	// 更新当前变量或字段值。
	PermUserCreate = "sys:user:create"
	// 更新当前变量或字段值。
	PermUserView = "sys:user:view"
	// 更新当前变量或字段值。
	PermUserUpdate = "sys:user:update"
	// 更新当前变量或字段值。
	PermUserDelete = "sys:user:delete"
	// 更新当前变量或字段值。
	PermUserProfile = "sys:user:profile"

	// 更新当前变量或字段值。
	PermUserRoleView = "sys:user_role:view"
	// 更新当前变量或字段值。
	PermUserRoleBind = "sys:user_role:bind"
	// 更新当前变量或字段值。
	PermUserRoleAdd = "sys:user_role:add"
	// 更新当前变量或字段值。
	PermUserRoleRemove = "sys:user_role:remove"
	// 更新当前变量或字段值。
	PermAuditLoginLogList = "sys:audit:login_log:list"

	// 更新当前变量或字段值。
	PermClientUserList = "client:user:list"
	// 更新当前变量或字段值。
	PermClientPostList = "client:post:list"
	// 更新当前变量或字段值。
	PermClientHotCommentList = "client:comment:hot:list"

	// 更新当前变量或字段值。
	PermBizBannerList = "biz:banner:list"
	// 更新当前变量或字段值。
	PermBizBroadcastList = "biz:broadcast:list"
	// 更新当前变量或字段值。
	PermBizSpecialLotteryList = "biz:special_lottery:list"
	// 更新当前变量或字段值。
	PermBizLotteryInfoList = "biz:lottery_info:list"
	// 更新当前变量或字段值。
	PermBizOfficialPostList = "biz:official_post:list"
	// 更新当前变量或字段值。
	PermBizExternalLinkList = "biz:external_link:list"
	// 更新当前变量或字段值。
	PermBizHomePopupList = "biz:home_popup:list"
	// 更新当前变量或字段值。
	PermBizSMSChannelList = "biz:sms_channel:list"
)
