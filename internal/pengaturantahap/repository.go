package pengaturantahap

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetAll ambil semua 3 row pengaturan tahap (tahap 1, 2, 3)
func (r *Repository) GetAll() ([]PengaturanTahap, error) {
	var result []PengaturanTahap
	err := r.db.Order("tahap ASC").Find(&result).Error
	return result, err
}

// GetByTahap ambil pengaturan untuk satu tahap (1, 2, atau 3)
func (r *Repository) GetByTahap(tahap uint) (*PengaturanTahap, error) {
	var p PengaturanTahap
	err := r.db.Where("tahap = ?", tahap).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// IsOpen cek apakah tahap tertentu sedang dibuka.
// Dipakai oleh middleware checkTahapOpen di luar package ini.
func (r *Repository) IsOpen(tahap uint) (bool, error) {
	var p PengaturanTahap
	err := r.db.Select("is_open").Where("tahap = ?", tahap).First(&p).Error
	if err != nil {
		return false, err
	}
	return p.IsOpen, nil
}

// Update simpan perubahan pengaturan satu tahap.
// Hanya field yang ada di map yang diupdate (partial update).
func (r *Repository) Update(tahap uint, updates map[string]interface{}) (*PengaturanTahap, error) {
	err := r.db.Model(&PengaturanTahap{}).
		Where("tahap = ?", tahap).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return r.GetByTahap(tahap)
}

// CountPernnahDibuka cek apakah suatu tahap pernah dibuka (is_open pernah = 1).
// Dipakai validasi urutan: tahap 2 tidak bisa dibuka sebelum tahap 1 pernah dibuka.
// Cara deteksi: cek apakah tanggal_buka terisi ATAU saat ini is_open = true.
func (r *Repository) PernahDibuka(tahap uint) (bool, error) {
	var count int64
	err := r.db.Model(&PengaturanTahap{}).
		Where("tahap = ? AND (is_open = 1 OR tanggal_buka IS NOT NULL)", tahap).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
