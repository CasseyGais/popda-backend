package auth

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Avatar    string    `json:"avatar"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRole struct {
	UserID uint `json:"user_id" gorm:"not null"`
	RoleID uint `json:"role_id" gorm:"not null"`
}

type Role struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
}

type UserTerritory struct {
	UserID      uint `json:"user_id" gorm:"not null"`
	TerritoryID uint `json:"territory_id" gorm:"not null"`
}

type Territory struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	Type string `json:"type" gorm:"not null"`
}

type Kontingen struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	TerritoryID       uint       `json:"territory_id" gorm:"not null"`
	NamaKontingen     string     `json:"nama_kontingen" gorm:"not null"`
	Tahap1Status      string     `json:"tahap1_status"`
	Tahap1SubmittedAt *time.Time `json:"tahap1_submitted_at"`
	Tahap2Status      string     `json:"tahap2_status"`
	Tahap2SubmittedAt *time.Time `json:"tahap2_submitted_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
