package masterofficial

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// GetKontingenIDByTerritory untuk resolve superadmin → kontingen
func (r *Repository) GetKontingenIDByTerritory(territoryID uint) (uint, error) {
	var kontingenID uint
	err := r.DB.Table("kontingen").
		Select("id").
		Where("territory_id = ?", territoryID).
		Scan(&kontingenID).Error
	if err != nil {
		return 0, err
	}
	if kontingenID == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return kontingenID, nil
}

// GetByKontingenID ambil semua official milik kontingen
func (r *Repository) GetByKontingenID(kontingenID uint) ([]MasterOfficial, error) {
	var data []MasterOfficial
	err := r.DB.Where("kontingen_id = ?", kontingenID).
		Order("nama_lengkap ASC").
		Find(&data).Error
	return data, err
}

// GetByID ambil satu official
func (r *Repository) GetByID(id uint) (*MasterOfficial, error) {
	var data MasterOfficial
	err := r.DB.First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// Create insert official baru
func (r *Repository) Create(o *MasterOfficial) error {
	return r.DB.Create(o).Error
}

// Update simpan perubahan official
func (r *Repository) Update(o *MasterOfficial) error {
	return r.DB.Save(o).Error
}

// Delete hapus official by id
func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&MasterOfficial{}, id).Error
}

// UpdateFile update satu kolom file/foto
func (r *Repository) UpdateFile(id uint, column, path string) error {
	return r.DB.Model(&MasterOfficial{}).Where("id = ?", id).Update(column, path).Error
}

// ===== TRX PENDAFTARAN OFFICIAL =====

// GetTrxByKontingen ambil semua trx official milik kontingen (via JOIN)
func (r *Repository) GetTrxByKontingen(kontingenID uint) ([]TrxPendaftaranOfficial, error) {
	var data []TrxPendaftaranOfficial
	err := r.DB.Table("trx_pendaftaran_official tpo").
		Select("tpo.*").
		Joins("JOIN master_official mo ON mo.id = tpo.official_id").
		Where("mo.kontingen_id = ?", kontingenID).
		Find(&data).Error
	return data, err
}

// CreateTrx insert ke trx_pendaftaran_official
func (r *Repository) CreateTrx(trx *TrxPendaftaranOfficial) error {
	return r.DB.Create(trx).Error
}

// DeleteTrx hapus trx by id
func (r *Repository) DeleteTrx(id uint) error {
	return r.DB.Delete(&TrxPendaftaranOfficial{}, id).Error
}

// DeleteTrxByOfficial hapus semua trx milik satu official (dipanggil saat official dihapus)
func (r *Repository) DeleteTrxByOfficial(officialID uint) error {
	return r.DB.Where("official_id = ?", officialID).Delete(&TrxPendaftaranOfficial{}).Error
}
