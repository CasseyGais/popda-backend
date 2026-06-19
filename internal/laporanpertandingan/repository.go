package laporanpertandingan

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ===== LAPORAN PERTANDINGAN =====

// GetAll ambil semua laporan dengan filter opsional, join nama cabor + kontingen
func (r *Repository) GetAll(filter FilterLaporan) ([]LaporanDetail, error) {
	q := r.db.Table("laporan_pertandingan lp").
		Select(`
			lp.*,
			mc.nama           AS nama_cabor,
			mn.nama           AS nama_nomor,
			ka.nama_kontingen AS nama_kontingen_a,
			kb.nama_kontingen AS nama_kontingen_b
		`).
		Joins("LEFT JOIN master_cabor mc ON mc.id = lp.cabor_id").
		Joins("LEFT JOIN master_nomor mn ON mn.id = lp.nomor_id").
		Joins("LEFT JOIN kontingen ka ON ka.id = lp.kontingen_a_id").
		Joins("LEFT JOIN kontingen kb ON kb.id = lp.kontingen_b_id").
		Order("lp.tanggal_pertandingan DESC, lp.waktu_pertandingan DESC")

	if filter.CaborID > 0 {
		q = q.Where("lp.cabor_id = ?", filter.CaborID)
	}
	if filter.NomorID > 0 {
		q = q.Where("lp.nomor_id = ?", filter.NomorID)
	}
	if filter.Babak != "" {
		q = q.Where("lp.babak = ?", filter.Babak)
	}
	if filter.Tanggal != "" {
		q = q.Where("lp.tanggal_pertandingan = ?", filter.Tanggal)
	}
	if filter.Pemenang != "" {
		q = q.Where("lp.pemenang = ?", filter.Pemenang)
	}

	var result []LaporanDetail
	err := q.Scan(&result).Error
	return result, err
}

// GetByID ambil satu laporan by primary key + join
func (r *Repository) GetByID(id uint) (*LaporanDetail, error) {
	var result LaporanDetail
	err := r.db.Table("laporan_pertandingan lp").
		Select(`
			lp.*,
			mc.nama           AS nama_cabor,
			mn.nama           AS nama_nomor,
			ka.nama_kontingen AS nama_kontingen_a,
			kb.nama_kontingen AS nama_kontingen_b
		`).
		Joins("LEFT JOIN master_cabor mc ON mc.id = lp.cabor_id").
		Joins("LEFT JOIN master_nomor mn ON mn.id = lp.nomor_id").
		Joins("LEFT JOIN kontingen ka ON ka.id = lp.kontingen_a_id").
		Joins("LEFT JOIN kontingen kb ON kb.id = lp.kontingen_b_id").
		Where("lp.id = ?", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create insert laporan baru
func (r *Repository) Create(l *LaporanPertandingan) error {
	return r.db.Create(l).Error
}

// Update simpan perubahan laporan
func (r *Repository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&LaporanPertandingan{}).Where("id = ?", id).Updates(updates).Error
}

// Delete hapus laporan (child atlet ikut terhapus via CASCADE)
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&LaporanPertandingan{}, id).Error
}

// UpdateFile update kolom foto_bukti atau video_bukti
func (r *Repository) UpdateFile(id uint, column, path string) error {
	return r.db.Model(&LaporanPertandingan{}).Where("id = ?", id).Update(column, path).Error
}

// ===== LAPORAN PERTANDINGAN ATLET =====

// GetAtletBySisi ambil atlet dalam laporan per sisi.
// JOIN chain: laporan_pertandingan_atlet → master_atlet
//             → trx_pendaftaran_atlet (match cabor+nomor dari laporan) → master_cabor + master_nomor
// Subquery dihindari — pakai JOIN ke laporan_pertandingan untuk dapat cabor_id & nomor_id
func (r *Repository) GetAtletBySisi(laporanID uint, sisi string) ([]AtletSisiItem, error) {
	var result []AtletSisiItem
	err := r.db.Table("laporan_pertandingan_atlet lpa").
		Select(`
			lpa.id,
			lpa.atlet_id,
			ma.nama_lengkap,
			COALESCE(tpa.cabor_id, 0)   AS cabor_id,
			COALESCE(mc.nama, '')        AS nama_cabor,
			COALESCE(tpa.nomor_id, 0)   AS nomor_id,
			COALESCE(mn.nama, '')        AS nama_nomor,
			lpa.urutan
		`).
		Joins("JOIN master_atlet ma ON ma.id = lpa.atlet_id").
		Joins("JOIN laporan_pertandingan lp ON lp.id = lpa.laporan_id").
		Joins(`LEFT JOIN trx_pendaftaran_atlet tpa
			ON tpa.atlet_id = lpa.atlet_id
			AND tpa.cabor_id = lp.cabor_id
			AND tpa.nomor_id = lp.nomor_id`).
		Joins("LEFT JOIN master_cabor mc ON mc.id = tpa.cabor_id").
		Joins("LEFT JOIN master_nomor mn ON mn.id = tpa.nomor_id").
		Where("lpa.laporan_id = ? AND lpa.sisi = ?", laporanID, sisi).
		Order("lpa.urutan ASC").
		Scan(&result).Error
	return result, err
}

// ReplaceAtletSisi hapus semua atlet satu sisi lalu insert yang baru
// Dipakai saat create (pertama kali) dan update (replace)
func (r *Repository) ReplaceAtletSisi(laporanID uint, sisi string, atletIDs []uint) error {
	// Hapus semua atlet sisi ini dulu
	if err := r.db.Where("laporan_id = ? AND sisi = ?", laporanID, sisi).
		Delete(&LaporanPertandinganAtlet{}).Error; err != nil {
		return err
	}
	if len(atletIDs) == 0 {
		return nil
	}
	// Insert semua sekaligus
	var records []LaporanPertandinganAtlet
	for i, atletID := range atletIDs {
		records = append(records, LaporanPertandinganAtlet{
			LaporanID: laporanID,
			Sisi:      sisi,
			AtletID:   atletID,
			Urutan:    uint8(i + 1),
		})
	}
	return r.db.Create(&records).Error
}

// DeleteAtletByLaporan hapus semua atlet suatu laporan (dipanggil manual jika perlu)
func (r *Repository) DeleteAtletByLaporan(laporanID uint) error {
	return r.db.Where("laporan_id = ?", laporanID).Delete(&LaporanPertandinganAtlet{}).Error
}

// ===== DROPDOWN QUERIES =====

// GetKontingenDropdown ambil semua kontingen untuk dropdown Tim A / Tim B
func (r *Repository) GetKontingenDropdown() ([]KontingenDropdownItem, error) {
	var result []KontingenDropdownItem
	err := r.db.Table("kontingen").
		Select("id, nama_kontingen, territory_id").
		Order("nama_kontingen ASC").
		Scan(&result).Error
	return result, err
}

// GetCaborDropdown ambil semua cabor aktif untuk dropdown
func (r *Repository) GetCaborDropdown() ([]CaborDropdownItem, error) {
	var result []CaborDropdownItem
	err := r.db.Table("master_cabor").
		Select("id, nama, is_active").
		Where("is_active = 1").
		Order("nama ASC").
		Scan(&result).Error
	return result, err
}

// GetNomorDropdown ambil nomor/kelas aktif, bisa difilter by cabor_id
func (r *Repository) GetNomorDropdown(caborID uint) ([]NomorDropdownItem, error) {
	var result []NomorDropdownItem
	q := r.db.Table("master_nomor").
		Select("id, cabor_id, nama, jenis_kelamin, tipe").
		Where("is_active = 1").
		Order("cabor_id ASC, jenis_kelamin ASC, nama ASC")
	if caborID > 0 {
		q = q.Where("cabor_id = ?", caborID)
	}
	err := q.Scan(&result).Error
	return result, err
}

// GetAtletTerdaftarDropdown ambil atlet yang sudah terdaftar di trx_pendaftaran_atlet
// Bisa difilter by kontingen_id, cabor_id, nomor_id
// Dipakai untuk dropdown atlet sisi A/B saat input laporan pertandingan
func (r *Repository) GetAtletTerdaftarDropdown(kontingenID, caborID, nomorID uint) ([]AtletTerdaftarDropdownItem, error) {
	var result []AtletTerdaftarDropdownItem
	q := r.db.Table("trx_pendaftaran_atlet tpa").
		Select(`
			DISTINCT tpa.atlet_id,
			ma.nama_lengkap,
			ma.kontingen_id,
			k.nama_kontingen,
			tpa.cabor_id,
			tpa.nomor_id
		`).
		Joins("JOIN master_atlet ma ON ma.id = tpa.atlet_id").
		Joins("JOIN kontingen k ON k.id = ma.kontingen_id").
		Order("k.nama_kontingen ASC, ma.nama_lengkap ASC")

	if kontingenID > 0 {
		q = q.Where("ma.kontingen_id = ?", kontingenID)
	}
	if caborID > 0 {
		q = q.Where("tpa.cabor_id = ?", caborID)
	}
	if nomorID > 0 {
		q = q.Where("tpa.nomor_id = ?", nomorID)
	}

	err := q.Scan(&result).Error
	return result, err
}
