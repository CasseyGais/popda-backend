package cabor

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Cabor, error) {
	var cabors []Cabor
	// Ambil semua cabor — superadmin panel perlu lihat semua data
	err := r.DB.Find(&cabors).Error
	return cabors, err
}

func (r *Repository) GetAllIncludeInactive() ([]Cabor, error) {
	var cabors []Cabor
	err := r.DB.Find(&cabors).Error
	return cabors, err
}

func (r *Repository) GetByID(id uint) (*Cabor, error) {
	var cabor Cabor
	err := r.DB.First(&cabor, id).Error
	if err != nil {
		return nil, err
	}
	return &cabor, nil
}

func (r *Repository) Create(cabor *Cabor) error {
	// Pakai Omit updated_at karena kolom tidak ada di tabel
	return r.DB.Omit("updated_at").Create(cabor).Error
}

func (r *Repository) Update(cabor *Cabor) error {
	// Pakai Updates dengan map agar hanya kolom yang ada yang di-update
	return r.DB.Model(cabor).Updates(map[string]interface{}{
		"nama":        cabor.Nama,
		"max_putra":   cabor.MaxPutra,
		"max_putri":   cabor.MaxPutri,
		"max_pelatih": cabor.MaxPelatih,
		"is_active":   cabor.IsActive,
	}).Error
}

func (r *Repository) Delete(id uint) error {
	// Hard delete — hapus permanen dari DB
	return r.DB.Delete(&Cabor{}, id).Error
}

func (r *Repository) HardDelete(id uint) error {
	return r.DB.Delete(&Cabor{}, id).Error
}
