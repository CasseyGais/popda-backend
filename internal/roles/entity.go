package roles

type Role struct {
	ID          uint   `gorm:"primaryKey;column:id"`
	Name        string `gorm:"size:50;not null;column:name"`
	Description string `gorm:"type:text;column:description"`
}

func (Role) TableName() string {
	return "roles"
}
