package cabor

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Cabor, error) {
	var cabors []Cabor
	err := r.DB.Where("is_active = ?", true).Find(&cabors).Error
	return cabors, err
}

func (r *Repository) GetByID(id uint) (*Cabor, error) {
	var cabor Cabor
	err := r.DB.First(&cabor, id).Error
	if err != nil {
		return nil, err
	}
	return &cabor, nil
}

func (r *Repository) Create(cabor *Cabor) error {
	return r.DB.Create(cabor).Error
}

func (r *Repository) Update(cabor *Cabor) error {
	return r.DB.Save(cabor).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Model(&Cabor{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *Repository) HardDelete(id uint) error {
	return r.DB.Delete(&Cabor{}, id).Error
}
