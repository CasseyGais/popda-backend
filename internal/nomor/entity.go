package nomor

import (
	"popda_bulutangkis/internal/cabor"
	"time"
)

type Nomor struct {
	ID        uint   `gorm:"primaryKey"`
	Nama      string `gorm:"size:100;not null;column:nama"`
	CaborID   uint   `gorm:"not null;column:cabor_id"`
	IsActive  bool   `gorm:"default:true;column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relations
	Cabor cabor.Cabor `gorm:"foreignKey:CaborID"`
}

func (Nomor) TableName() string {
	return "master_nomor"
}
