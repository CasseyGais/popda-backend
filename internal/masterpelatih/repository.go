package masterpelatih

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]MasterPelatih, error) {
	var pelatihs []MasterPelatih
	err := r.DB.Find(&pelatihs).Error
	return pelatihs, err
}

func (r *Repository) GetByID(id uint) (*MasterPelatih, error) {
	var pelatih MasterPelatih
	err := r.DB.First(&pelatih, id).Error
	if err != nil {
		return nil, err
	}
	return &pelatih, nil
}

func (r *Repository) GetByKontingenID(kontingenID uint) ([]MasterPelatih, error) {
	var pelatihs []MasterPelatih
	err := r.DB.Where("kontingen_id = ?", kontingenID).Find(&pelatihs).Error
	return pelatihs, err
}

func (r *Repository) Create(pelatih *MasterPelatih) error {
	return r.DB.Create(pelatih).Error
}

func (r *Repository) Update(pelatih *MasterPelatih) error {
	return r.DB.Save(pelatih).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&MasterPelatih{}, id).Error
}

func (r *Repository) UpdateFoto(id uint, foto string) error {
	return r.DB.Model(&MasterPelatih{}).Where("id = ?", id).Update("foto", foto).Error
}

func (r *Repository) UpdateSertifikat(id uint, sertifikat string) error {
	return r.DB.Model(&MasterPelatih{}).Where("id = ?", id).Update("sertifikat", sertifikat).Error
}
