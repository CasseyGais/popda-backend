package roles

// Tabel roles: id, name, description (tidak ada created_at / updated_at)
type Role struct {
	ID          uint   `gorm:"primaryKey;column:id" json:"id"`
	Name        string `gorm:"size:50;not null;column:name" json:"name"`
	Description string `gorm:"type:text;column:description" json:"description"`
}

func (Role) TableName() string {
	return "roles"
}
