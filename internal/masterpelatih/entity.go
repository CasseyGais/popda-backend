package masterpelatih

import "time"

// Tabel master_pelatih: id, kontingen_id, nama, no_hp, sertifikat, foto, created_at
// Tidak ada updated_at
type MasterPelatih struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID uint      `gorm:"not null;column:kontingen_id" json:"kontingen_id"`
	Nama        string    `gorm:"size:150;column:nama" json:"nama"`
	NoHP        string    `gorm:"size:30;column:no_hp" json:"no_hp"`
	Sertifikat  string    `gorm:"size:255;column:sertifikat" json:"sertifikat"`
	Foto        string    `gorm:"size:255;column:foto" json:"foto"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (MasterPelatih) TableName() string {
	return "master_pelatih"
}
