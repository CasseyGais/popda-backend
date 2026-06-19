package validasipendaftaran

import "time"

// KontingenValidasi adalah hasil query JOIN kontingen + territories
// untuk halaman Validasi Pendaftaran superadmin.
type KontingenValidasi struct {
	KontingenID   uint    `gorm:"column:kontingen_id" json:"kontingen_id"`
	TerritoryID   uint    `gorm:"column:territory_id" json:"territory_id"`
	NamaKontingen string  `gorm:"column:nama_kontingen" json:"nama_kontingen"`

	// Tahap 1
	Tahap1SubmitStatus    string     `gorm:"column:tahap1_status" json:"submit_status_1"`
	Tahap1SubmittedAt     *time.Time `gorm:"column:tahap1_submitted_at" json:"submitted_at_1"`
	Tahap1ValidasiStatus  *string    `gorm:"column:tahap1_validasi_status" json:"validasi_status_1"`
	Tahap1ValidasiCatatan *string    `gorm:"column:tahap1_validasi_catatan" json:"validasi_catatan_1"`
	Tahap1ValidasiAt      *time.Time `gorm:"column:tahap1_validasi_at" json:"validasi_at_1"`

	// Tahap 2
	Tahap2SubmitStatus    string     `gorm:"column:tahap2_status" json:"submit_status_2"`
	Tahap2SubmittedAt     *time.Time `gorm:"column:tahap2_submitted_at" json:"submitted_at_2"`
	Tahap2ValidasiStatus  *string    `gorm:"column:tahap2_validasi_status" json:"validasi_status_2"`
	Tahap2ValidasiCatatan *string    `gorm:"column:tahap2_validasi_catatan" json:"validasi_catatan_2"`
	Tahap2ValidasiAt      *time.Time `gorm:"column:tahap2_validasi_at" json:"validasi_at_2"`

	// Tahap 3
	Tahap3SubmitStatus    string     `gorm:"column:tahap3_status" json:"submit_status_3"`
	Tahap3SubmittedAt     *time.Time `gorm:"column:tahap3_submitted_at" json:"submitted_at_3"`
	Tahap3ValidasiStatus  *string    `gorm:"column:tahap3_validasi_status" json:"validasi_status_3"`
	Tahap3ValidasiCatatan *string    `gorm:"column:tahap3_validasi_catatan" json:"validasi_catatan_3"`
	Tahap3ValidasiAt      *time.Time `gorm:"column:tahap3_validasi_at" json:"validasi_at_3"`
}

// KontingenRow dipakai untuk cek status submit di SetValidasi.
type KontingenRow struct {
	ID            uint   `gorm:"column:id"`
	NamaKontingen string `gorm:"column:nama_kontingen"`
	Tahap1Status  string `gorm:"column:tahap1_status"`
	Tahap2Status  string `gorm:"column:tahap2_status"`
	Tahap3Status  string `gorm:"column:tahap3_status"`
}

// KontingenStatusRow dipakai untuk widget dashboard — hanya field validasi.
type KontingenStatusRow struct {
	KontingenID           uint    `gorm:"column:id"`
	NamaKontingen         string  `gorm:"column:nama_kontingen"`
	Tahap1ValidasiStatus  *string `gorm:"column:tahap1_validasi_status"`
	Tahap1ValidasiCatatan *string `gorm:"column:tahap1_validasi_catatan"`
	Tahap2ValidasiStatus  *string `gorm:"column:tahap2_validasi_status"`
	Tahap2ValidasiCatatan *string `gorm:"column:tahap2_validasi_catatan"`
	Tahap3ValidasiStatus  *string `gorm:"column:tahap3_validasi_status"`
	Tahap3ValidasiCatatan *string `gorm:"column:tahap3_validasi_catatan"`
}

// SetValidasiRequest adalah body PUT /admin/validasi-pendaftaran/:id/tahap/:tahap
type SetValidasiRequest struct {
	Status  string  `json:"status" binding:"required"` // "VALID" atau "REVISI"
	Catatan *string `json:"catatan"`
}

// ValidasiResult adalah response setelah PUT berhasil
type ValidasiResult struct {
	KontingenID   uint       `json:"kontingen_id"`
	NamaKontingen string     `json:"nama_kontingen"`
	Tahap         int        `json:"tahap"`
	Status        string     `json:"status"`
	Catatan       *string    `json:"catatan"`
	ValidasiAt    *time.Time `json:"validasi_at"`
}
