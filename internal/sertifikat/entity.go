package sertifikat

import "time"

// Sertifikat memetakan tabel sertifikat.
// Tepat satu dari AtletID / PelatihID / OfficialID harus terisi — sisanya nil.
// NamaPenerima diisi otomatis dari nama_lengkap tabel master saat Create.
type Sertifikat struct {
	ID               uint       `gorm:"primaryKey;column:id" json:"id"`
	TipePenerima     string     `gorm:"column:tipe_penerima;not null" json:"tipe_penerima"`
	AtletID          *uint      `gorm:"column:atlet_id" json:"atlet_id"`
	PelatihID        *uint      `gorm:"column:pelatih_id" json:"pelatih_id"`
	OfficialID       *uint      `gorm:"column:official_id" json:"official_id"`
	NamaPenerima     string     `gorm:"column:nama_penerima;not null" json:"nama_penerima"`
	Judul            string     `gorm:"column:judul;not null" json:"judul"`
	NomorSertifikat  *string    `gorm:"column:nomor_sertifikat" json:"nomor_sertifikat"`
	TanggalTerbit    string     `gorm:"column:tanggal_terbit;not null" json:"tanggal_terbit"`
	FileSertifikat   *string    `gorm:"column:file_sertifikat" json:"file_sertifikat"`
	Catatan          *string    `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedAt        time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (Sertifikat) TableName() string {
	return "sertifikat"
}

// ===== REQUEST STRUCTS =====

// CreateSertifikatRequest body POST /admin/sertifikat.
// NamaPenerima TIDAK perlu dikirim — diisi otomatis dari nama_lengkap tabel master.
type CreateSertifikatRequest struct {
	TipePenerima    string  `json:"tipe_penerima" binding:"required,oneof=ATLET PELATIH OFFICIAL"`
	AtletID         *uint   `json:"atlet_id"`
	PelatihID       *uint   `json:"pelatih_id"`
	OfficialID      *uint   `json:"official_id"`
	Judul           string  `json:"judul" binding:"required"`
	NomorSertifikat *string `json:"nomor_sertifikat"`
	TanggalTerbit   string  `json:"tanggal_terbit" binding:"required"`
	Catatan         *string `json:"catatan"`
}

// UpdateSertifikatRequest body PUT /admin/sertifikat/:id.
// Partial update — semua field opsional.
// NamaPenerima tidak bisa diubah via PUT.
type UpdateSertifikatRequest struct {
	Judul           string  `json:"judul"`
	NomorSertifikat *string `json:"nomor_sertifikat"`
	TanggalTerbit   string  `json:"tanggal_terbit"`
	Catatan         *string `json:"catatan"`
}
