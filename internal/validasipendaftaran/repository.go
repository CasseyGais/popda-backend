package validasipendaftaran

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetAll ambil semua kontingen beserta status validasi ketiga tahap.
// Filter opsional: status (PENDING/VALID/REVISI), tahap (1/2/3), territory_id.
func (r *Repository) GetAll(filterStatus string, filterTahap int, filterTerritoryID uint) ([]KontingenValidasi, error) {
	q := r.db.Table("kontingen k").
		Select(`
			k.id                       AS kontingen_id,
			k.territory_id,
			k.nama_kontingen,
			k.tahap1_status,
			k.tahap1_submitted_at,
			k.tahap1_validasi_status,
			k.tahap1_validasi_catatan,
			k.tahap1_validasi_at,
			k.tahap2_status,
			k.tahap2_submitted_at,
			k.tahap2_validasi_status,
			k.tahap2_validasi_catatan,
			k.tahap2_validasi_at,
			k.tahap3_status,
			k.tahap3_submitted_at,
			k.tahap3_validasi_status,
			k.tahap3_validasi_catatan,
			k.tahap3_validasi_at
		`).
		Order("k.nama_kontingen ASC")

	if filterTerritoryID > 0 {
		q = q.Where("k.territory_id = ?", filterTerritoryID)
	}

	// Filter by status: cari di tahap manapun atau tahap tertentu
	if filterStatus != "" {
		upperStatus := strings.ToUpper(filterStatus)
		if filterTahap > 0 {
			switch filterTahap {
			case 1:
				q = q.Where("k.tahap1_validasi_status = ?", upperStatus)
			case 2:
				q = q.Where("k.tahap2_validasi_status = ?", upperStatus)
			case 3:
				q = q.Where("k.tahap3_validasi_status = ?", upperStatus)
			}
		} else {
			q = q.Where(
				"k.tahap1_validasi_status = ? OR k.tahap2_validasi_status = ? OR k.tahap3_validasi_status = ?",
				upperStatus, upperStatus, upperStatus,
			)
		}
	} else if filterTahap > 0 {
		// Filter tahap saja (tanpa filter status) — hanya yang sudah SUBMITTED
		switch filterTahap {
		case 1:
			q = q.Where("k.tahap1_status = 'SUBMITTED'")
		case 2:
			q = q.Where("k.tahap2_status = 'SUBMITTED'")
		case 3:
			q = q.Where("k.tahap3_status = 'SUBMITTED'")
		}
	}

	var result []KontingenValidasi
	err := q.Scan(&result).Error
	return result, err
}

// GetKontingen ambil satu kontingen untuk keperluan cek status & nama.
func (r *Repository) GetKontingen(kontingenID uint) (*KontingenRow, error) {
	var k KontingenRow
	err := r.db.Table("kontingen").
		Select("id, nama_kontingen, tahap1_status, tahap2_status, tahap3_status").
		Where("id = ?", kontingenID).
		First(&k).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// GetStatusByKontingen ambil hanya field validasi untuk widget dashboard.
func (r *Repository) GetStatusByKontingen(kontingenID uint) (*KontingenStatusRow, error) {
	var result KontingenStatusRow
	err := r.db.Table("kontingen").
		Select(`id, nama_kontingen,
			tahap1_validasi_status, tahap1_validasi_catatan,
			tahap2_validasi_status, tahap2_validasi_catatan,
			tahap3_validasi_status, tahap3_validasi_catatan`).
		Where("id = ?", kontingenID).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SetValidasi simpan status validasi (VALID/REVISI) untuk satu tahap kontingen.
func (r *Repository) SetValidasi(kontingenID uint, tahap int, status string, catatan *string) error {
	now := time.Now()
	var updates map[string]interface{}

	switch tahap {
	case 1:
		updates = map[string]interface{}{
			"tahap1_validasi_status":  status,
			"tahap1_validasi_catatan": catatan,
			"tahap1_validasi_at":      now,
		}
	case 2:
		updates = map[string]interface{}{
			"tahap2_validasi_status":  status,
			"tahap2_validasi_catatan": catatan,
			"tahap2_validasi_at":      now,
		}
	case 3:
		updates = map[string]interface{}{
			"tahap3_validasi_status":  status,
			"tahap3_validasi_catatan": catatan,
			"tahap3_validasi_at":      now,
		}
	}

	return r.db.Table("kontingen").
		Where("id = ?", kontingenID).
		Updates(updates).Error
}

// GetKontingenIDByTerritory cari kontingen_id berdasarkan territory_id (untuk superadmin).
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

// ===== REKAP PENDAFTARAN =====

// RekapCabor dipakai untuk cabor_terpilih di rekap
type RekapCabor struct {
	CaborID   uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NamaCabor string `gorm:"column:nama_cabor" json:"nama_cabor"`
	Putra     int    `gorm:"column:putra" json:"putra"`
	Putri     int    `gorm:"column:putri" json:"putri"`
	Pelatih   int    `gorm:"column:pelatih" json:"pelatih"`
	TotalAtlet int   `gorm:"column:total_atlet" json:"total_atlet"`
}

// RekapNomor dipakai untuk nomor_terdaftar di rekap
type RekapNomor struct {
	NomorID      uint   `gorm:"column:nomor_id" json:"nomor_id"`
	CaborID      uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NamaCabor    string `gorm:"column:nama_cabor" json:"nama_cabor"`
	NamaNomor    string `gorm:"column:nama_nomor" json:"nama_nomor"`
	JenisKelamin string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Tipe         string `gorm:"column:tipe" json:"tipe"`
}

// RekapAtlet dipakai untuk list atlet di rekap
type RekapAtlet struct {
	ID            uint   `gorm:"column:id" json:"id"`
	NamaLengkap   string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	JenisKelamin  string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	TanggalLahir  string `gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	NISN          string `gorm:"column:nisn" json:"nisn"`
	Sekolah       string `gorm:"column:sekolah" json:"sekolah"`
	KelasJurusan  string `gorm:"column:kelas_jurusan" json:"kelas_jurusan"`
	KabupatenKota string `gorm:"column:kabupaten_kota" json:"kabupaten_kota"`
	NoHP          string `gorm:"column:no_hp" json:"no_hp"`
	Status        string `gorm:"column:status" json:"status"`
}

// RekapAtletTrx dipakai untuk trx masing-masing atlet di rekap
type RekapAtletTrx struct {
	AtletID      uint   `gorm:"column:atlet_id" json:"atlet_id"`
	CaborID      uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NamaCabor    string `gorm:"column:nama_cabor" json:"nama_cabor"`
	NomorID      uint   `gorm:"column:nomor_id" json:"nomor_id"`
	NamaNomor    string `gorm:"column:nama_nomor" json:"nama_nomor"`
	JenisKelamin string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Tipe         string `gorm:"column:tipe" json:"tipe"`
}

// RekapPelatih dipakai untuk list pelatih di rekap
type RekapPelatih struct {
	ID            uint   `gorm:"column:id" json:"id"`
	NamaLengkap   string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	JenisKelamin  string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Jabatan       string `gorm:"column:jabatan" json:"jabatan"`
	NoHP          string `gorm:"column:no_hp" json:"no_hp"`
	Status        string `gorm:"column:status" json:"status"`
}

// RekapPelatihTrx dipakai untuk trx masing-masing pelatih di rekap
type RekapPelatihTrx struct {
	PelatihID uint   `gorm:"column:pelatih_id" json:"pelatih_id"`
	CaborID   uint   `gorm:"column:cabor_id" json:"cabor_id"`
	NamaCabor string `gorm:"column:nama_cabor" json:"nama_cabor"`
}

// RekapOfficial dipakai untuk list official di rekap
type RekapOfficial struct {
	ID            uint   `gorm:"column:id" json:"id"`
	NamaLengkap   string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	JenisKelamin  string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Jabatan       string `gorm:"column:jabatan" json:"jabatan"`
	NoHP          string `gorm:"column:no_hp" json:"no_hp"`
	Status        string `gorm:"column:status" json:"status"`
}

// GetRekap ambil semua data pendaftaran kontingen untuk halaman Rekap Pendaftaran.
func (r *Repository) GetRekap(kontingenID uint) (map[string]interface{}, error) {
	// 1. Data kontingen + validasi
	var k struct {
		ID                    uint    `gorm:"column:id"`
		TerritoryID           uint    `gorm:"column:territory_id"`
		NamaKontingen         string  `gorm:"column:nama_kontingen"`
		Tahap1Status          string  `gorm:"column:tahap1_status"`
		Tahap1ValidasiStatus  *string `gorm:"column:tahap1_validasi_status"`
		Tahap1ValidasiCatatan *string `gorm:"column:tahap1_validasi_catatan"`
		Tahap2Status          string  `gorm:"column:tahap2_status"`
		Tahap2ValidasiStatus  *string `gorm:"column:tahap2_validasi_status"`
		Tahap2ValidasiCatatan *string `gorm:"column:tahap2_validasi_catatan"`
		Tahap3Status          string  `gorm:"column:tahap3_status"`
		Tahap3ValidasiStatus  *string `gorm:"column:tahap3_validasi_status"`
		Tahap3ValidasiCatatan *string `gorm:"column:tahap3_validasi_catatan"`
	}
	if err := r.db.Table("kontingen").Where("id = ?", kontingenID).First(&k).Error; err != nil {
		return nil, err
	}

	// 2. Cabor terpilih (dari tahap 1)
	var caborList []RekapCabor
	r.db.Table("trx_kontingen_cabor tkc").
		Select("tkc.cabor_id, mc.nama AS nama_cabor, tkc.putra, tkc.putri, tkc.pelatih, tkc.total_atlet").
		Joins("JOIN master_cabor mc ON mc.id = tkc.cabor_id").
		Where("tkc.kontingen_id = ?", kontingenID).
		Order("mc.nama ASC").
		Scan(&caborList)

	// 3. Nomor terdaftar (dari tahap 2)
	var nomorList []RekapNomor
	r.db.Table("trx_kontingen_nomor tkn").
		Select("tkn.nomor_id, mn.cabor_id, mc.nama AS nama_cabor, mn.nama AS nama_nomor, mn.jenis_kelamin, mn.tipe").
		Joins("JOIN master_nomor mn ON mn.id = tkn.nomor_id").
		Joins("JOIN master_cabor mc ON mc.id = mn.cabor_id").
		Where("tkn.kontingen_id = ?", kontingenID).
		Order("mc.nama ASC, mn.nama ASC").
		Scan(&nomorList)

	// 4. Atlet
	var atletList []RekapAtlet
	r.db.Table("master_atlet").
		Select("id, nama_lengkap, jenis_kelamin, tanggal_lahir, nisn, sekolah, kelas_jurusan, kabupaten_kota, no_hp, status").
		Where("kontingen_id = ?", kontingenID).
		Order("nama_lengkap ASC").
		Scan(&atletList)

	// 4b. Trx atlet — semua sekaligus, lalu group di Go
	var atletTrxList []RekapAtletTrx
	r.db.Table("trx_pendaftaran_atlet tpa").
		Select("tpa.atlet_id, tpa.cabor_id, mc.nama AS nama_cabor, tpa.nomor_id, mn.nama AS nama_nomor, mn.jenis_kelamin, mn.tipe").
		Joins("JOIN master_nomor mn ON mn.id = tpa.nomor_id").
		Joins("JOIN master_cabor mc ON mc.id = tpa.cabor_id").
		Joins("JOIN master_atlet ma ON ma.id = tpa.atlet_id").
		Where("ma.kontingen_id = ?", kontingenID).
		Scan(&atletTrxList)

	// Build map atlet_id → []trx
	atletTrxMap := make(map[uint][]RekapAtletTrx)
	for _, t := range atletTrxList {
		atletTrxMap[t.AtletID] = append(atletTrxMap[t.AtletID], t)
	}

	// Gabungkan atlet + trx
	atletsWithTrx := make([]map[string]interface{}, 0, len(atletList))
	for _, a := range atletList {
		trx := atletTrxMap[a.ID]
		if trx == nil {
			trx = []RekapAtletTrx{}
		}
		atletsWithTrx = append(atletsWithTrx, map[string]interface{}{
			"id":             a.ID,
			"nama_lengkap":   a.NamaLengkap,
			"jenis_kelamin":  a.JenisKelamin,
			"tanggal_lahir":  a.TanggalLahir,
			"nisn":           a.NISN,
			"sekolah":        a.Sekolah,
			"kelas_jurusan":  a.KelasJurusan,
			"kabupaten_kota": a.KabupatenKota,
			"no_hp":          a.NoHP,
			"status":         a.Status,
			"trx":            trx,
		})
	}

	// 5. Pelatih
	var pelatihList []RekapPelatih
	r.db.Table("master_pelatih").
		Select("id, nama_lengkap, jenis_kelamin, jabatan, no_hp, status").
		Where("kontingen_id = ?", kontingenID).
		Order("nama_lengkap ASC").
		Scan(&pelatihList)

	var pelatihTrxList []RekapPelatihTrx
	r.db.Table("trx_pendaftaran_pelatih tpp").
		Select("tpp.pelatih_id, tpp.cabor_id, mc.nama AS nama_cabor").
		Joins("JOIN master_cabor mc ON mc.id = tpp.cabor_id").
		Joins("JOIN master_pelatih mp ON mp.id = tpp.pelatih_id").
		Where("mp.kontingen_id = ?", kontingenID).
		Scan(&pelatihTrxList)

	pelatihTrxMap := make(map[uint][]RekapPelatihTrx)
	for _, t := range pelatihTrxList {
		pelatihTrxMap[t.PelatihID] = append(pelatihTrxMap[t.PelatihID], t)
	}

	pelatihsWithTrx := make([]map[string]interface{}, 0, len(pelatihList))
	for _, p := range pelatihList {
		trx := pelatihTrxMap[p.ID]
		if trx == nil {
			trx = []RekapPelatihTrx{}
		}
		pelatihsWithTrx = append(pelatihsWithTrx, map[string]interface{}{
			"id":            p.ID,
			"nama_lengkap":  p.NamaLengkap,
			"jenis_kelamin": p.JenisKelamin,
			"jabatan":       p.Jabatan,
			"no_hp":         p.NoHP,
			"status":        p.Status,
			"trx":           trx,
		})
	}

	// 6. Official
	var officialList []RekapOfficial
	r.db.Table("master_official").
		Select("id, nama_lengkap, jenis_kelamin, jabatan, no_hp, status").
		Where("kontingen_id = ?", kontingenID).
		Order("nama_lengkap ASC").
		Scan(&officialList)

	// Cek mana official yang sudah terdaftar di trx
	type OfficialTrxID struct {
		OfficialID uint `gorm:"column:official_id"`
	}
	var officialTrxIDs []OfficialTrxID
	r.db.Table("trx_pendaftaran_official tpo").
		Select("tpo.official_id").
		Joins("JOIN master_official mo ON mo.id = tpo.official_id").
		Where("mo.kontingen_id = ?", kontingenID).
		Scan(&officialTrxIDs)

	terdaftarSet := make(map[uint]bool)
	for _, o := range officialTrxIDs {
		terdaftarSet[o.OfficialID] = true
	}

	officialsWithTrx := make([]map[string]interface{}, 0, len(officialList))
	for _, o := range officialList {
		var trx interface{}
		if terdaftarSet[o.ID] {
			trx = []map[string]interface{}{{"keterangan": "Terdaftar sebagai official kontingen"}}
		} else {
			trx = []map[string]interface{}{}
		}
		officialsWithTrx = append(officialsWithTrx, map[string]interface{}{
			"id":            o.ID,
			"nama_lengkap":  o.NamaLengkap,
			"jenis_kelamin": o.JenisKelamin,
			"jabatan":       o.Jabatan,
			"no_hp":         o.NoHP,
			"status":        o.Status,
			"trx":           trx,
		})
	}

	return map[string]interface{}{
		"kontingen_id":   k.ID,
		"territory_id":   k.TerritoryID,
		"nama_kontingen": k.NamaKontingen,
		"validasi": map[string]interface{}{
			"tahap1": map[string]interface{}{"status": k.Tahap1ValidasiStatus, "catatan": k.Tahap1ValidasiCatatan},
			"tahap2": map[string]interface{}{"status": k.Tahap2ValidasiStatus, "catatan": k.Tahap2ValidasiCatatan},
			"tahap3": map[string]interface{}{"status": k.Tahap3ValidasiStatus, "catatan": k.Tahap3ValidasiCatatan},
		},
		"cabor_terpilih":  caborList,
		"nomor_terdaftar": nomorList,
		"atlets":          atletsWithTrx,
		"pelatihs":        pelatihsWithTrx,
		"officials":       officialsWithTrx,
	}, nil
}
