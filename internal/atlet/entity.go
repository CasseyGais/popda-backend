package atlet

import (
	"time"
)

type Atlet struct {
	ID              uint      `gorm:"primaryKey;column:id"`
	KontingenID     uint      `gorm:"not null;column:kontingen_id;index"`
	SekolahID       uint      `gorm:"not null;column:sekolah_id;index"`
	NISN            string    `gorm:"size:50;column:nisn"`
	Nama            string    `gorm:"size:150;column:nama"`
	JenisKelamin    string    `gorm:"type:enum('PUTRA','PUTRI');column:jenis_kelamin"`
	TanggalLahir    *time.Time `gorm:"type:date;column:tanggal_lahir"`
	Kelas           string    `gorm:"size:20;column:kelas"`
	Tinggi          *int      `gorm:"column:tinggi"`
	Berat           *float64  `gorm:"type:decimal(5,2);column:berat"`
	Foto            string    `gorm:"size:255;column:foto"`
	StatusVerifikasi string    `gorm:"type:enum('PENDING','VALID','DITOLAK');default:'PENDING';column:status_verifikasi"`
	CreatedAt       time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (Atlet) TableName() string {
	return "atlet"
}
