package kontingen

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ================= GET =================
func (r *Repository) GetByKontingenID(kontingenID uint) (*IdentitasKontingen, error) {
	var data IdentitasKontingen

	err := r.db.
		Where("kontingen_id = ?", kontingenID).
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

// ================= CREATE =================
func (r *Repository) Create(data *IdentitasKontingen) error {
	return r.db.Create(data).Error
}

// ================= KONTINGEN HELPERS =================
func (r *Repository) KontingenExists(kontingenID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Kontingen{}).Where("id = ?", kontingenID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateKontingenForID(kontingenID uint) error {
	// Buat kontingen dengan territory_id yang sama
	kontingen := Kontingen{
		ID:            kontingenID,
		TerritoryID:   kontingenID, // Fallback: gunakan kontingenID sebagai territoryID
		NamaKontingen: fmt.Sprintf("Kontingen %d", kontingenID),
	}

	return r.db.Create(&kontingen).Error
}

// ================= UPDATE =================
func (r *Repository) Update(kontingenID uint, data *IdentitasKontingen) error {
	return r.db.
		Model(&IdentitasKontingen{}).
		Where("kontingen_id = ?", kontingenID).
		Updates(map[string]interface{}{
			"kepala_nama":    data.KepalaNama,
			"kepala_jabatan": data.KepalaJabatan,
			"kepala_nip":     data.KepalaNIP,
			"kepala_telepon": data.KepalaTelepon,
			"kepala_foto":    data.KepalaFoto,
			"pic_nama":       data.PICNama,
			"pic_jabatan":    data.PICJabatan,
			"pic_telepon":    data.PICTelepon,
			"pic_foto":       data.PICFoto,
			"alamat":         data.Alamat,
			"email_instansi": data.EmailInstansi,
			"phone_instansi": data.PhoneInstansi,
			"updated_at":     gorm.Expr("NOW()"),
		}).Error
}

// ================= COUNT =================
func (r *Repository) CountByKontingenID(kontingenID uint) (int64, error) {
	var count int64

	err := r.db.
		Model(&IdentitasKontingen{}).
		Where("kontingen_id = ?", kontingenID).
		Count(&count).Error

	return count, err
}

// ================= KONTINGEN METHODS =================
func (r *Repository) GetKontingenByID(id uint) (*Kontingen, error) {
	var kontingen Kontingen

	err := r.db.
		Where("id = ?", id).
		First(&kontingen).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &kontingen, nil
}

func (r *Repository) GetKontingenByTerritoryID(territoryID uint) (*Kontingen, error) {
	var kontingen Kontingen

	err := r.db.
		Where("territory_id = ?", territoryID).
		First(&kontingen).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &kontingen, nil
}

func (r *Repository) CreateKontingen(kontingen *Kontingen) error {
	return r.db.Create(kontingen).Error
}

func (r *Repository) UpdateKontingen(id uint, kontingen *Kontingen) error {
	updates := map[string]interface{}{
		"territory_id":   kontingen.TerritoryID,
		"nama_kontingen": kontingen.NamaKontingen,
		"updated_at":     gorm.Expr("NOW()"),
	}

	// Hanya update status jika nilainya valid (DRAFT atau SUBMITTED)
	if kontingen.Tahap1Status == TahapStatusDraft || kontingen.Tahap1Status == TahapStatusSubmitted {
		updates["tahap1_status"] = kontingen.Tahap1Status
	}
	if kontingen.Tahap1Submitted != nil {
		updates["tahap1_submitted_at"] = kontingen.Tahap1Submitted
	}
	if kontingen.Tahap2Status == TahapStatusDraft || kontingen.Tahap2Status == TahapStatusSubmitted {
		updates["tahap2_status"] = kontingen.Tahap2Status
	}
	if kontingen.Tahap2Submitted != nil {
		updates["tahap2_submitted_at"] = kontingen.Tahap2Submitted
	}

	return r.db.
		Model(&Kontingen{}).
		Where("id = ?", id).
		Updates(updates).Error
}
