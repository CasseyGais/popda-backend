package users

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"size:150;not null;column:name"`
	Email     string    `gorm:"size:150;not null;column:email"`
	Password  string    `gorm:"size:255;not null;column:password"`
	IsActive  bool      `gorm:"default:true;column:is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	Avatar    string    `gorm:"size:255;column:avatar"`
}

func (User) TableName() string {
	return "users"
}
