package nomor

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Nomor, error) {
	var nomors []Nomor
	// Ambil semua nomor aktif
	err := r.DB.Preload("Cabor").Where("is_active = ?", true).Find(&nomors).Error
	return nomors, err
}

func (r *Repository) GetByID(id uint) (*Nomor, error) {
	var nomor Nomor
	err := r.DB.Preload("Cabor").First(&nomor, id).Error
	if err != nil {
		return nil, err
	}
	return &nomor, nil
}

func (r *Repository) GetByCaborID(caborID uint) ([]Nomor, error) {
	var nomors []Nomor
	err := r.DB.Preload("Cabor").Where("cabor_id = ? AND is_active = ?", caborID, true).Find(&nomors).Error
	return nomors, err
}

func (r *Repository) Create(nomor *Nomor) error {
	// Omit updated_at karena kolom tidak ada di tabel
	return r.DB.Omit("updated_at").Create(nomor).Error
}

func (r *Repository) Update(nomor *Nomor) error {
	return r.DB.Model(nomor).Updates(map[string]interface{}{
		"cabor_id":      nomor.CaborID,
		"nama":          nomor.Nama,
		"jenis_kelamin": nomor.JenisKelamin,
		"tipe":          nomor.Tipe,
		"is_active":     nomor.IsActive,
	}).Error
}

func (r *Repository) Delete(id uint) error {
	// Hard delete — hapus permanen dari DB
	return r.DB.Delete(&Nomor{}, id).Error
}

func (r *Repository) HardDelete(id uint) error {
	return r.DB.Delete(&Nomor{}, id).Error
}
