package masterpelatih

import (
	"gorm.io/gorm"
)

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

// GetByKontingenID ambil semua pelatih milik kontingen
func (r *Repository) GetByKontingenID(kontingenID uint) ([]MasterPelatih, error) {
	var data []MasterPelatih
	err := r.DB.Where("kontingen_id = ?", kontingenID).
		Order("nama_lengkap ASC").
		Find(&data).Error
	return data, err
}

// GetByID ambil satu pelatih by primary key
func (r *Repository) GetByID(id uint) (*MasterPelatih, error) {
	var data MasterPelatih
	err := r.DB.First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// Create insert pelatih baru
func (r *Repository) Create(p *MasterPelatih) error {
	return r.DB.Create(p).Error
}

// Update simpan perubahan pelatih
func (r *Repository) Update(p *MasterPelatih) error {
	return r.DB.Save(p).Error
}

// Delete hapus pelatih by id
func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&MasterPelatih{}, id).Error
}

// UpdateFile update satu kolom file/foto
func (r *Repository) UpdateFile(id uint, column, path string) error {
	return r.DB.Model(&MasterPelatih{}).Where("id = ?", id).Update(column, path).Error
}

// ===== TRX PENDAFTARAN PELATIH =====

// GetTrxByKontingen ambil semua trx pelatih milik kontingen (via JOIN)
func (r *Repository) GetTrxByKontingen(kontingenID uint) ([]TrxPendaftaranPelatih, error) {
	var data []TrxPendaftaranPelatih
	err := r.DB.Table("trx_pendaftaran_pelatih tpp").
		Select("tpp.*").
		Joins("JOIN master_pelatih mp ON mp.id = tpp.pelatih_id").
		Where("mp.kontingen_id = ?", kontingenID).
		Find(&data).Error
	return data, err
}

// CreateTrx insert ke trx_pendaftaran_pelatih
func (r *Repository) CreateTrx(trx *TrxPendaftaranPelatih) error {
	return r.DB.Create(trx).Error
}

// DeleteTrx hapus trx by id
func (r *Repository) DeleteTrx(id uint) error {
	return r.DB.Delete(&TrxPendaftaranPelatih{}, id).Error
}

// DeleteTrxByPelatih hapus semua trx milik satu pelatih (dipanggil saat pelatih dihapus)
func (r *Repository) DeleteTrxByPelatih(pelatihID uint) error {
	return r.DB.Where("pelatih_id = ?", pelatihID).Delete(&TrxPendaftaranPelatih{}).Error
}

// GetCaborKontingen ambil cabor_id yang sudah dipilih kontingen di tahap 1
func (r *Repository) GetCaborKontingen(kontingenID uint) ([]uint, error) {
	var caborIDs []uint
	err := r.DB.Table("trx_kontingen_cabor").
		Where("kontingen_id = ?", kontingenID).
		Pluck("cabor_id", &caborIDs).Error
	return caborIDs, err
}
