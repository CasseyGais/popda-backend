package tahap3

import "time"

// ===== MASTER TABLES =====

// MasterAtlet memetakan tabel master_atlet di popda_2026.
// jenis_kelamin: enum('L','P')
// status: enum('draft','terdaftar','terverifikasi','ditolak')
type MasterAtlet struct {
	ID                       uint       `gorm:"primaryKey;column:id" json:"id"`
	KontingenID              uint       `gorm:"column:kontingen_id;not null" json:"kontingen_id"`
	NamaLengkap              string     `gorm:"column:nama_lengkap;not null" json:"nama_lengkap"`
	JenisKelamin             string     `gorm:"column:jenis_kelamin;not null" json:"jenis_kelamin"`
	TanggalLahir             string     `gorm:"column:tanggal_lahir;not null" json:"tanggal_lahir"`
	TempatLahir              string     `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	NISN                     string     `gorm:"column:nisn;not null;uniqueIndex" json:"nisn"`
	NIS                      string     `gorm:"column:nis" json:"nis"`
	Sekolah                  string     `gorm:"column:sekolah;not null" json:"sekolah"`
	KelasJurusan             string     `gorm:"column:kelas_jurusan" json:"kelas_jurusan"`
	Alamat                   string     `gorm:"column:alamat;type:text" json:"alamat"`
	KabupatenKota            string     `gorm:"column:kabupaten_kota;not null" json:"kabupaten_kota"`
	NoHP                     string     `gorm:"column:no_hp" json:"no_hp"`
	NamaOrtuWali             string     `gorm:"column:nama_ortu_wali" json:"nama_ortu_wali"`
	Status                   string     `gorm:"column:status;default:'draft'" json:"status"`
	Foto                     string     `gorm:"column:foto" json:"foto"`
	FileKartuPelajar         string     `gorm:"column:file_kartu_pelajar" json:"file_kartu_pelajar"`
	FileAkteKelahiran        string     `gorm:"column:file_akte_kelahiran" json:"file_akte_kelahiran"`
	FileKK                   string     `gorm:"column:file_kk" json:"file_kk"`
	FileSuratKeteranganSekolah string   `gorm:"column:file_surat_keterangan_sekolah" json:"file_surat_keterangan_sekolah"`
	FileSuratIzinOrtu        string     `gorm:"column:file_surat_izin_ortu" json:"file_surat_izin_ortu"`
	PrestasiSebelumnya       string     `gorm:"column:prestasi_sebelumnya;type:text" json:"prestasi_sebelumnya"`
	Catatan                  string     `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedAt                time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt                time.Time  `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (MasterAtlet) TableName() string { return "master_atlet" }

// MasterPelatih memetakan tabel master_pelatih di popda_2026.
// jenis_kelamin: enum('L','P')
// status: enum('draft','terdaftar','terverifikasi','ditolak')
type MasterPelatih struct {
	ID                    uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID           uint      `gorm:"column:kontingen_id;not null" json:"kontingen_id"`
	NamaLengkap           string    `gorm:"column:nama_lengkap;not null" json:"nama_lengkap"`
	JenisKelamin          string    `gorm:"column:jenis_kelamin;not null" json:"jenis_kelamin"`
	TanggalLahir          string    `gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	TempatLahir           string    `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	NIK                   string    `gorm:"column:nik;uniqueIndex" json:"nik"`
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

func (MasterPelatih) TableName() string { return "master_pelatih" }

// MasterOfficial memetakan tabel master_official di popda_2026.
// jenis_kelamin: enum('L','P')
// status: enum('draft','terdaftar','terverifikasi','ditolak')
type MasterOfficial struct {
	ID              uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID     uint      `gorm:"column:kontingen_id;not null" json:"kontingen_id"`
	NamaLengkap     string    `gorm:"column:nama_lengkap;not null" json:"nama_lengkap"`
	JenisKelamin    string    `gorm:"column:jenis_kelamin;not null" json:"jenis_kelamin"`
	TanggalLahir    string    `gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	TempatLahir     string    `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	NIK             string    `gorm:"column:nik;uniqueIndex" json:"nik"`
	SekolahAsal     string    `gorm:"column:sekolah_asal" json:"sekolah_asal"`
	Jabatan         string    `gorm:"column:jabatan;not null" json:"jabatan"`
	Alamat          string    `gorm:"column:alamat;type:text" json:"alamat"`
	KabupatenKota   string    `gorm:"column:kabupaten_kota;not null" json:"kabupaten_kota"`
	NoHP            string    `gorm:"column:no_hp;not null" json:"no_hp"`
	Email           string    `gorm:"column:email" json:"email"`
	Status          string    `gorm:"column:status;default:'draft'" json:"status"`
	Foto            string    `gorm:"column:foto" json:"foto"`
	FileKTP         string    `gorm:"column:file_ktp" json:"file_ktp"`
	FileSuratTugas  string    `gorm:"column:file_surat_tugas" json:"file_surat_tugas"`
	Catatan         string    `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedAt       time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (MasterOfficial) TableName() string { return "master_official" }

// ===== TRANSAKSI TABLES =====

// TrxPendaftaranAtlet memetakan tabel trx_pendaftaran_atlet.
// Kolom: id, atlet_id, cabor_id, nomor_id, created_at, updated_at
// kontingen_id diakses via JOIN ke master_atlet
type TrxPendaftaranAtlet struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	AtletID   uint      `gorm:"column:atlet_id;not null" json:"atlet_id"`
	CaborID   uint      `gorm:"column:cabor_id;not null" json:"cabor_id"`
	NomorID   uint      `gorm:"column:nomor_id;not null" json:"nomor_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (TrxPendaftaranAtlet) TableName() string { return "trx_pendaftaran_atlet" }

// TrxPendaftaranPelatih memetakan tabel trx_pendaftaran_pelatih.
// Kolom: id, pelatih_id, cabor_id, created_at, updated_at
type TrxPendaftaranPelatih struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	PelatihID uint      `gorm:"column:pelatih_id;not null" json:"pelatih_id"`
	CaborID   uint      `gorm:"column:cabor_id;not null" json:"cabor_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (TrxPendaftaranPelatih) TableName() string { return "trx_pendaftaran_pelatih" }

// TrxPendaftaranOfficial memetakan tabel trx_pendaftaran_official.
// Kolom: id, official_id, created_at, updated_at
type TrxPendaftaranOfficial struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	OfficialID uint      `gorm:"column:official_id;not null" json:"official_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (TrxPendaftaranOfficial) TableName() string { return "trx_pendaftaran_official" }

// ===== STATUS =====

// Kontingen untuk baca/update status tahap di tabel kontingen.
// Kolom: id, territory_id, nama_kontingen, tahap1_status, tahap1_submitted_at,
//        tahap2_status, tahap2_submitted_at, tahap3_status, tahap3_submitted_at,
//        created_at, updated_at
type Kontingen struct {
	ID                uint       `gorm:"primaryKey;column:id" json:"id"`
	TerritoryID       uint       `gorm:"column:territory_id" json:"territory_id"`
	NamaKontingen     string     `gorm:"column:nama_kontingen" json:"nama_kontingen"`
	Tahap1Status      string     `gorm:"column:tahap1_status" json:"tahap1_status"`
	Tahap1SubmittedAt *time.Time `gorm:"column:tahap1_submitted_at" json:"tahap1_submitted_at"`
	Tahap2Status      string     `gorm:"column:tahap2_status" json:"tahap2_status"`
	Tahap2SubmittedAt *time.Time `gorm:"column:tahap2_submitted_at" json:"tahap2_submitted_at"`
	Tahap3Status      string     `gorm:"column:tahap3_status" json:"tahap3_status"`
	Tahap3SubmittedAt *time.Time `gorm:"column:tahap3_submitted_at" json:"tahap3_submitted_at"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Kontingen) TableName() string { return "kontingen" }

// ===== REQUEST STRUCTS =====

type CreateAtletRequest struct {
	NamaLengkap              string `json:"nama_lengkap" binding:"required"`
	JenisKelamin             string `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TanggalLahir             string `json:"tanggal_lahir" binding:"required"`
	TempatLahir              string `json:"tempat_lahir"`
	NISN                     string `json:"nisn" binding:"required"`
	NIS                      string `json:"nis"`
	Sekolah                  string `json:"sekolah" binding:"required"`
	KelasJurusan             string `json:"kelas_jurusan"`
	Alamat                   string `json:"alamat"`
	KabupatenKota            string `json:"kabupaten_kota" binding:"required"`
	NoHP                     string `json:"no_hp"`
	NamaOrtuWali             string `json:"nama_ortu_wali"`
	PrestasiSebelumnya       string `json:"prestasi_sebelumnya"`
	Catatan                  string `json:"catatan"`
}

type UpdateAtletRequest struct {
	NamaLengkap              string `json:"nama_lengkap"`
	JenisKelamin             string `json:"jenis_kelamin" binding:"omitempty,oneof=L P"`
	TanggalLahir             string `json:"tanggal_lahir"`
	TempatLahir              string `json:"tempat_lahir"`
	NISN                     string `json:"nisn"`
	NIS                      string `json:"nis"`
	Sekolah                  string `json:"sekolah"`
	KelasJurusan             string `json:"kelas_jurusan"`
	Alamat                   string `json:"alamat"`
	KabupatenKota            string `json:"kabupaten_kota"`
	NoHP                     string `json:"no_hp"`
	NamaOrtuWali             string `json:"nama_ortu_wali"`
	PrestasiSebelumnya       string `json:"prestasi_sebelumnya"`
	Catatan                  string `json:"catatan"`
}

type CreatePelatihRequest struct {
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

type UpdatePelatihRequest struct {
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

type CreateOfficialRequest struct {
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

type UpdateOfficialRequest struct {
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

// CreateTrxAtletRequest — daftarkan atlet ke nomor tertentu
type CreateTrxAtletRequest struct {
	AtletID uint `json:"atlet_id" binding:"required"`
	CaborID uint `json:"cabor_id" binding:"required"`
	NomorID uint `json:"nomor_id" binding:"required"`
}

// CreateTrxPelatihRequest — daftarkan pelatih ke cabor tertentu
type CreateTrxPelatihRequest struct {
	PelatihID uint `json:"pelatih_id" binding:"required"`
	CaborID   uint `json:"cabor_id" binding:"required"`
}

// CreateTrxOfficialRequest — daftarkan official
type CreateTrxOfficialRequest struct {
	OfficialID uint `json:"official_id" binding:"required"`
}
