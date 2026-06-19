package sertifikat

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetAll ambil semua sertifikat dengan filter opsional.
// Filter by tipe_penerima, atlet_id, pelatih_id, atau official_id.
func (r *Repository) GetAll(filter map[string]interface{}) ([]Sertifikat, error) {
	var data []Sertifikat
	q := r.db.Order("created_at DESC")
	for k, v := range filter {
		q = q.Where(k+" = ?", v)
	}
	err := q.Find(&data).Error
	return data, err
}

// GetByID ambil satu sertifikat by primary key.
func (r *Repository) GetByID(id uint) (*Sertifikat, error) {
	var data Sertifikat
	err := r.db.First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// Create insert sertifikat baru.
func (r *Repository) Create(s *Sertifikat) error {
	return r.db.Create(s).Error
}

// Update simpan perubahan sertifikat (partial — hanya field yang berubah).
func (r *Repository) Update(id uint, updates map[string]interface{}) (*Sertifikat, error) {
	err := r.db.Model(&Sertifikat{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// Delete hapus sertifikat by id (hard delete).
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Sertifikat{}, id).Error
}

// UpdateFile update kolom file_sertifikat.
func (r *Repository) UpdateFile(id uint, path string) error {
	return r.db.Model(&Sertifikat{}).Where("id = ?", id).Update("file_sertifikat", path).Error
}

// ===== LOOKUP NAMA PENERIMA dari tabel master =====

// GetNamaAtlet ambil nama_lengkap dari master_atlet by id.
func (r *Repository) GetNamaAtlet(id uint) (string, error) {
	var nama string
	err := r.db.Table("master_atlet").
		Select("nama_lengkap").
		Where("id = ?", id).
		Scan(&nama).Error
	if err != nil || nama == "" {
		return "", gorm.ErrRecordNotFound
	}
	return nama, nil
}

// GetNamaPelatih ambil nama_lengkap dari master_pelatih by id.
func (r *Repository) GetNamaPelatih(id uint) (string, error) {
	var nama string
	err := r.db.Table("master_pelatih").
		Select("nama_lengkap").
		Where("id = ?", id).
		Scan(&nama).Error
	if err != nil || nama == "" {
		return "", gorm.ErrRecordNotFound
	}
	return nama, nil
}

// GetNamaOfficial ambil nama_lengkap dari master_official by id.
func (r *Repository) GetNamaOfficial(id uint) (string, error) {
	var nama string
	err := r.db.Table("master_official").
		Select("nama_lengkap").
		Where("id = ?", id).
		Scan(&nama).Error
	if err != nil || nama == "" {
		return "", gorm.ErrRecordNotFound
	}
	return nama, nil
}

// ===== DROPDOWN DATA untuk form sertifikat =====
// Endpoint khusus SUPERADMIN/STAFF_LAPANGAN — tidak difilter by kontingen.

// PenerimaSingkat adalah data minimal untuk dropdown form sertifikat.
type PenerimaSingkat struct {
	ID            uint   `gorm:"column:id" json:"id"`
	NamaLengkap   string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	KontingenID   uint   `gorm:"column:kontingen_id" json:"kontingen_id"`
	NamaKontingen string `gorm:"column:nama_kontingen" json:"nama_kontingen"`
}

// GetAllAtletSingkat ambil id + nama_lengkap + kontingen semua atlet — untuk dropdown.
func (r *Repository) GetAllAtletSingkat() ([]PenerimaSingkat, error) {
	var data []PenerimaSingkat
	err := r.db.Table("master_atlet ma").
		Select("ma.id, ma.nama_lengkap, ma.kontingen_id, k.nama_kontingen").
		Joins("LEFT JOIN kontingen k ON k.id = ma.kontingen_id").
		Order("k.nama_kontingen ASC, ma.nama_lengkap ASC").
		Scan(&data).Error
	return data, err
}

// GetAllPelatihSingkat ambil id + nama_lengkap + kontingen semua pelatih — untuk dropdown.
func (r *Repository) GetAllPelatihSingkat() ([]PenerimaSingkat, error) {
	var data []PenerimaSingkat
	err := r.db.Table("master_pelatih mp").
		Select("mp.id, mp.nama_lengkap, mp.kontingen_id, k.nama_kontingen").
		Joins("LEFT JOIN kontingen k ON k.id = mp.kontingen_id").
		Order("k.nama_kontingen ASC, mp.nama_lengkap ASC").
		Scan(&data).Error
	return data, err
}

// GetAllOfficialSingkat ambil id + nama_lengkap + kontingen semua official — untuk dropdown.
func (r *Repository) GetAllOfficialSingkat() ([]PenerimaSingkat, error) {
	var data []PenerimaSingkat
	err := r.db.Table("master_official mo").
		Select("mo.id, mo.nama_lengkap, mo.kontingen_id, k.nama_kontingen").
		Joins("LEFT JOIN kontingen k ON k.id = mo.kontingen_id").
		Order("k.nama_kontingen ASC, mo.nama_lengkap ASC").
		Scan(&data).Error
	return data, err
}

// ===== DROPDOWN PENERIMA =====

// DropdownItem dipakai untuk populate dropdown di form buat sertifikat
type DropdownItem struct {
	ID            uint   `gorm:"column:id" json:"id"`
	NamaLengkap   string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	KontingenID   uint   `gorm:"column:kontingen_id" json:"kontingen_id"`
	NamaKontingen string `gorm:"column:nama_kontingen" json:"nama_kontingen"`
}

// GetAtletDropdown ambil semua atlet dari semua kontingen untuk dropdown.
func (r *Repository) GetAtletDropdown() ([]DropdownItem, error) {
	var data []DropdownItem
	err := r.db.Table("master_atlet ma").
		Select("ma.id, ma.nama_lengkap, ma.kontingen_id, k.nama_kontingen").
		Joins("JOIN kontingen k ON k.id = ma.kontingen_id").
		Order("k.nama_kontingen ASC, ma.nama_lengkap ASC").
		Scan(&data).Error
	return data, err
}

// GetPelatihDropdown ambil semua pelatih dari semua kontingen untuk dropdown.
func (r *Repository) GetPelatihDropdown() ([]DropdownItem, error) {
	var data []DropdownItem
	err := r.db.Table("master_pelatih mp").
		Select("mp.id, mp.nama_lengkap, mp.kontingen_id, k.nama_kontingen").
		Joins("JOIN kontingen k ON k.id = mp.kontingen_id").
		Order("k.nama_kontingen ASC, mp.nama_lengkap ASC").
		Scan(&data).Error
	return data, err
}

// GetOfficialDropdown ambil semua official dari semua kontingen untuk dropdown.
func (r *Repository) GetOfficialDropdown() ([]DropdownItem, error) {
	var data []DropdownItem
	err := r.db.Table("master_official mo").
		Select("mo.id, mo.nama_lengkap, mo.kontingen_id, k.nama_kontingen").
		Joins("JOIN kontingen k ON k.id = mo.kontingen_id").
		Order("k.nama_kontingen ASC, mo.nama_lengkap ASC").
		Scan(&data).Error
	return data, err
}
