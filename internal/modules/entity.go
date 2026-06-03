package modules

import "time"

// Module memetakan tabel `modules` di database popda_2026.
// Setiap modul merepresentasikan satu fitur/halaman dalam sistem.
type Module struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	Name      string    `gorm:"size:100;not null;column:name" json:"name"`
	Label     string    `gorm:"size:150;not null;column:label" json:"label"`
	Code      string    `gorm:"size:50;not null;uniqueIndex;column:code" json:"code"`
	URL       string    `gorm:"size:255;column:url" json:"url"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (Module) TableName() string {
	return "modules"
}

// CreateModuleRequest adalah body untuk POST /admin/modules
type CreateModuleRequest struct {
	Name  string `json:"name" binding:"required"`
	Label string `json:"label" binding:"required"`
	Code  string `json:"code" binding:"required"`
	URL   string `json:"url"`
}

// UpdateModuleRequest adalah body untuk PUT /admin/modules/:id
type UpdateModuleRequest struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Code  string `json:"code"`
	URL   string `json:"url"`
}
