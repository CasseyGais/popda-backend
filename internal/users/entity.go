package users

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	Name      string    `gorm:"size:150;not null;column:name" json:"name"`
	Email     string    `gorm:"size:150;not null;column:email" json:"email"`
	Password  string    `gorm:"size:255;not null;column:password" json:"-"` // tidak pernah dikirim ke client
	IsActive  bool      `gorm:"default:true;column:is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	Avatar    string    `gorm:"size:255;column:avatar" json:"avatar"`
}

func (User) TableName() string {
	return "users"
}
