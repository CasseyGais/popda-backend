package transaksi

import "time"

type TrxKontingenCabor struct {
	ID           uint      `gorm:"primaryKey"`
	KontingenID  uint      `gorm:"column:kontingen_id"`
	CaborID      uint      `gorm:"column:cabor_id"`
	Putra        int       `gorm:"column:putra"`
	Putri        int       `gorm:"column:putri"`
	Pelatih      int       `gorm:"column:pelatih"`
	TotalAtlet   int       `gorm:"column:total_atlet"`
	TotalPersonel int      `gorm:"column:total_personel"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (TrxKontingenCabor) TableName() string {
	return "trx_kontingen_cabor"
}

type TrxKontingenNomor struct {
	ID          uint      `gorm:"primaryKey"`
	KontingenID uint      `gorm:"column:kontingen_id"`
	NomorID     uint      `gorm:"column:nomor_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (TrxKontingenNomor) TableName() string {
	return "trx_kontingen_nomor"
}

type TrxPendaftaranAtlet struct {
	ID        uint      `gorm:"primaryKey"`
	AtletID   uint      `gorm:"column:atlet_id"`
	NomorID   uint      `gorm:"column:nomor_id"`
	KelasID   *uint     `gorm:"column:kelas_id"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (TrxPendaftaranAtlet) TableName() string {
	return "trx_pendaftaran_atlet"
}
