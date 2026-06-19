package tahap2

import "time"

// TrxKontingenNomor memetakan tabel trx_kontingen_nomor.
// Kolom: id, kontingen_id, nomor_id, created_at
// Tidak ada updated_at di tabel ini.
type TrxKontingenNomor struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID uint      `gorm:"column:kontingen_id;not null" json:"kontingen_id"`
	NomorID     uint      `gorm:"column:nomor_id;not null" json:"nomor_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (TrxKontingenNomor) TableName() string {
	return "trx_kontingen_nomor"
}

// NomorWithStatus adalah BUKAN tabel — ini scan result dari JOIN query.
// Menggabungkan data dari:
//   - master_nomor     (id, cabor_id, nama, jenis_kelamin, tipe)
//   - master_cabor     (nama sebagai nama_cabor)
//   - trx_kontingen_nomor (untuk cek apakah kontingen sudah daftar)
//
// Field harus cocok dengan alias di SELECT query di repository.
type NomorWithStatus struct {
	NomorID      uint   `gorm:"column:nomor_id" json:"nomor_id"`
	CaborID      uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NamaCabor    string `gorm:"column:nama_cabor" json:"nama_cabor"`
	NamaNomor    string `gorm:"column:nama_nomor" json:"nama_nomor"`
	JenisKelamin string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Tipe         string `gorm:"column:tipe" json:"tipe"`
	Terdaftar    bool   `gorm:"column:terdaftar" json:"terdaftar"`
}

// Kontingen untuk baca/update tahap2_status dari tabel kontingen.
type Kontingen struct {
	ID                    uint       `gorm:"primaryKey;column:id" json:"id"`
	TerritoryID           uint       `gorm:"column:territory_id" json:"territory_id"`
	NamaKontingen         string     `gorm:"column:nama_kontingen" json:"nama_kontingen"`
	Tahap1Status          string     `gorm:"column:tahap1_status" json:"tahap1_status"`
	Tahap1SubmittedAt     *time.Time `gorm:"column:tahap1_submitted_at" json:"tahap1_submitted_at"`
	Tahap1ValidasiStatus  *string    `gorm:"column:tahap1_validasi_status" json:"tahap1_validasi_status"`
	Tahap1ValidasiCatatan *string    `gorm:"column:tahap1_validasi_catatan" json:"tahap1_validasi_catatan"`
	Tahap1ValidasiAt      *time.Time `gorm:"column:tahap1_validasi_at" json:"tahap1_validasi_at"`
	Tahap2Status          string     `gorm:"column:tahap2_status" json:"tahap2_status"`
	Tahap2SubmittedAt     *time.Time `gorm:"column:tahap2_submitted_at" json:"tahap2_submitted_at"`
	Tahap2ValidasiStatus  *string    `gorm:"column:tahap2_validasi_status" json:"tahap2_validasi_status"`
	Tahap2ValidasiCatatan *string    `gorm:"column:tahap2_validasi_catatan" json:"tahap2_validasi_catatan"`
	Tahap2ValidasiAt      *time.Time `gorm:"column:tahap2_validasi_at" json:"tahap2_validasi_at"`
	Tahap3Status          string     `gorm:"column:tahap3_status" json:"tahap3_status"`
	Tahap3SubmittedAt     *time.Time `gorm:"column:tahap3_submitted_at" json:"tahap3_submitted_at"`
	Tahap3ValidasiStatus  *string    `gorm:"column:tahap3_validasi_status" json:"tahap3_validasi_status"`
	Tahap3ValidasiCatatan *string    `gorm:"column:tahap3_validasi_catatan" json:"tahap3_validasi_catatan"`
	Tahap3ValidasiAt      *time.Time `gorm:"column:tahap3_validasi_at" json:"tahap3_validasi_at"`
	CreatedAt             time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Kontingen) TableName() string {
	return "kontingen"
}
