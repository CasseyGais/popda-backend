package masterofficial

import "time"

// MasterOfficial memetakan tabel master_official (struktur lengkap)
type MasterOfficial struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID   uint      `gorm:"column:kontingen_id;not null" json:"kontingen_id"`
	NamaLengkap   string    `gorm:"column:nama_lengkap;not null" json:"nama_lengkap"`
	JenisKelamin  string    `gorm:"column:jenis_kelamin;not null" json:"jenis_kelamin"`
	TanggalLahir  string    `gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	TempatLahir   string    `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	NIK           string    `gorm:"column:nik" json:"nik"`
	SekolahAsal   string    `gorm:"column:sekolah_asal" json:"sekolah_asal"`
	Jabatan       string    `gorm:"column:jabatan;not null" json:"jabatan"`
	Alamat        string    `gorm:"column:alamat;type:text" json:"alamat"`
	KabupatenKota string    `gorm:"column:kabupaten_kota;not null" json:"kabupaten_kota"`
	NoHP          string    `gorm:"column:no_hp;not null" json:"no_hp"`
	Email         string    `gorm:"column:email" json:"email"`
	Status        string    `gorm:"column:status;default:'draft'" json:"status"`
	Foto          string    `gorm:"column:foto" json:"foto"`
	FileKTP       string    `gorm:"column:file_ktp" json:"file_ktp"`
	FileSuratTugas string   `gorm:"column:file_surat_tugas" json:"file_surat_tugas"`
	Catatan       string    `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedAt     time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (MasterOfficial) TableName() string {
	return "master_official"
}

// TrxPendaftaranOfficial memetakan tabel trx_pendaftaran_official
type TrxPendaftaranOfficial struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	OfficialID uint      `gorm:"column:official_id;not null" json:"official_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (TrxPendaftaranOfficial) TableName() string {
	return "trx_pendaftaran_official"
}

// ===== REQUEST STRUCTS =====

type CreateMasterOfficialRequest struct {
	NamaLengkap   string `json:"nama_lengkap" binding:"required"`
	JenisKelamin  string `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TanggalLahir  string `json:"tanggal_lahir"`
	TempatLahir   string `json:"tempat_lahir"`
	NIK           string `json:"nik"`
	SekolahAsal   string `json:"sekolah_asal"`
	Jabatan       string `json:"jabatan" binding:"required"`
	Alamat        string `json:"alamat"`
	KabupatenKota string `json:"kabupaten_kota" binding:"required"`
	NoHP          string `json:"no_hp" binding:"required"`
	Email         string `json:"email"`
	Catatan       string `json:"catatan"`
}

type UpdateMasterOfficialRequest struct {
	NamaLengkap   string `json:"nama_lengkap"`
	JenisKelamin  string `json:"jenis_kelamin" binding:"omitempty,oneof=L P"`
	TanggalLahir  string `json:"tanggal_lahir"`
	TempatLahir   string `json:"tempat_lahir"`
	NIK           string `json:"nik"`
	SekolahAsal   string `json:"sekolah_asal"`
	Jabatan       string `json:"jabatan"`
	Alamat        string `json:"alamat"`
	KabupatenKota string `json:"kabupaten_kota"`
	NoHP          string `json:"no_hp"`
	Email         string `json:"email"`
	Catatan       string `json:"catatan"`
}
