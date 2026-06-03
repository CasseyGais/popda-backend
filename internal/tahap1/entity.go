package tahap1

import "time"


type TrxKontingenCabor struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	KontingenID   uint      `gorm:"column:kontingen_id;not null" json:"kontingen_id"`
	CaborID       uint      `gorm:"column:cabor_id;not null" json:"cabor_id"`
	Putra         int       `gorm:"column:putra;default:0" json:"putra"`
	Putri         int       `gorm:"column:putri;default:0" json:"putri"`
	Pelatih       int       `gorm:"column:pelatih;default:0" json:"pelatih"`
	TotalAtlet    int       `gorm:"column:total_atlet;default:0" json:"total_atlet"`
	TotalPersonel int       `gorm:"column:total_personel;default:0" json:"total_personel"`
	CreatedAt     time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (TrxKontingenCabor) TableName() string {
	return "trx_kontingen_cabor"
}

// Kontingen dipakai untuk baca/update status tahap1 di tabel kontingen
type Kontingen struct {
	ID                uint       `gorm:"primaryKey;column:id" json:"id"`
	TerritoryID       uint       `gorm:"column:territory_id" json:"territory_id"`
	NamaKontingen     string     `gorm:"column:nama_kontingen" json:"nama_kontingen"`
	Tahap1Status      string     `gorm:"column:tahap1_status" json:"tahap1_status"`
	Tahap1SubmittedAt *time.Time `gorm:"column:tahap1_submitted_at" json:"tahap1_submitted_at"`
	Tahap2Status      string     `gorm:"column:tahap2_status" json:"tahap2_status"`
	Tahap2SubmittedAt *time.Time `gorm:"column:tahap2_submitted_at" json:"tahap2_submitted_at"`
	Tahap3Status      string     `gorm:"column:tahap3_status" json:"tahap3_status"`
	Tahap3SubmittedAt *time.Time `gorm:"column:tahap3_submitted_at" json:"tahap3_submitted_at"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Kontingen) TableName() string {
	return "kontingen"
}
