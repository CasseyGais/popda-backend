package territories

type Territory struct {
	ID   uint   `gorm:"primaryKey;column:id"`
	Name string `gorm:"size:100;not null;column:name"`
	Type string `gorm:"type:enum('PROVINSI','KABUPATEN','KOTA');not null;column:type"`
}

func (Territory) TableName() string {
	return "territories"
}
