package masterpelatih

import "time"

// MasterPelatih memetakan tabel master_pelatih (struktur lengkap)
type MasterPelatih struct {
	ID                    uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID           uint      `gorm:"column:kontingen_id;not null" json:"kontingen_id"`
	NamaLengkap           string    `gorm:"column:nama_lengkap;not null" json:"nama_lengkap"`
	JenisKelamin          string    `gorm:"column:jenis_kelamin;not null" json:"jenis_kelamin"`
	TanggalLahir          string    `gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	TempatLahir           string    `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	NIK                   string    `gorm:"column:nik" json:"nik"`
	SekolahAsal           string    `gorm:"column:sekolah_asal" json:"sekolah_asal"`
	Profesi               string    `gorm:"column:profesi" json:"profesi"`
	Jabatan               string    `gorm:"column:jabatan" json:"jabatan"`
	Alamat                string    `gorm:"column:alamat;type:text" json:"alamat"`
	KabupatenKota         string    `gorm:"column:kabupaten_kota;not null" json:"kabupaten_kota"`
	NoHP                  string    `gorm:"column:no_hp;not null" json:"no_hp"`
	Email                 string    `gorm:"column:email" json:"email"`
	NamaIstriSuami        string    `gorm:"column:nama_istri_suami" json:"nama_istri_suami"`
	Status                string    `gorm:"column:status;default:'draft'" json:"status"`
	Foto                  string    `gorm:"column:foto" json:"foto"`
	FileKTP               string    `gorm:"column:file_ktp" json:"file_ktp"`
	FileSuratTugas        string    `gorm:"column:file_surat_tugas" json:"file_surat_tugas"`
	FileSertifikatPelatih string    `gorm:"column:file_sertifikat_pelatih" json:"file_sertifikat_pelatih"`
	PrestasiSebelumnya    string    `gorm:"column:prestasi_sebelumnya;type:text" json:"prestasi_sebelumnya"`
	Catatan               string    `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedAt             time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (MasterPelatih) TableName() string {
	return "master_pelatih"
}

// TrxPendaftaranPelatih memetakan tabel trx_pendaftaran_pelatih
type TrxPendaftaranPelatih struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	PelatihID uint      `gorm:"column:pelatih_id;not null" json:"pelatih_id"`
	CaborID   uint      `gorm:"column:cabor_id;not null" json:"cabor_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (TrxPendaftaranPelatih) TableName() string {
	return "trx_pendaftaran_pelatih"
}

// Kontingen untuk resolveKontingenID
type Kontingen struct {
	ID          uint   `gorm:"primaryKey;column:id" json:"id"`
	TerritoryID uint   `gorm:"column:territory_id" json:"territory_id"`
}

func (Kontingen) TableName() string {
	return "kontingen"
}

// ===== REQUEST STRUCTS =====

type CreateMasterPelatihRequest struct {
	NamaLengkap           string `json:"nama_lengkap" binding:"required"`
	JenisKelamin          string `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TanggalLahir          string `json:"tanggal_lahir"`
	TempatLahir           string `json:"tempat_lahir"`
	NIK                   string `json:"nik"`
	SekolahAsal           string `json:"sekolah_asal"`
	Profesi               string `json:"profesi"`
	Jabatan               string `json:"jabatan"`
	Alamat                string `json:"alamat"`
	KabupatenKota         string `json:"kabupaten_kota" binding:"required"`
	NoHP                  string `json:"no_hp" binding:"required"`
	Email                 string `json:"email"`
	NamaIstriSuami        string `json:"nama_istri_suami"`
	PrestasiSebelumnya    string `json:"prestasi_sebelumnya"`
	Catatan               string `json:"catatan"`
}

type UpdateMasterPelatihRequest struct {
	NamaLengkap           string `json:"nama_lengkap"`
	JenisKelamin          string `json:"jenis_kelamin" binding:"omitempty,oneof=L P"`
	TanggalLahir          string `json:"tanggal_lahir"`
	TempatLahir           string `json:"tempat_lahir"`
	NIK                   string `json:"nik"`
	SekolahAsal           string `json:"sekolah_asal"`
	Profesi               string `json:"profesi"`
	Jabatan               string `json:"jabatan"`
	Alamat                string `json:"alamat"`
	KabupatenKota         string `json:"kabupaten_kota"`
	NoHP                  string `json:"no_hp"`
	Email                 string `json:"email"`
	NamaIstriSuami        string `json:"nama_istri_suami"`
	PrestasiSebelumnya    string `json:"prestasi_sebelumnya"`
	Catatan               string `json:"catatan"`
}
