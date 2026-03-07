package tahap1

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMeta(kontingenID uint) (*KontingenMeta, error) {
	var meta KontingenMeta
	err := r.db.Where("kontingen_id = ?", kontingenID).First(&meta).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &KontingenMeta{KontingenID: kontingenID}, nil
		}
		return nil, err
	}
	return &meta, nil
}

func (r *Repository) SaveMeta(kontingenID uint, atlet, pelatih, official int) error {
	var meta KontingenMeta
	err := r.db.Where("kontingen_id = ?", kontingenID).First(&meta).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		meta = KontingenMeta{
			KontingenID:    kontingenID,
			JumlahAtlet:    atlet,
			JumlahPelatih:  pelatih,
			JumlahOfficial: official,
		}
		return r.db.Create(&meta).Error
	}

	meta.JumlahAtlet = atlet
	meta.JumlahPelatih = pelatih
	meta.JumlahOfficial = official

	return r.db.Save(&meta).Error
}

func (r *Repository) SaveCabor(kontingenID uint, list []uint) error {
	r.db.Where("kontingen_id = ?", kontingenID).Delete(&CaborKontingen{})

	for _, c := range list {
		r.db.Create(&CaborKontingen{
			KontingenID: kontingenID,
			CaborID:     c,
		})
	}
	return nil
}

func (r *Repository) GetCabor(kontingenID uint) ([]uint, error) {
	var result []uint
	err := r.db.Model(&CaborKontingen{}).
		Where("kontingen_id = ?", kontingenID).
		Pluck("cabor_id", &result).Error
	return result, err
}

func (r *Repository) SetSubmitted(kontingenID uint) error {
	now := time.Now()
	return r.db.Model(&KontingenMeta{}).
		Where("kontingen_id = ?", kontingenID).
		Updates(map[string]interface{}{
			"tahap1_submitted": true,
			"submitted_at":     now,
		}).Error
}
