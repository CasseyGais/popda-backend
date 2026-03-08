package permissions

type Permission struct {
	ID          uint   `gorm:"primaryKey;column:id"`
	Name        string `gorm:"size:100;not null;column:name"`
	Description string `gorm:"type:text;column:description"`
}

func (Permission) TableName() string {
	return "permissions"
}
