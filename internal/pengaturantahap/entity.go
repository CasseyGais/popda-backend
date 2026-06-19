package pengaturantahap

import "time"

// PengaturanTahap merepresentasikan satu row di tabel pengaturan_tahap.
// is_open adalah kontrol utama — tanggal_buka/tutup bersifat informatif.
type PengaturanTahap struct {
	ID           uint       `gorm:"primaryKey;column:id" json:"id"`
	Tahap        uint       `gorm:"column:tahap;uniqueIndex" json:"tahap"`
	IsOpen       bool       `gorm:"column:is_open;default:0" json:"is_open"`
	TanggalBuka  *string    `gorm:"column:tanggal_buka" json:"tanggal_buka"`
	TanggalTutup *string    `gorm:"column:tanggal_tutup" json:"tanggal_tutup"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (PengaturanTahap) TableName() string {
	return "pengaturan_tahap"
}

// UpdatePengaturanRequest adalah body request PUT /admin/pengaturan-tahap/:tahap.
// Semua field opsional — hanya yang dikirim yang diupdate.
type UpdatePengaturanRequest struct {
	IsOpen       *bool   `json:"is_open"`
	TanggalBuka  *string `json:"tanggal_buka"`
	TanggalTutup *string `json:"tanggal_tutup"`
}
