package cabor

import (
	"time"
)

type Cabor struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	Nama       string    `gorm:"size:150;not null;column:nama"`
	MaxPutra   int       `gorm:"default:0;column:max_putra"`
	MaxPutri   int       `gorm:"default:0;column:max_putri"`
	MaxPelatih int       `gorm:"default:0;column:max_pelatih"`
	IsActive   bool      `gorm:"default:true;column:is_active"`
	CreatedAt  time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (Cabor) TableName() string {
	return "master_cabor"
}
