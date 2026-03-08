package nomor

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Nomor, error) {
	var nomors []Nomor
	err := r.DB.Preload("Cabor").Where("is_active = ?", true).Find(&nomors).Error
	return nomors, err
}

func (r *Repository) GetByID(id uint) (*Nomor, error) {
	var nomor Nomor
	err := r.DB.Preload("Cabor").First(&nomor, id).Error
	if err != nil {
		return nil, err
	}
	return &nomor, nil
}

func (r *Repository) GetByCaborID(caborID uint) ([]Nomor, error) {
	var nomors []Nomor
	err := r.DB.Preload("Cabor").Where("cabor_id = ? AND is_active = ?", caborID, true).Find(&nomors).Error
	return nomors, err
}

func (r *Repository) Create(nomor *Nomor) error {
	return r.DB.Create(nomor).Error
}

func (r *Repository) Update(nomor *Nomor) error {
	return r.DB.Save(nomor).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Model(&Nomor{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *Repository) HardDelete(id uint) error {
	return r.DB.Delete(&Nomor{}, id).Error
}
