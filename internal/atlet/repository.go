package atlet

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Atlet, error) {
	var atlets []Atlet
	err := r.DB.Find(&atlets).Error
	return atlets, err
}

func (r *Repository) GetByID(id uint) (*Atlet, error) {
	var atlet Atlet
	err := r.DB.First(&atlet, id).Error
	if err != nil {
		return nil, err
	}
	return &atlet, nil
}

func (r *Repository) GetByKontingenID(kontingenID uint) ([]Atlet, error) {
	var atlets []Atlet
	err := r.DB.Where("kontingen_id = ?", kontingenID).Find(&atlets).Error
	return atlets, err
}

func (r *Repository) GetBySekolahID(sekolahID uint) ([]Atlet, error) {
	var atlets []Atlet
	err := r.DB.Where("sekolah_id = ?", sekolahID).Find(&atlets).Error
	return atlets, err
}

func (r *Repository) GetByStatus(status string) ([]Atlet, error) {
	var atlets []Atlet
	err := r.DB.Where("status_verifikasi = ?", status).Find(&atlets).Error
	return atlets, err
}

func (r *Repository) Create(atlet *Atlet) error {
	return r.DB.Create(atlet).Error
}

func (r *Repository) Update(atlet *Atlet) error {
	return r.DB.Save(atlet).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&Atlet{}, id).Error
}

func (r *Repository) UpdateStatus(id uint, status string) error {
	return r.DB.Model(&Atlet{}).Where("id = ?", id).Update("status_verifikasi", status).Error
}

func (r *Repository) UpdateFoto(id uint, foto string) error {
	return r.DB.Model(&Atlet{}).Where("id = ?", id).Update("foto", foto).Error
}
