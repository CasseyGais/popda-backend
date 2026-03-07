package master

type Kelas struct {
	ID       uint `gorm:"primaryKey"`
	NomorID  uint `gorm:"column:nomor_id"`
	Nama     string `gorm:"column:nama"`
	BeratMin int   `gorm:"column:berat_min"`
	BeratMax int   `gorm:"column:berat_max"`
}

func (Kelas) TableName() string {
	return "master_kelas"
}
