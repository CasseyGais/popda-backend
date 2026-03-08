package kontingenidentitas

import (
	"time"
)

type KontingenIdentitas struct {
	ID            uint      `gorm:"primaryKey;column:id"`
	KontingenID   uint      `gorm:"not null;column:kontingen_id"`
	KepalaNama    string    `gorm:"size:150;column:kepala_nama"`
	KepalaJabatan string    `gorm:"size:150;column:kepala_jabatan"`
	KepalaNIP     string    `gorm:"size:50;column:kepala_nip"`
	KepalaTelepon string    `gorm:"size:30;column:kepala_telepon"`
	KepalaFoto    string    `gorm:"size:255;column:kepala_foto"`
	PICNama       string    `gorm:"size:150;column:pic_nama"`
	PICJabatan    string    `gorm:"size:150;column:pic_jabatan"`
	PICTelepon    string    `gorm:"size:30;column:pic_telepon"`
	PICFoto       string    `gorm:"size:255;column:pic_foto"`
	Alamat        string    `gorm:"type:text;column:alamat"`
	EmailInstansi  string    `gorm:"size:150;column:email_instansi"`
	PhoneInstansi string    `gorm:"size:30;column:phone_instansi"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (KontingenIdentitas) TableName() string {
	return "kontingen_identitas"
}
