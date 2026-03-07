package tahap2

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
		Tahap2Submitted bool
		SubmittedAt     *time.Time
	}

	var meta Meta
	err := r.db.Table("kontingen_metas").
		Select("tahap2_submitted, submitted_at").
		Where("kontingen_id = ?", kontingenID).
		First(&meta).Error

	if err != nil {
		return false, nil, err
	}

	return meta.Tahap2Submitted, meta.SubmittedAt, nil
}

func (r *Repository) GetEventsWithStatus(kontingenID uint) ([]EventWithStatus, error) {

	var events []EventWithStatus

	err := r.db.Table("master_nomor mn").
		Select(`
			mn.id as event_id,
			mc.nama as cabor,
			mn.nama as nama_event,
			mn.jenis_kelamin,
			CASE WHEN ten.id IS NOT NULL THEN true ELSE false END as ikut
		`).
		Joins("JOIN master_cabor mc ON mn.cabor_id = mc.id").
		Joins("LEFT JOIN trx_kontingen_nomor ten ON ten.nomor_id = mn.id AND ten.kontingen_id = ?", kontingenID).
		Where("mn.is_active = ?", 1).
		Scan(&events).Error

	return events, err
}

func (r *Repository) SaveSelectedEvents(kontingenID uint, nomorIDs []uint) error {

	if err := r.db.Where("kontingen_id = ?", kontingenID).
		Delete(&EntryNumber{}).Error; err != nil {
		return err
	}

	if len(nomorIDs) == 0 {
		return nil
	}

	var records []EntryNumber
	for _, id := range nomorIDs {
		records = append(records, EntryNumber{
			KontingenID: kontingenID,
			NomorID:     id,
		})
	}

	return r.db.Create(&records).Error
}

func (r *Repository) SetSubmitted(kontingenID uint) error {
	return r.db.Table("kontingen_metas").
		Where("kontingen_id = ?", kontingenID).
		Updates(map[string]interface{}{
			"tahap2_submitted": true,
			"submitted_at":     time.Now(),
		}).Error
}
