package userroles

import (
	"gorm.io/gorm"
)

type UserRole struct {
	UserID uint `gorm:"primaryKey;column:user_id"`
	RoleID uint `gorm:"primaryKey;column:role_id"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]UserRole, error) {
	var userRoles []UserRole
	err := r.DB.Find(&userRoles).Error
	return userRoles, err
}

func (r *Repository) Create(userRole *UserRole) error {
	return r.DB.Create(userRole).Error
}

func (r *Repository) Delete(userID, roleID uint) error {
	return r.DB.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{}).Error
}

func (r *Repository) DeleteByUserID(userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&UserRole{}).Error
}

func (r *Repository) DeleteByRoleID(roleID uint) error {
	return r.DB.Where("role_id = ?", roleID).Delete(&UserRole{}).Error
}

func (r *Repository) GetByUserID(userID uint) ([]UserRole, error) {
	var userRoles []UserRole
	err := r.DB.Where("user_id = ?", userID).Find(&userRoles).Error
	return userRoles, err
}

func (r *Repository) GetByRoleID(roleID uint) ([]UserRole, error) {
	var userRoles []UserRole
	err := r.DB.Where("role_id = ?", roleID).Find(&userRoles).Error
	return userRoles, err
}

func (r *Repository) GetRoleIDsByUserID(userID uint) ([]uint, error) {
	var roleIDs []uint
	err := r.DB.Table("user_roles").
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

func (r *Repository) GetUserIDsByRoleID(roleID uint) ([]uint, error) {
	var userIDs []uint
	err := r.DB.Table("user_roles").
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
