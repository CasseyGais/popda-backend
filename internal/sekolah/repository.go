package sekolah

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Sekolah, error) {
	var sekolahs []Sekolah
	err := r.DB.Where("is_active = ?", true).Find(&sekolahs).Error
	return sekolahs, err
}

func (r *Repository) GetByID(id uint) (*Sekolah, error) {
	var sekolah Sekolah
	err := r.DB.First(&sekolah, id).Error
	if err != nil {
		return nil, err
	}
	return &sekolah, nil
}

func (r *Repository) GetByNPSN(npsn string) (*Sekolah, error) {
	var sekolah Sekolah
	err := r.DB.Where("npsn = ? AND is_active = ?", npsn, true).First(&sekolah).Error
	if err != nil {
		return nil, err
	}
	return &sekolah, nil
}

func (r *Repository) Search(keyword string) ([]Sekolah, error) {
	var sekolahs []Sekolah
	err := r.DB.Where("is_active = ? AND (nama LIKE ? OR npsn LIKE ? OR kabupaten LIKE ?)", 
		true, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Find(&sekolahs).Error
	return sekolahs, err
}

func (r *Repository) Create(sekolah *Sekolah) error {
	return r.DB.Create(sekolah).Error
}

func (r *Repository) Update(sekolah *Sekolah) error {
	return r.DB.Save(sekolah).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Model(&Sekolah{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *Repository) HardDelete(id uint) error {
	return r.DB.Delete(&Sekolah{}, id).Error
}
