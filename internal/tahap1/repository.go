package tahap1

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

// ===== TrxKontingenCabor =====

// GetCabor ambil semua trx_kontingen_cabor milik kontingen
func (r *Repository) GetCabor(kontingenID uint) ([]TrxKontingenCabor, error) {
	var result []TrxKontingenCabor
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&result).Error
	return result, err
}

// UpsertCabor insert jika belum ada, update jika sudah ada (berdasarkan kontingen_id + cabor_id)
func (r *Repository) UpsertCabor(trx *TrxKontingenCabor) error {
	var existing TrxKontingenCabor
	err := r.db.Where("kontingen_id = ? AND cabor_id = ?", trx.KontingenID, trx.CaborID).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Insert baru — DB trigger akan validasi kuota otomatis
		return r.db.Omit("updated_at").Create(trx).Error
	}

	// Update yang sudah ada — trigger tidak aktif, validasi sudah dilakukan di service
	return r.db.Model(&existing).Updates(map[string]interface{}{
		"putra":          trx.Putra,
		"putri":          trx.Putri,
		"pelatih":        trx.Pelatih,
		"total_atlet":    trx.TotalAtlet,
		"total_personel": trx.TotalPersonel,
	}).Error
}

// DeleteCabor hapus satu cabor dari kontingen
func (r *Repository) DeleteCabor(kontingenID, caborID uint) error {
	return r.db.Where("kontingen_id = ? AND cabor_id = ?", kontingenID, caborID).
		Delete(&TrxKontingenCabor{}).Error
}

// DeleteAllCabor hapus semua cabor milik kontingen
func (r *Repository) DeleteAllCabor(kontingenID uint) error {
	return r.db.Where("kontingen_id = ?", kontingenID).Delete(&TrxKontingenCabor{}).Error
}

// GetMasterCaborKuota ambil data kuota max dari master_cabor untuk validasi di service
func (r *Repository) GetMasterCaborKuota(caborID uint) (maxPutra, maxPutri, maxPelatih int, err error) {
	var result struct {
		MaxPutra   int `gorm:"column:max_putra"`
		MaxPutri   int `gorm:"column:max_putri"`
		MaxPelatih int `gorm:"column:max_pelatih"`
	}
	err = r.db.Table("master_cabor").
		Select("max_putra, max_putri, max_pelatih").
		Where("id = ?", caborID).
		Scan(&result).Error
	return result.MaxPutra, result.MaxPutri, result.MaxPelatih, err
}

// GetKontingenIDByTerritory cari kontingen_id berdasarkan territory_id
// Dipakai superadmin yang tidak punya kontingen_id di JWT
func (r *Repository) GetKontingenIDByTerritory(territoryID uint) (uint, error) {
	var kontingenID uint
	err := r.db.Table("kontingen").
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

// GetKontingen ambil data kontingen (termasuk status tahap)
func (r *Repository) GetKontingen(kontingenID uint) (*Kontingen, error) {
	var k Kontingen
	err := r.db.First(&k, kontingenID).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// SetTahap1Submitted set tahap1_status = SUBMITTED, isi submitted_at,
// dan otomatis set tahap1_validasi_status = PENDING untuk review superadmin.
func (r *Repository) SetTahap1Submitted(kontingenID uint) error {
	now := time.Now()
	pending := "PENDING"
	return r.db.Model(&Kontingen{}).
		Where("id = ?", kontingenID).
		Updates(map[string]interface{}{
			"tahap1_status":           "SUBMITTED",
			"tahap1_submitted_at":     now,
			"tahap1_validasi_status":  pending,
			"tahap1_validasi_catatan": nil,
			"tahap1_validasi_at":      nil,
		}).Error
}

// GetCaborWithNama ambil cabor beserta nama cabor (join ke master_cabor) untuk keperluan export
func (r *Repository) GetCaborWithNama(kontingenID uint) ([]CaborExportRow, error) {
	var result []CaborExportRow
	err := r.db.Table("trx_kontingen_cabor tkc").
		Select(`
			mc.nama        AS nama_cabor,
			tkc.putra,
			tkc.putri,
			tkc.pelatih,
			tkc.total_atlet,
			tkc.total_personel
		`).
		Joins("JOIN master_cabor mc ON mc.id = tkc.cabor_id").
		Where("tkc.kontingen_id = ?", kontingenID).
		Order("mc.nama ASC").
		Scan(&result).Error
	return result, err
}

// ResetTahap1 kembalikan tahap1_status ke DRAFT dan reset validasi ke NULL.
// Dipanggil superadmin agar kontingen bisa edit dan submit ulang.
func (r *Repository) ResetTahap1(kontingenID uint) error {
	return r.db.Model(&Kontingen{}).
		Where("id = ?", kontingenID).
		Updates(map[string]interface{}{
			"tahap1_status":           "DRAFT",
			"tahap1_submitted_at":     nil,
			"tahap1_validasi_status":  nil,
			"tahap1_validasi_catatan": nil,
			"tahap1_validasi_at":      nil,
		}).Error
}
