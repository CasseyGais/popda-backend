package masterofficial

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]MasterOfficial, error) {
	var officials []MasterOfficial
	err := r.DB.Find(&officials).Error
	return officials, err
}

func (r *Repository) GetByID(id uint) (*MasterOfficial, error) {
	var official MasterOfficial
	err := r.DB.First(&official, id).Error
	if err != nil {
		return nil, err
	}
	return &official, nil
}

func (r *Repository) GetByKontingenID(kontingenID uint) ([]MasterOfficial, error) {
	var officials []MasterOfficial
	err := r.DB.Where("kontingen_id = ?", kontingenID).Find(&officials).Error
	return officials, err
}

func (r *Repository) Create(official *MasterOfficial) error {
	// Omit updated_at karena tidak ada di tabel
	return r.DB.Omit("updated_at").Create(official).Error
}

func (r *Repository) Update(official *MasterOfficial) error {
	return r.DB.Model(official).Updates(map[string]interface{}{
		"kontingen_id": official.KontingenID,
		"nama":         official.Nama,
		"jabatan":      official.Jabatan,
		"no_hp":        official.NoHP,
	}).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&MasterOfficial{}, id).Error
}
