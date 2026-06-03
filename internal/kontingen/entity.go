package kontingen

import "time"

// TahapStatus adalah nilai valid untuk kolom tahap1_status dan tahap2_status
// sesuai ENUM('DRAFT','SUBMITTED') di database
type TahapStatus string

const (
	TahapStatusDraft     TahapStatus = "DRAFT"
	TahapStatusSubmitted TahapStatus = "SUBMITTED"
)

type Kontingen struct {
	ID              uint        `gorm:"primaryKey;column:id" json:"id"`
	TerritoryID     uint        `gorm:"column:territory_id;not null" json:"territory_id"`
	NamaKontingen   string      `gorm:"column:nama_kontingen;size:150" json:"nama_kontingen"`
	Tahap1Status    TahapStatus `gorm:"column:tahap1_status;type:enum('DRAFT','SUBMITTED');default:DRAFT" json:"tahap1_status"`
	Tahap1Submitted *time.Time  `gorm:"column:tahap1_submitted_at" json:"tahap1_submitted_at"`
	Tahap2Status    TahapStatus `gorm:"column:tahap2_status;type:enum('DRAFT','SUBMITTED');default:DRAFT" json:"tahap2_status"`
	Tahap2Submitted *time.Time  `gorm:"column:tahap2_submitted_at" json:"tahap2_submitted_at"`
	CreatedAt       time.Time   `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Kontingen) TableName() string {
	return "kontingen"
}

// IdentitasKontingen adalah alias ke tabel kontingen_identitas.
// Struct utama ada di package kontingenidentitas, ini dipakai internal kontingen package.
type IdentitasKontingen struct {
	ID          uint `gorm:"primaryKey;column:id" json:"id"`
	KontingenID uint `gorm:"column:kontingen_id;uniqueIndex;not null" json:"kontingen_id"` // UNIQUE per SQL schema

	// Data Kepala
	KepalaNama    string `gorm:"column:kepala_nama;size:150" json:"kepala_nama"`
	KepalaJabatan string `gorm:"column:kepala_jabatan;size:150" json:"kepala_jabatan"`
	KepalaNIP     string `gorm:"column:kepala_nip;size:50" json:"kepala_nip"`
	KepalaTelepon string `gorm:"column:kepala_telepon;size:30" json:"kepala_telepon"`
	KepalaFoto    string `gorm:"column:kepala_foto;size:255" json:"kepala_foto"` // Path: /uploads/kepala/

	// Data PIC (Person In Charge)
	PICNama    string `gorm:"column:pic_nama;size:150" json:"pic_nama"`
	PICJabatan string `gorm:"column:pic_jabatan;size:150" json:"pic_jabatan"`
	PICTelepon string `gorm:"column:pic_telepon;size:30" json:"pic_telepon"`
	PICFoto    string `gorm:"column:pic_foto;size:255" json:"pic_foto"` // Path: /uploads/pic/

	// Data Instansi
	Alamat        string `gorm:"column:alamat;type:text" json:"alamat"`
	EmailInstansi string `gorm:"column:email_instansi;size:150" json:"email_instansi"`
	PhoneInstansi string `gorm:"column:phone_instansi;size:30" json:"phone_instansi"`

	// Timestamps
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (IdentitasKontingen) TableName() string {
	return "kontingen_identitas"
}
