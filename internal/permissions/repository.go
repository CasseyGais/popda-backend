package permissions

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Permission, error) {
	var permissions []Permission
	err := r.DB.Find(&permissions).Error
	return permissions, err
}

func (r *Repository) GetByID(id uint) (*Permission, error) {
	var permission Permission
	err := r.DB.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *Repository) GetByName(name string) (*Permission, error) {
	var permission Permission
	err := r.DB.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *Repository) Create(permission *Permission) error {
	return r.DB.Create(permission).Error
}

func (r *Repository) Update(permission *Permission) error {
	return r.DB.Save(permission).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&Permission{}, id).Error
}

func (r *Repository) GetByRoleID(roleID uint) ([]Permission, error) {
	var permissions []Permission
	err := r.DB.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}
