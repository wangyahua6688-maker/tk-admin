package models

import "time"

// SystemMessage 系统消息模型。
// 设计说明：
// 1. 一条消息只对应一个接收用户，便于做“已读/未读”管理；
// 2. 通过 biz_type + biz_id 记录业务来源，便于后续回溯来源（用户/角色/权限）。
type SystemMessage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index:idx_sys_msg_user_read,priority:1" json:"user_id"`
	Title      string    `gorm:"size:200;not null" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Level      string    `gorm:"size:20;not null;default:'info'" json:"level"` // info/success/warning
	IsRead     bool      `gorm:"not null;default:false;index:idx_sys_msg_user_read,priority:2" json:"is_read"`
	OperatorID uint      `gorm:"not null;default:0;index" json:"operator_id"`
	BizType    string    `gorm:"size:50;default:'';index" json:"biz_type"`
	BizID      uint      `gorm:"not null;default:0;index" json:"biz_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 指定系统消息表名。
func (SystemMessage) TableName() string {
	return "sys_system_messages"
}
