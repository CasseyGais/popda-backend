package modules

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Module, error) {
	var modules []Module
	err := r.DB.Order("id ASC").Find(&modules).Error
	return modules, err
}

func (r *Repository) GetByID(id uint) (*Module, error) {
	var module Module
	err := r.DB.First(&module, id).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *Repository) GetByCode(code string) (*Module, error) {
	var module Module
	err := r.DB.Where("code = ?", code).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *Repository) Create(module *Module) error {
	return r.DB.Create(module).Error
}

func (r *Repository) Update(module *Module) error {
	return r.DB.Save(module).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&Module{}, id).Error
}
