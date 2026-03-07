package master

import "time"

type Nomor struct {
	ID        uint   `gorm:"primaryKey"`
	CaborID   uint   `gorm:"column:cabor_id;index"`
	Name      string `gorm:"size:200;not null"`
	Gender    string `gorm:"size:20;column:jenis_kelamin"` // PUTRA, PUTRI, CAMPURAN
	Tipe      string `gorm:"size:20;column:tipe"`          // INDIVIDU, BEREGU
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Nomor) TableName() string {
	return "master_nomor"
}
