package tahap3

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

// ===== KONTINGEN =====

func (r *Repository) GetKontingen(kontingenID uint) (*Kontingen, error) {
	var k Kontingen
	err := r.db.First(&k, kontingenID).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

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

// ===== MASTER ATLET =====

func (r *Repository) GetAtlets(kontingenID uint) ([]MasterAtlet, error) {
	var data []MasterAtlet
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&data).Error
	return data, err
}

func (r *Repository) GetAtletByID(id uint) (*MasterAtlet, error) {
	var data MasterAtlet
	err := r.db.First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *Repository) CreateAtlet(atlet *MasterAtlet) error {
	return r.db.Create(atlet).Error
}

func (r *Repository) UpdateAtlet(atlet *MasterAtlet) error {
	return r.db.Save(atlet).Error
}

func (r *Repository) DeleteAtlet(id uint) error {
	return r.db.Delete(&MasterAtlet{}, id).Error
}

func (r *Repository) UpdateAtletFoto(id uint, path string) error {
	return r.db.Model(&MasterAtlet{}).Where("id = ?", id).Update("foto", path).Error
}

func (r *Repository) UpdateAtletFile(id uint, column, path string) error {
	return r.db.Model(&MasterAtlet{}).Where("id = ?", id).Update(column, path).Error
}

// ===== MASTER PELATIH =====

func (r *Repository) GetPelatihs(kontingenID uint) ([]MasterPelatih, error) {
	var data []MasterPelatih
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&data).Error
	return data, err
}

func (r *Repository) GetPelatihByID(id uint) (*MasterPelatih, error) {
	var data MasterPelatih
	err := r.db.First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *Repository) CreatePelatih(pelatih *MasterPelatih) error {
	return r.db.Create(pelatih).Error
}

func (r *Repository) UpdatePelatih(pelatih *MasterPelatih) error {
	return r.db.Save(pelatih).Error
}

func (r *Repository) DeletePelatih(id uint) error {
	return r.db.Delete(&MasterPelatih{}, id).Error
}

func (r *Repository) UpdatePelatihFile(id uint, column, path string) error {
	return r.db.Model(&MasterPelatih{}).Where("id = ?", id).Update(column, path).Error
}

// ===== MASTER OFFICIAL =====

func (r *Repository) GetOfficials(kontingenID uint) ([]MasterOfficial, error) {
	var data []MasterOfficial
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&data).Error
	return data, err
}

func (r *Repository) GetOfficialByID(id uint) (*MasterOfficial, error) {
	var data MasterOfficial
	err := r.db.First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *Repository) CreateOfficial(official *MasterOfficial) error {
	return r.db.Create(official).Error
}

func (r *Repository) UpdateOfficial(official *MasterOfficial) error {
	return r.db.Save(official).Error
}

func (r *Repository) DeleteOfficial(id uint) error {
	return r.db.Delete(&MasterOfficial{}, id).Error
}

func (r *Repository) UpdateOfficialFile(id uint, column, path string) error {
	return r.db.Model(&MasterOfficial{}).Where("id = ?", id).Update(column, path).Error
}

// ===== TRX PENDAFTARAN ATLET =====

func (r *Repository) GetTrxAtlets(kontingenID uint) ([]TrxPendaftaranAtlet, error) {
	var data []TrxPendaftaranAtlet
	// Join ke master_atlet untuk filter by kontingen_id
	err := r.db.Table("trx_pendaftaran_atlet tpa").
		Select("tpa.*").
		Joins("JOIN master_atlet ma ON ma.id = tpa.atlet_id").
		Where("ma.kontingen_id = ?", kontingenID).
		Scan(&data).Error
	return data, err
}

func (r *Repository) CreateTrxAtlet(trx *TrxPendaftaranAtlet) error {
	return r.db.Create(trx).Error
}

func (r *Repository) DeleteTrxAtlet(id uint) error {
	return r.db.Delete(&TrxPendaftaranAtlet{}, id).Error
}

// ===== TRX PENDAFTARAN PELATIH =====

func (r *Repository) GetTrxPelatihs(kontingenID uint) ([]TrxPendaftaranPelatih, error) {
	var data []TrxPendaftaranPelatih
	err := r.db.Table("trx_pendaftaran_pelatih tpp").
		Select("tpp.*").
		Joins("JOIN master_pelatih mp ON mp.id = tpp.pelatih_id").
		Where("mp.kontingen_id = ?", kontingenID).
		Scan(&data).Error
	return data, err
}

func (r *Repository) CreateTrxPelatih(trx *TrxPendaftaranPelatih) error {
	return r.db.Create(trx).Error
}

func (r *Repository) DeleteTrxPelatih(id uint) error {
	return r.db.Delete(&TrxPendaftaranPelatih{}, id).Error
}

// ===== TRX PENDAFTARAN OFFICIAL =====

func (r *Repository) GetTrxOfficials(kontingenID uint) ([]TrxPendaftaranOfficial, error) {
	var data []TrxPendaftaranOfficial
	err := r.db.Table("trx_pendaftaran_official tpo").
		Select("tpo.*").
		Joins("JOIN master_official mo ON mo.id = tpo.official_id").
		Where("mo.kontingen_id = ?", kontingenID).
		Scan(&data).Error
	return data, err
}

func (r *Repository) CreateTrxOfficial(trx *TrxPendaftaranOfficial) error {
	return r.db.Create(trx).Error
}

func (r *Repository) DeleteTrxOfficial(id uint) error {
	return r.db.Delete(&TrxPendaftaranOfficial{}, id).Error
}

// ===== SUBMIT TAHAP 3 =====

// SetTahap3Submitted set tahap3_status = SUBMITTED, isi submitted_at,
// dan otomatis set tahap3_validasi_status = PENDING untuk review superadmin.
func (r *Repository) SetTahap3Submitted(kontingenID uint) error {
	now := time.Now()
	pending := "PENDING"
	return r.db.Model(&Kontingen{}).
		Where("id = ?", kontingenID).
		Updates(map[string]interface{}{
			"tahap3_status":           "SUBMITTED",
			"tahap3_submitted_at":     now,
			"tahap3_validasi_status":  pending,
			"tahap3_validasi_catatan": nil,
			"tahap3_validasi_at":      nil,
		}).Error
}

// BulkInsertTrxAtlets insert semua atlet kontingen ke trx_pendaftaran_atlet
// Dipanggil saat submit tahap 3 — hanya atlet yang belum ada di trx
func (r *Repository) BulkInsertTrxAtlets(kontingenID uint) error {
	// Ambil semua atlet kontingen yang belum ada di trx_pendaftaran_atlet
	type AtletNomor struct {
		AtletID uint
		CaborID uint
		NomorID uint
	}
	var pairs []AtletNomor
	// Ambil atlet berdasarkan trx_kontingen_nomor (nomor yang dipilih di tahap 2)
	err := r.db.Raw(`
		SELECT ma.id AS atlet_id, tkc.cabor_id, tkn.nomor_id
		FROM master_atlet ma
		JOIN trx_kontingen_cabor tkc ON tkc.kontingen_id = ma.kontingen_id
		JOIN trx_kontingen_nomor tkn ON tkn.kontingen_id = ma.kontingen_id
			AND tkn.nomor_id IN (
				SELECT id FROM master_nomor mn WHERE mn.cabor_id = tkc.cabor_id
			)
		LEFT JOIN trx_pendaftaran_atlet tpa ON tpa.atlet_id = ma.id AND tpa.nomor_id = tkn.nomor_id
		WHERE ma.kontingen_id = ? AND tpa.id IS NULL
	`, kontingenID).Scan(&pairs).Error
	if err != nil || len(pairs) == 0 {
		return err
	}
	var records []TrxPendaftaranAtlet
	for _, p := range pairs {
		records = append(records, TrxPendaftaranAtlet{
			AtletID: p.AtletID,
			CaborID: p.CaborID,
			NomorID: p.NomorID,
		})
	}
	return r.db.Create(&records).Error
}

// BulkInsertTrxPelatihs insert semua pelatih kontingen ke trx_pendaftaran_pelatih
func (r *Repository) BulkInsertTrxPelatihs(kontingenID uint) error {
	type PelatihCabor struct {
		PelatihID uint
		CaborID   uint
	}
	var pairs []PelatihCabor
	err := r.db.Raw(`
		SELECT mp.id AS pelatih_id, tkc.cabor_id
		FROM master_pelatih mp
		JOIN trx_kontingen_cabor tkc ON tkc.kontingen_id = mp.kontingen_id
		LEFT JOIN trx_pendaftaran_pelatih tpp ON tpp.pelatih_id = mp.id AND tpp.cabor_id = tkc.cabor_id
		WHERE mp.kontingen_id = ? AND tpp.id IS NULL
	`, kontingenID).Scan(&pairs).Error
	if err != nil || len(pairs) == 0 {
		return err
	}
	var records []TrxPendaftaranPelatih
	for _, p := range pairs {
		records = append(records, TrxPendaftaranPelatih{
			PelatihID: p.PelatihID,
			CaborID:   p.CaborID,
		})
	}
	return r.db.Create(&records).Error
}

// BulkInsertTrxOfficials insert semua official kontingen ke trx_pendaftaran_official
func (r *Repository) BulkInsertTrxOfficials(kontingenID uint) error {
	type OfficialRow struct {
		OfficialID uint
	}
	var rows []OfficialRow
	err := r.db.Raw(`
		SELECT mo.id AS official_id
		FROM master_official mo
		LEFT JOIN trx_pendaftaran_official tpo ON tpo.official_id = mo.id
		WHERE mo.kontingen_id = ? AND tpo.id IS NULL
	`, kontingenID).Scan(&rows).Error
	if err != nil || len(rows) == 0 {
		return err
	}
	var records []TrxPendaftaranOfficial
	for _, r2 := range rows {
		records = append(records, TrxPendaftaranOfficial{OfficialID: r2.OfficialID})
	}
	return r.db.Create(&records).Error
}

// ===== STATISTIK ATLET =====

// ===== STATISTIK ATLET =====

// GetStatistikAtlet hitung total seluruh atlet dari semua kontingen.
func (r *Repository) GetStatistikAtlet() (map[string]interface{}, error) {
	var total int64
	if err := r.db.Table("master_atlet").Count(&total).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"total_atlet": total,
	}, nil
}

// ===== REFERENCE DATA — cabor & nomor yang dipilih kontingen =====

// CaborTerpilih adalah cabor yang dipilih kontingen di tahap 1
type CaborTerpilih struct {
	CaborID   uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NamaCabor string `gorm:"column:nama_cabor" json:"nama_cabor"`
}

// NomorTerdaftar adalah nomor yang dicentang kontingen di tahap 2
type NomorTerdaftar struct {
	NomorID      uint   `gorm:"column:nomor_id" json:"nomor_id"`
	CaborID      uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NamaCabor    string `gorm:"column:nama_cabor" json:"nama_cabor"`
	NamaNomor    string `gorm:"column:nama_nomor" json:"nama_nomor"`
	JenisKelamin string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Tipe         string `gorm:"column:tipe" json:"tipe"`
}

// GetCaborTerpilih ambil cabor yang sudah dipilih kontingen di tahap 1
func (r *Repository) GetCaborTerpilih(kontingenID uint) ([]CaborTerpilih, error) {
	var result []CaborTerpilih
	err := r.db.Table("trx_kontingen_cabor tkc").
		Select("tkc.cabor_id, mc.nama AS nama_cabor").
		Joins("JOIN master_cabor mc ON mc.id = tkc.cabor_id").
		Where("tkc.kontingen_id = ?", kontingenID).
		Order("mc.nama ASC").
		Scan(&result).Error
	return result, err
}

// GetNomorTerdaftar ambil nomor pertandingan yang dicentang kontingen di tahap 2
func (r *Repository) GetNomorTerdaftar(kontingenID uint) ([]NomorTerdaftar, error) {
	var result []NomorTerdaftar
	err := r.db.Table("trx_kontingen_nomor tkn").
		Select(`
			tkn.nomor_id,
			mn.cabor_id,
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
