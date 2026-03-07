package tahap1

import "time"

type CaborKontingen struct {
	ID          uint `gorm:"primaryKey"`
	KontingenID uint `gorm:"column:kontingen_id;index"`
	CaborID     uint `gorm:"column:cabor_id"`
	Putra       int  `gorm:"column:putra"`
	Putri       int  `gorm:"column:putri"`
	Pelatih     int  `gorm:"column:pelatih"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (CaborKontingen) TableName() string {
	return "trx_kontingen_cabor"
}

type Tahap1Document struct {
	ID              uint   `gorm:"primaryKey"`
	KontingenID     uint   `gorm:"column:kontingen_id;index"`
	SuratPernyataan string `gorm:"size:255"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Tahap1Document) TableName() string {
	return "tahap1_documents"
}

type KontingenMeta struct {
	ID              uint `gorm:"primaryKey"`
	KontingenID     uint `gorm:"column:kontingen_id;index"`
	JumlahAtlet     int
	JumlahPelatih   int
	JumlahOfficial  int
	Tahap1Submitted bool
	SubmittedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (KontingenMeta) TableName() string {
	return "kontingen_metas"
}
