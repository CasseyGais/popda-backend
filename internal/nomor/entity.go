package nomor

import (
	"popda_bulutangkis/internal/cabor"
	"time"
)

// Tabel master_nomor: id, cabor_id, nama, jenis_kelamin, tipe, is_active, created_at
// Tidak ada updated_at
type Nomor struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	CaborID      uint      `gorm:"not null;column:cabor_id" json:"cabor_id"`
	Nama         string    `gorm:"size:255;not null;column:nama" json:"nama"`
	JenisKelamin string    `gorm:"type:enum('PUTRA','PUTRI','CAMPURAN');not null;column:jenis_kelamin" json:"jenis_kelamin"`
	Tipe         string    `gorm:"type:enum('INDIVIDU','BEREGU');default:'INDIVIDU';column:tipe" json:"tipe"`
	IsActive     bool      `gorm:"default:true;column:is_active" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	// Relasi
	Cabor cabor.Cabor `gorm:"foreignKey:CaborID" json:"cabor"`
}

func (Nomor) TableName() string {
	return "master_nomor"
}
