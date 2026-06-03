package rolepermissions

import (
	"time"
)

// RolePermission represents the junction table between roles and permissions
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;column:role_id" json:"role_id"`
	PermissionID uint `gorm:"primaryKey;column:permission_id" json:"permission_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName returns the table name for RolePermission
func (RolePermission) TableName() string {
	return "role_permissions"
}
