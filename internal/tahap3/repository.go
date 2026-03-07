package tahap3

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMeta(kontingenID uint) (bool, *time.Time, error) {

	type Meta struct {
		Tahap3Submitted   bool
		Tahap3SubmittedAt *time.Time
	}

	var meta Meta
	err := r.db.Table("kontingen_metas").
		Select("tahap3_submitted, tahap3_submitted_at").
		Where("kontingen_id = ?", kontingenID).
		First(&meta).Error

	if err != nil {
		return false, nil, err
	}

	return meta.Tahap3Submitted, meta.Tahap3SubmittedAt, nil
}

func (r *Repository) GetAtlets(kontingenID uint) ([]Atlet, error) {
	var data []Atlet
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&data).Error
	return data, err
}

func (r *Repository) GetPelatihs(kontingenID uint) ([]Pelatih, error) {
	var data []Pelatih
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&data).Error
	return data, err
}

func (r *Repository) GetOfficials(kontingenID uint) ([]Official, error) {
	var data []Official
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&data).Error
	return data, err
}

func (r *Repository) SetSubmitted(kontingenID uint) error {
	return r.db.Table("kontingen_metas").
		Where("kontingen_id = ?", kontingenID).
		Updates(map[string]interface{}{
			"tahap3_submitted":    true,
			"tahap3_submitted_at": time.Now(),
		}).Error
}
