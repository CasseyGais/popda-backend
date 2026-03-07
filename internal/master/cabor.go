package master

import "time"

type Cabor struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"size:150;unique;not null"`
	MaxPutra   int    `gorm:"column:max_putra"`
	MaxPutri   int    `gorm:"column:max_putri"`
	MaxPelatih int    `gorm:"column:max_pelatih"`
	IsActive   bool   `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Cabor) TableName() string {
	return "master_cabor"
}
