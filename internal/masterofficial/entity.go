package masterofficial

import (
	"time"
)

type MasterOfficial struct {
	ID          uint      `gorm:"primaryKey;column:id"`
	KontingenID uint      `gorm:"not null;column:kontingen_id"`
	Nama        string    `gorm:"size:150;column:nama"`
	Jabatan     string    `gorm:"size:100;column:jabatan"`
	NoHP        string    `gorm:"size:30;column:no_hp"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at"`
}

func (MasterOfficial) TableName() string {
	return "master_official"
}
