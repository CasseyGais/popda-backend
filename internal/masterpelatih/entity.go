package masterpelatih

import (
	"time"
)

type MasterPelatih struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	KontingenID uint      `gorm:"not null;column:kontingen_id"`
	Nama       string    `gorm:"size:150;column:nama"`
	NoHP       string    `gorm:"size:30;column:no_hp"`
	Sertifikat string    `gorm:"size:255;column:sertifikat"`
	Foto       string    `gorm:"size:255;column:foto"`
	CreatedAt  time.Time `gorm:"autoCreateTime;column:created_at"`
}

func (MasterPelatih) TableName() string {
	return "master_pelatih"
}
