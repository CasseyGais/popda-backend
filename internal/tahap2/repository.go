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

// ===== Kontingen =====

// GetKontingen ambil data kontingen beserta status tahap2
func (r *Repository) GetKontingen(kontingenID uint) (*Kontingen, error) {
	var k Kontingen
	err := r.db.First(&k, kontingenID).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// GetKontingenIDByTerritory cari kontingen_id berdasarkan territory_id (untuk superadmin)
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

// SetTahap2Submitted set tahap2_status = SUBMITTED, isi submitted_at,
// dan otomatis set tahap2_validasi_status = PENDING untuk review superadmin.
func (r *Repository) SetTahap2Submitted(kontingenID uint) error {
	now := time.Now()
	pending := "PENDING"
	return r.db.Model(&Kontingen{}).
		Where("id = ?", kontingenID).
		Updates(map[string]interface{}{
			"tahap2_status":           "SUBMITTED",
			"tahap2_submitted_at":     now,
			"tahap2_validasi_status":  pending,
			"tahap2_validasi_catatan": nil,
			"tahap2_validasi_at":      nil,
		}).Error
}

// ===== Nomor =====

// GetNomorByCabor ambil semua nomor aktif dari cabor yang dipilih di tahap 1,
// beserta status apakah kontingen sudah mendaftar di nomor tsb.
func (r *Repository) GetNomorByCabor(kontingenID uint) ([]NomorWithStatus, error) {
	var result []NomorWithStatus

	err := r.db.Table("master_nomor mn").
		Select(`
			mn.id       AS nomor_id,
			mn.cabor_id AS cabor_id,
			mc.nama     AS nama_cabor,
			mn.nama     AS nama_nomor,
			mn.jenis_kelamin,
			mn.tipe,
			CASE WHEN tkn.id IS NOT NULL THEN TRUE ELSE FALSE END AS terdaftar
		`).
		Joins("JOIN master_cabor mc ON mc.id = mn.cabor_id").
		Joins(`JOIN trx_kontingen_cabor tkc 
			ON tkc.cabor_id = mn.cabor_id 
			AND tkc.kontingen_id = ?`, kontingenID).
		Joins(`LEFT JOIN trx_kontingen_nomor tkn 
			ON tkn.nomor_id = mn.id 
			AND tkn.kontingen_id = ?`, kontingenID).
		Where("mn.is_active = ?", true).
		Order("mc.nama ASC, mn.jenis_kelamin ASC, mn.nama ASC").
		Scan(&result).Error

	return result, err
}

// ===== TrxKontingenNomor =====

// GetTerdaftar ambil semua nomor_id yang sudah didaftarkan kontingen
func (r *Repository) GetTerdaftar(kontingenID uint) ([]uint, error) {
	var nomorIDs []uint
	err := r.db.Table("trx_kontingen_nomor").
		Where("kontingen_id = ?", kontingenID).
		Pluck("nomor_id", &nomorIDs).Error
	return nomorIDs, err
}

// DaftarNomor tambah satu nomor ke daftar kontingen
func (r *Repository) DaftarNomor(kontingenID, nomorID uint) error {
	trx := &TrxKontingenNomor{
		KontingenID: kontingenID,
		NomorID:     nomorID,
	}
	return r.db.Omit("updated_at").Create(trx).Error
}

// BatalNomor hapus satu nomor dari daftar kontingen
func (r *Repository) BatalNomor(kontingenID, nomorID uint) error {
	return r.db.Where("kontingen_id = ? AND nomor_id = ?", kontingenID, nomorID).
		Delete(&TrxKontingenNomor{}).Error
}

// HapusSemua hapus semua nomor kontingen (untuk replace-all)
func (r *Repository) HapusSemua(kontingenID uint) error {
	return r.db.Where("kontingen_id = ?", kontingenID).
		Delete(&TrxKontingenNomor{}).Error
}

// IsNomorDariCaborKontingen cek apakah nomor_id berasal dari cabor yang dipilih kontingen di tahap 1
func (r *Repository) IsNomorDariCaborKontingen(kontingenID, nomorID uint) (bool, error) {
	var count int64
	err := r.db.Table("master_nomor mn").
		Joins("JOIN trx_kontingen_cabor tkc ON tkc.cabor_id = mn.cabor_id").
		Where("mn.id = ? AND tkc.kontingen_id = ?", nomorID, kontingenID).
		Count(&count).Error
	return count > 0, err
}

// GetNomorTerdaftarForExport ambil nomor yang sudah didaftarkan kontingen beserta detail untuk export
func (r *Repository) GetNomorTerdaftarForExport(kontingenID uint) ([]NomorExportRow, error) {
	var result []NomorExportRow
	err := r.db.Table("trx_kontingen_nomor tkn").
		Select(`
			mc.nama  AS nama_cabor,
			mn.nama  AS nama_nomor,
			mn.jenis_kelamin,
			mn.tipe
		`).
		Joins("JOIN master_nomor mn ON mn.id = tkn.nomor_id").
		Joins("JOIN master_cabor mc ON mc.id = mn.cabor_id").
		Where("tkn.kontingen_id = ?", kontingenID).
		Order("mc.nama ASC, mn.jenis_kelamin ASC, mn.nama ASC").
		Scan(&result).Error
	return result, err
}
