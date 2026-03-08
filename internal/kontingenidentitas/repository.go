package kontingenidentitas

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]KontingenIdentitas, error) {
	var identitas []KontingenIdentitas
	err := r.DB.Find(&identitas).Error
	return identitas, err
}

func (r *Repository) GetByID(id uint) (*KontingenIdentitas, error) {
	var identitas KontingenIdentitas
	err := r.DB.First(&identitas, id).Error
	if err != nil {
		return nil, err
	}
	return &identitas, nil
}

func (r *Repository) GetByKontingenID(kontingenID uint) (*KontingenIdentitas, error) {
	var identitas KontingenIdentitas
	err := r.DB.Where("kontingen_id = ?", kontingenID).First(&identitas).Error
	if err != nil {
		return nil, err
	}
	return &identitas, nil
}

func (r *Repository) Create(identitas *KontingenIdentitas) error {
	return r.DB.Create(identitas).Error
}

func (r *Repository) Update(identitas *KontingenIdentitas) error {
	return r.DB.Save(identitas).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&KontingenIdentitas{}, id).Error
}

func (r *Repository) UpdateKepalaFoto(id uint, foto string) error {
	return r.DB.Model(&KontingenIdentitas{}).Where("id = ?", id).Update("kepala_foto", foto).Error
}

func (r *Repository) UpdatePICFoto(id uint, foto string) error {
	return r.DB.Model(&KontingenIdentitas{}).Where("id = ?", id).Update("pic_foto", foto).Error
}
