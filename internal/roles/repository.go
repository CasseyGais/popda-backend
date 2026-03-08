package roles

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Role, error) {
	var roles []Role
	err := r.DB.Find(&roles).Error
	return roles, err
}

func (r *Repository) GetByID(id uint) (*Role, error) {
	var role Role
	err := r.DB.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) GetByName(name string) (*Role, error) {
	var role Role
	err := r.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) Create(role *Role) error {
	return r.DB.Create(role).Error
}

func (r *Repository) Update(role *Role) error {
	return r.DB.Save(role).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&Role{}, id).Error
}

func (r *Repository) GetByUserID(userID uint) ([]Role, error) {
	var roles []Role
	err := r.DB.Table("roles").
		Joins("INNER JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *Repository) AssignPermission(roleID, permissionID uint) error {
	assignment := map[string]interface{}{
		"role_id":       roleID,
		"permission_id": permissionID,
	}
	return r.DB.Table("role_permissions").Create(assignment).Error
}

func (r *Repository) RemovePermission(roleID, permissionID uint) error {
	return r.DB.Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&map[string]interface{}{}).Error
}

func (r *Repository) GetRolePermissions(roleID uint) ([]uint, error) {
	var permissionIDs []uint
	err := r.DB.Table("role_permissions").
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissionIDs).Error
	return permissionIDs, err
}
