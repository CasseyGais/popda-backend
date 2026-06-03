package rolepermissions

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]RolePermission, error) {
	var rolePermissions []RolePermission
	err := r.DB.Find(&rolePermissions).Error
	return rolePermissions, err
}

func (r *Repository) Create(rolePermission *RolePermission) error {
	return r.DB.Create(rolePermission).Error
}

func (r *Repository) Delete(roleID, permissionID uint) error {
	return r.DB.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&RolePermission{}).Error
}

func (r *Repository) DeleteByRoleID(roleID uint) error {
	return r.DB.Where("role_id = ?", roleID).Delete(&RolePermission{}).Error
}

func (r *Repository) DeleteByPermissionID(permissionID uint) error {
	return r.DB.Where("permission_id = ?", permissionID).Delete(&RolePermission{}).Error
}

func (r *Repository) GetByRoleID(roleID uint) ([]RolePermission, error) {
	var rolePermissions []RolePermission
	err := r.DB.Where("role_id = ?", roleID).Find(&rolePermissions).Error
	return rolePermissions, err
}

func (r *Repository) GetByPermissionID(permissionID uint) ([]RolePermission, error) {
	var rolePermissions []RolePermission
	err := r.DB.Where("permission_id = ?", permissionID).Find(&rolePermissions).Error
	return rolePermissions, err
}
