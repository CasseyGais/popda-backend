package cabor

import "time"

type Cabor struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	Nama       string    `gorm:"size:150;not null;column:nama" json:"nama"`
	MaxPutra   int       `gorm:"default:0;column:max_putra" json:"max_putra"`
	MaxPutri   int       `gorm:"default:0;column:max_putri" json:"max_putri"`
	MaxPelatih int       `gorm:"default:0;column:max_pelatih" json:"max_pelatih"`
	IsActive   bool      `gorm:"default:true;column:is_active" json:"is_active"`
	CreatedAt  time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	// Tidak ada updated_at di tabel master_cabor
}

func (Cabor) TableName() string {
	return "master_cabor"
}
