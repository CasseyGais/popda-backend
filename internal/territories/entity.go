package territories

type Territory struct {
	ID   uint   `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"size:100;not null;column:name" json:"name"`
	Type string `gorm:"type:enum('PROVINSI','KABUPATEN','KOTA');not null;column:type" json:"type"`
}

func (Territory) TableName() string {
	return "territories"
}
