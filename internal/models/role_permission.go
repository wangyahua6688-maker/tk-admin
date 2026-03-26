package models

// RolePermission 角色-权限关联模型。
type RolePermission struct {
	// 处理当前语句逻辑。
	RoleID uint `gorm:"primaryKey;column:role_id"`
	// 处理当前语句逻辑。
	PermissionID uint `gorm:"primaryKey;column:permission_id"`
}

// TableName 返回模型对应的数据表名。
func (RolePermission) TableName() string {
	return "sys_role_permissions"
}
