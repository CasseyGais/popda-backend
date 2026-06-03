package kontingenidentitas

import (
	"time"
)

// KontingenIdentitas memetakan tabel kontingen_identitas di database.
// kontingen_id bersifat UNIQUE (satu kontingen hanya punya satu identitas).
type KontingenIdentitas struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID   uint      `gorm:"not null;uniqueIndex;column:kontingen_id" json:"kontingen_id"` // UNIQUE per SQL schema
	KepalaNama    string    `gorm:"size:150;column:kepala_nama" json:"kepala_nama"`
	KepalaJabatan string    `gorm:"size:150;column:kepala_jabatan" json:"kepala_jabatan"`
	KepalaNIP     string    `gorm:"size:50;column:kepala_nip" json:"kepala_nip"`
	KepalaTelepon string    `gorm:"size:30;column:kepala_telepon" json:"kepala_telepon"`
	KepalaFoto    string    `gorm:"size:255;column:kepala_foto" json:"kepala_foto"`
	PICNama       string    `gorm:"size:150;column:pic_nama" json:"pic_nama"`
	PICJabatan    string    `gorm:"size:150;column:pic_jabatan" json:"pic_jabatan"`
	PICTelepon    string    `gorm:"size:30;column:pic_telepon" json:"pic_telepon"`
	PICFoto       string    `gorm:"size:255;column:pic_foto" json:"pic_foto"`
	Alamat        string    `gorm:"type:text;column:alamat" json:"alamat"`
	EmailInstansi string    `gorm:"size:150;column:email_instansi" json:"email_instansi"`
	PhoneInstansi string    `gorm:"size:30;column:phone_instansi" json:"phone_instansi"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (KontingenIdentitas) TableName() string {
	return "kontingen_identitas"
}
