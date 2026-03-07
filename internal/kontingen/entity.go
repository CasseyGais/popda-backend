package kontingen

import "time"

type Kontingen struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TerritoryID     uint       `gorm:"column:territory_id" json:"territory_id"`
	NamaKontingen   string     `gorm:"column:nama_kontingen" json:"nama_kontingen"`
	Tahap1Status    string     `gorm:"column:tahap1_status" json:"tahap1_status"`
	Tahap1Submitted *time.Time `gorm:"column:tahap1_submitted_at" json:"tahap1_submitted_at"`
	Tahap2Status    string     `gorm:"column:tahap2_status" json:"tahap2_status"`
	Tahap2Submitted *time.Time `gorm:"column:tahap2_submitted_at" json:"tahap2_submitted_at"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Kontingen) TableName() string {
	return "kontingen"
}

type IdentitasKontingen struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	KontingenID uint `gorm:"column:kontingen_id;index" json:"kontingen_id"`

	// Data Kepala
	KepalaNama    string `gorm:"column:kepala_nama" json:"kepala_nama"`       // Nama lengkap kepala
	KepalaJabatan string `gorm:"column:kepala_jabatan" json:"kepala_jabatan"` // Jabatan kepala
	KepalaNIP     string `gorm:"column:kepala_nip" json:"kepala_nip"`         // NIP kepala
	KepalaTelepon string `gorm:"column:kepala_telepon" json:"kepala_telepon"` // Nomor telepon kepala
	KepalaFoto    string `gorm:"column:kepala_foto" json:"kepala_foto"`       // Path foto kepala  (uploads/foto/kepala/)

	// Data PIC (Person In Charge)
	PICNama    string `gorm:"column:pic_nama" json:"pic_nama"`       // Nama lengkap PIC
	PICJabatan string `gorm:"column:pic_jabatan" json:"pic_jabatan"` // Jabatan PIC
	PICTelepon string `gorm:"column:pic_telepon" json:"pic_telepon"` // Nomor telepon PIC
	PICFoto    string `gorm:"column:pic_foto" json:"pic_foto"`       // Path foto PIC (uploads/foto/pic/)

	// Data instansi
	Alamat        string `gorm:"type:text" json:"alamat"`                     // Alamat instansi
	EmailInstansi string `gorm:"column:email_instansi" json:"email_instansi"` // Email instansi
	PhoneInstansi string `gorm:"column:phone_instansi" json:"phone_instansi"` // Telepon instansi

	// Status dan timestamps
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"` // Waktu diupdate

	// Foreign key relation
	Kontingen Kontingen `gorm:"foreignKey:KontingenID"`
}

func (IdentitasKontingen) TableName() string {
	return "kontingen_identitas"
}
