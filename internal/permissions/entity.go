package permissions

import "time"

// Permission memetakan tabel `permissions` di database popda_2026.
// Setiap permission terikat ke satu module via module_id.
// Format name: <module_code>.<action>, contoh: cabor.read, atlet.create
type Permission struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	ModuleID    *uint     `gorm:"column:module_id" json:"module_id"`
	Name        string    `gorm:"size:100;not null;column:name" json:"name"`
	Description string    `gorm:"type:text;column:description" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

// CreatePermissionRequest adalah body untuk POST /admin/permissions
type CreatePermissionRequest struct {
	ModuleID    *uint  `json:"module_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdatePermissionRequest adalah body untuk PUT /admin/permissions/:id
type UpdatePermissionRequest struct {
	ModuleID    *uint  `json:"module_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
