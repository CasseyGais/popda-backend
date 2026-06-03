package sekolah

import "time"

// Tabel master_sekolah: id, nama, npsn, alamat, kabupaten, created_at
// Tidak ada is_active dan updated_at
type Sekolah struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	Name      string    `gorm:"size:200;column:nama" json:"nama"`
	NPSN      string    `gorm:"size:20;column:npsn" json:"npsn"`
	Alamat    string    `gorm:"type:text;column:alamat" json:"alamat"`
	Kabupaten string    `gorm:"size:150;column:kabupaten" json:"kabupaten"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (Sekolah) TableName() string {
	return "master_sekolah"
}
