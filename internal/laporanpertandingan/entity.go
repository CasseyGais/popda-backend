package laporanpertandingan

import (
	"strings"
	"time"
)

// TanggalDate adalah custom type untuk kolom DATE di MariaDB.
// Serialize ke JSON sebagai "YYYY-MM-DD" (bukan timestamp dengan timezone).
type TanggalDate struct {
	time.Time
}

func (t TanggalDate) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

func (t *TanggalDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		return nil
	}
	// Support: "2026-06-12", "2026-06-12T00:00:00+07:00", "2026-06-12T00:00:00Z"
	for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01-02T15:04:05Z07:00"} {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			t.Time = parsed
			return nil
		}
	}
	return nil
}

// LaporanPertandingan memetakan tabel laporan_pertandingan.
// id: bigint(20) UNSIGNED — pakai uint
// cabor_id: int(11) — pakai uint (cocok dengan master_cabor.id)
// nomor_id: bigint(20) — pakai uint
// kontingen_a/b_id: bigint(20) — pakai uint
// created_by: int(11) — pakai *uint (nullable, FK ke users.id)
type LaporanPertandingan struct {
	ID                   uint        `gorm:"primaryKey;column:id" json:"id"`
	TanggalPertandingan  TanggalDate `gorm:"column:tanggal_pertandingan;not null" json:"tanggal_pertandingan"`
	WaktuPertandingan    string      `gorm:"column:waktu_pertandingan;not null" json:"waktu_pertandingan"`
	Venue                string     `gorm:"column:venue;not null" json:"venue"`
	CaborID              uint       `gorm:"column:cabor_id;not null" json:"cabor_id"`
	NomorID              uint       `gorm:"column:nomor_id;not null" json:"nomor_id"`
	Babak                string     `gorm:"column:babak;not null" json:"babak"`
	KontingenAID         uint       `gorm:"column:kontingen_a_id;not null" json:"kontingen_a_id"`
	KontingenBID         *uint      `gorm:"column:kontingen_b_id" json:"kontingen_b_id"`
	HasilPertandingan    string     `gorm:"column:hasil_pertandingan;not null" json:"hasil_pertandingan"`
	Pemenang             string     `gorm:"column:pemenang;not null" json:"pemenang"`
	JuaraKe              *uint8     `gorm:"column:juara_ke" json:"juara_ke"`
	Wasit                string     `gorm:"column:wasit;not null" json:"wasit"`
	CatatanKhusus        *string    `gorm:"column:catatan_khusus;type:text" json:"catatan_khusus"`
	FotoBukti            *string    `gorm:"column:foto_bukti" json:"foto_bukti"`
	VideoBukti           *string    `gorm:"column:video_bukti" json:"video_bukti"`
	CreatedBy            *uint      `gorm:"column:created_by" json:"created_by"`
	CreatedAt            time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (LaporanPertandingan) TableName() string {
	return "laporan_pertandingan"
}

// LaporanPertandinganAtlet memetakan tabel laporan_pertandingan_atlet.
// id: bigint(20) UNSIGNED — pakai uint
// laporan_id: bigint(20) UNSIGNED — pakai uint
// atlet_id: bigint(20) — pakai uint
type LaporanPertandinganAtlet struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	LaporanID uint      `gorm:"column:laporan_id;not null" json:"laporan_id"`
	Sisi      string    `gorm:"column:sisi;not null" json:"sisi"` // "A" atau "B"
	AtletID   uint      `gorm:"column:atlet_id;not null" json:"atlet_id"`
	Urutan    uint8     `gorm:"column:urutan;not null;default:1" json:"urutan"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (LaporanPertandinganAtlet) TableName() string {
	return "laporan_pertandingan_atlet"
}

// ===== RESPONSE STRUCTS =====

// LaporanDetail adalah response lengkap dengan data join untuk GET single/list
type LaporanDetail struct {
	LaporanPertandingan
	NamaCabor      string          `gorm:"column:nama_cabor" json:"nama_cabor"`
	NamaNomor      string          `gorm:"column:nama_nomor" json:"nama_nomor"`
	NamaKontingenA string          `gorm:"column:nama_kontingen_a" json:"nama_kontingen_a"`
	NamaKontingenB *string         `gorm:"column:nama_kontingen_b" json:"nama_kontingen_b"`
	AtletA         []AtletSisiItem `gorm:"-" json:"atlet_a"`
	AtletB         []AtletSisiItem `gorm:"-" json:"atlet_b"`
}

// AtletSisiItem adalah satu atlet dalam satu sisi pertandingan
// Data diambil via JOIN ke laporan_pertandingan_atlet → master_atlet + trx_pendaftaran_atlet
type AtletSisiItem struct {
	ID          uint   `gorm:"column:id" json:"id"`
	AtletID     uint   `gorm:"column:atlet_id" json:"atlet_id"`
	NamaLengkap string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	CaborID     uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NamaCabor   string `gorm:"column:nama_cabor" json:"nama_cabor"`
	NomorID     uint   `gorm:"column:nomor_id" json:"nomor_id"`
	NamaNomor   string `gorm:"column:nama_nomor" json:"nama_nomor"`
	Urutan      uint8  `gorm:"column:urutan" json:"urutan"`
}

// ===== REQUEST STRUCTS =====

// CreateLaporanRequest body POST /admin/laporan-pertandingan
type CreateLaporanRequest struct {
	TanggalPertandingan string  `json:"tanggal_pertandingan" binding:"required"` // YYYY-MM-DD
	WaktuPertandingan   string  `json:"waktu_pertandingan" binding:"required"`    // HH:MM:SS
	Venue               string  `json:"venue" binding:"required"`
	CaborID             uint    `json:"cabor_id" binding:"required"`
	NomorID             uint    `json:"nomor_id" binding:"required"`
	Babak               string  `json:"babak" binding:"required"`
	KontingenAID        uint    `json:"kontingen_a_id" binding:"required"`
	KontingenBID        *uint   `json:"kontingen_b_id"`
	HasilPertandingan   string  `json:"hasil_pertandingan" binding:"required"`
	Pemenang            string  `json:"pemenang" binding:"required"`
	JuaraKe             *uint8  `json:"juara_ke"`
	Wasit               string  `json:"wasit" binding:"required"`
	CatatanKhusus       *string `json:"catatan_khusus"`
	// AtletA dan AtletB: list atlet_id per sisi (opsional, untuk nomor individu/ganda)
	AtletA []uint `json:"atlet_a"` // ordered list of atlet_id sisi A
	AtletB []uint `json:"atlet_b"` // ordered list of atlet_id sisi B
}

// UpdateLaporanRequest body PUT /admin/laporan-pertandingan/:id — partial
type UpdateLaporanRequest struct {
	TanggalPertandingan string  `json:"tanggal_pertandingan"`
	WaktuPertandingan   string  `json:"waktu_pertandingan"`
	Venue               string  `json:"venue"`
	CaborID             uint    `json:"cabor_id"`
	NomorID             uint    `json:"nomor_id"`
	Babak               string  `json:"babak"`
	KontingenAID        uint    `json:"kontingen_a_id"`
	KontingenBID        *uint   `json:"kontingen_b_id"`
	HasilPertandingan   string  `json:"hasil_pertandingan"`
	Pemenang            string  `json:"pemenang"`
	JuaraKe             *uint8  `json:"juara_ke"`
	Wasit               string  `json:"wasit"`
	CatatanKhusus       *string `json:"catatan_khusus"`
	AtletA              []uint  `json:"atlet_a"` // jika dikirim, replace semua atlet sisi A
	AtletB              []uint  `json:"atlet_b"` // jika dikirim, replace semua atlet sisi B
}

// FilterLaporan untuk query params GET list
type FilterLaporan struct {
	CaborID   uint
	NomorID   uint
	Babak     string
	Tanggal   string
	Pemenang  string
}

// ===== DROPDOWN STRUCTS =====

// KontingenDropdownItem untuk dropdown Tim A / Tim B
type KontingenDropdownItem struct {
	ID            uint   `gorm:"column:id" json:"id"`
	NamaKontingen string `gorm:"column:nama_kontingen" json:"nama_kontingen"`
	TerritoryID   uint   `gorm:"column:territory_id" json:"territory_id"`
}

// CaborDropdownItem untuk dropdown cabang olahraga
type CaborDropdownItem struct {
	ID       uint   `gorm:"column:id" json:"id"`
	Nama     string `gorm:"column:nama" json:"nama"`
	IsActive bool   `gorm:"column:is_active" json:"is_active"`
}

// NomorDropdownItem untuk dropdown nomor/kelas per cabor
type NomorDropdownItem struct {
	ID           uint   `gorm:"column:id" json:"id"`
	CaborID      uint   `gorm:"column:cabor_id" json:"cabor_id"`
	Nama         string `gorm:"column:nama" json:"nama"`
	JenisKelamin string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Tipe         string `gorm:"column:tipe" json:"tipe"`
}

// AtletTerdaftarDropdownItem untuk dropdown atlet per kontingen + cabor + nomor
// Diambil dari trx_pendaftaran_atlet JOIN master_atlet JOIN kontingen
type AtletTerdaftarDropdownItem struct {
	AtletID       uint   `gorm:"column:atlet_id" json:"atlet_id"`
	NamaLengkap   string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	KontingenID   uint   `gorm:"column:kontingen_id" json:"kontingen_id"`
	NamaKontingen string `gorm:"column:nama_kontingen" json:"nama_kontingen"`
	CaborID       uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NomorID       uint   `gorm:"column:nomor_id" json:"nomor_id"`
}
