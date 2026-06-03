package masterofficial

import "time"

// Tabel master_official: id, kontingen_id, nama, jabatan, no_hp, created_at
// Tidak ada updated_at
type MasterOfficial struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID uint      `gorm:"not null;column:kontingen_id" json:"kontingen_id"`
	Nama        string    `gorm:"size:150;column:nama" json:"nama"`
	Jabatan     string    `gorm:"size:100;column:jabatan" json:"jabatan"`
	NoHP        string    `gorm:"size:30;column:no_hp" json:"no_hp"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (MasterOfficial) TableName() string {
	return "master_official"
}
