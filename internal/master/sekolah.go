package master

import "time"

type Sekolah struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:200;not null;column:nama"`
	NPSN      string `gorm:"size:50;column:npsn"`
	Alamat    string `gorm:"type:text;column:alamat"`
	Kabupaten string `gorm:"size:150;column:kabupaten"`
	CreatedAt time.Time
}

func (Sekolah) TableName() string {
	return "master_sekolah"
}
