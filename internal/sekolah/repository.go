package sekolah

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Sekolah, error) {
	var sekolahs []Sekolah
	err := r.DB.Find(&sekolahs).Error
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
	err := r.DB.Where("npsn = ?", npsn).First(&sekolah).Error
	if err != nil {
		return nil, err
	}
	return &sekolah, nil
}

func (r *Repository) Search(keyword string) ([]Sekolah, error) {
	var sekolahs []Sekolah
	err := r.DB.Where("nama LIKE ? OR npsn LIKE ? OR kabupaten LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Find(&sekolahs).Error
	return sekolahs, err
}

func (r *Repository) Create(sekolah *Sekolah) error {
	// Omit updated_at dan is_active karena tidak ada di tabel
	return r.DB.Omit("updated_at", "is_active").Create(sekolah).Error
}

func (r *Repository) Update(sekolah *Sekolah) error {
	return r.DB.Model(sekolah).Updates(map[string]interface{}{
		"nama":      sekolah.Name,
		"npsn":      sekolah.NPSN,
		"alamat":    sekolah.Alamat,
		"kabupaten": sekolah.Kabupaten,
	}).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&Sekolah{}, id).Error
}
