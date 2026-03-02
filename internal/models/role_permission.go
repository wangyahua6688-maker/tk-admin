package models

// RolePermission 角色-权限关联模型。
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;column:role_id"`
	PermissionID uint `gorm:"primaryKey;column:permission_id"`
}

func (RolePermission) TableName() string {
	return "sys_role_permissions"
}
