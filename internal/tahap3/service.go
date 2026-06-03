package tahap3

import (
	"errors"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetKontingenIDByTerritory untuk superadmin
func (s *Service) GetKontingenIDByTerritory(territoryID uint) (uint, error) {
	return s.repo.GetKontingenIDByTerritory(territoryID)
}

// GetData ambil semua data tahap 3 milik kontingen
func (s *Service) GetData(kontingenID uint) (map[string]interface{}, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, errors.New("kontingen tidak ditemukan")
	}

	atlets, _ := s.repo.GetAtlets(kontingenID)
	pelatihs, _ := s.repo.GetPelatihs(kontingenID)
	officials, _ := s.repo.GetOfficials(kontingenID)
	trxAtlets, _ := s.repo.GetTrxAtlets(kontingenID)
	trxPelatihs, _ := s.repo.GetTrxPelatihs(kontingenID)
	trxOfficials, _ := s.repo.GetTrxOfficials(kontingenID)

	return map[string]interface{}{
		"kontingen":     kontingen,
		"atlets":        atlets,
		"pelatihs":      pelatihs,
		"officials":     officials,
		"trx_atlets":    trxAtlets,
		"trx_pelatihs":  trxPelatihs,
		"trx_officials": trxOfficials,
	}, nil
}

// ===== ATLET =====

func (s *Service) GetAtlets(kontingenID uint) ([]MasterAtlet, error) {
	return s.repo.GetAtlets(kontingenID)
}

func (s *Service) GetAtletByID(id uint) (*MasterAtlet, error) {
	atlet, err := s.repo.GetAtletByID(id)
	if err != nil {
		return nil, errors.New("atlet tidak ditemukan")
	}
	return atlet, nil
}

func (s *Service) CreateAtlet(kontingenID uint, req *CreateAtletRequest) (*MasterAtlet, error) {
	atlet := &MasterAtlet{
		KontingenID:              kontingenID,
		NamaLengkap:              req.NamaLengkap,
		JenisKelamin:             req.JenisKelamin,
		TanggalLahir:             req.TanggalLahir,
		TempatLahir:              req.TempatLahir,
		NISN:                     req.NISN,
		NIS:                      req.NIS,
		Sekolah:                  req.Sekolah,
		KelasJurusan:             req.KelasJurusan,
		Alamat:                   req.Alamat,
		KabupatenKota:            req.KabupatenKota,
		NoHP:                     req.NoHP,
		NamaOrtuWali:             req.NamaOrtuWali,
		PrestasiSebelumnya:       req.PrestasiSebelumnya,
		Catatan:                  req.Catatan,
		Status:                   "draft",
	}
	if err := s.repo.CreateAtlet(atlet); err != nil {
		return nil, fmt.Errorf("gagal membuat atlet: %w", err)
	}
	return atlet, nil
}

func (s *Service) UpdateAtlet(id uint, req *UpdateAtletRequest) (*MasterAtlet, error) {
	atlet, err := s.repo.GetAtletByID(id)
	if err != nil {
		return nil, errors.New("atlet tidak ditemukan")
	}
	if req.NamaLengkap != "" { atlet.NamaLengkap = req.NamaLengkap }
	if req.JenisKelamin != "" { atlet.JenisKelamin = req.JenisKelamin }
	if req.TanggalLahir != "" { atlet.TanggalLahir = req.TanggalLahir }
	if req.TempatLahir != "" { atlet.TempatLahir = req.TempatLahir }
	if req.NISN != "" { atlet.NISN = req.NISN }
	if req.NIS != "" { atlet.NIS = req.NIS }
	if req.Sekolah != "" { atlet.Sekolah = req.Sekolah }
	if req.KelasJurusan != "" { atlet.KelasJurusan = req.KelasJurusan }
	if req.Alamat != "" { atlet.Alamat = req.Alamat }
	if req.KabupatenKota != "" { atlet.KabupatenKota = req.KabupatenKota }
	if req.NoHP != "" { atlet.NoHP = req.NoHP }
	if req.NamaOrtuWali != "" { atlet.NamaOrtuWali = req.NamaOrtuWali }
	if req.PrestasiSebelumnya != "" { atlet.PrestasiSebelumnya = req.PrestasiSebelumnya }
	if req.Catatan != "" { atlet.Catatan = req.Catatan }
	if err := s.repo.UpdateAtlet(atlet); err != nil {
		return nil, errors.New("gagal mengupdate atlet")
	}
	return atlet, nil
}

func (s *Service) DeleteAtlet(id uint) error {
	if _, err := s.repo.GetAtletByID(id); err != nil {
		return errors.New("atlet tidak ditemukan")
	}
	return s.repo.DeleteAtlet(id)
}

// ===== PELATIH =====

func (s *Service) GetPelatihs(kontingenID uint) ([]MasterPelatih, error) {
	return s.repo.GetPelatihs(kontingenID)
}

func (s *Service) GetPelatihByID(id uint) (*MasterPelatih, error) {
	p, err := s.repo.GetPelatihByID(id)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}
	return p, nil
}

func (s *Service) CreatePelatih(kontingenID uint, req *CreatePelatihRequest) (*MasterPelatih, error) {
	pelatih := &MasterPelatih{
		KontingenID:    kontingenID,
		NamaLengkap:    req.NamaLengkap,
		JenisKelamin:   req.JenisKelamin,
		TanggalLahir:   req.TanggalLahir,
		TempatLahir:    req.TempatLahir,
		NIK:            req.NIK,
		SekolahAsal:    req.SekolahAsal,
		Profesi:        req.Profesi,
		Jabatan:        req.Jabatan,
		Alamat:         req.Alamat,
		KabupatenKota:  req.KabupatenKota,
		NoHP:           req.NoHP,
		Email:          req.Email,
		NamaIstriSuami: req.NamaIstriSuami,
		PrestasiSebelumnya: req.PrestasiSebelumnya,
		Catatan:        req.Catatan,
		Status:         "draft",
	}
	if err := s.repo.CreatePelatih(pelatih); err != nil {
		return nil, fmt.Errorf("gagal membuat pelatih: %w", err)
	}
	return pelatih, nil
}

func (s *Service) UpdatePelatih(id uint, req *UpdatePelatihRequest) (*MasterPelatih, error) {
	pelatih, err := s.repo.GetPelatihByID(id)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}
	if req.NamaLengkap != "" { pelatih.NamaLengkap = req.NamaLengkap }
	if req.JenisKelamin != "" { pelatih.JenisKelamin = req.JenisKelamin }
	if req.TanggalLahir != "" { pelatih.TanggalLahir = req.TanggalLahir }
	if req.TempatLahir != "" { pelatih.TempatLahir = req.TempatLahir }
	if req.NIK != "" { pelatih.NIK = req.NIK }
	if req.SekolahAsal != "" { pelatih.SekolahAsal = req.SekolahAsal }
	if req.Profesi != "" { pelatih.Profesi = req.Profesi }
	if req.Jabatan != "" { pelatih.Jabatan = req.Jabatan }
	if req.Alamat != "" { pelatih.Alamat = req.Alamat }
	if req.KabupatenKota != "" { pelatih.KabupatenKota = req.KabupatenKota }
	if req.NoHP != "" { pelatih.NoHP = req.NoHP }
	if req.Email != "" { pelatih.Email = req.Email }
	if req.NamaIstriSuami != "" { pelatih.NamaIstriSuami = req.NamaIstriSuami }
	if req.PrestasiSebelumnya != "" { pelatih.PrestasiSebelumnya = req.PrestasiSebelumnya }
	if req.Catatan != "" { pelatih.Catatan = req.Catatan }
	if err := s.repo.UpdatePelatih(pelatih); err != nil {
		return nil, errors.New("gagal mengupdate pelatih")
	}
	return pelatih, nil
}

func (s *Service) DeletePelatih(id uint) error {
	if _, err := s.repo.GetPelatihByID(id); err != nil {
		return errors.New("pelatih tidak ditemukan")
	}
	return s.repo.DeletePelatih(id)
}

// ===== OFFICIAL =====

func (s *Service) GetOfficials(kontingenID uint) ([]MasterOfficial, error) {
	return s.repo.GetOfficials(kontingenID)
}

func (s *Service) GetOfficialByID(id uint) (*MasterOfficial, error) {
	o, err := s.repo.GetOfficialByID(id)
	if err != nil {
		return nil, errors.New("official tidak ditemukan")
	}
	return o, nil
}

func (s *Service) CreateOfficial(kontingenID uint, req *CreateOfficialRequest) (*MasterOfficial, error) {
	official := &MasterOfficial{
		KontingenID:   kontingenID,
		NamaLengkap:   req.NamaLengkap,
		JenisKelamin:  req.JenisKelamin,
		TanggalLahir:  req.TanggalLahir,
		TempatLahir:   req.TempatLahir,
		NIK:           req.NIK,
		SekolahAsal:   req.SekolahAsal,
		Jabatan:       req.Jabatan,
		Alamat:        req.Alamat,
		KabupatenKota: req.KabupatenKota,
		NoHP:          req.NoHP,
		Email:         req.Email,
		Catatan:       req.Catatan,
		Status:        "draft",
	}
	if err := s.repo.CreateOfficial(official); err != nil {
		return nil, fmt.Errorf("gagal membuat official: %w", err)
	}
	return official, nil
}

func (s *Service) UpdateOfficial(id uint, req *UpdateOfficialRequest) (*MasterOfficial, error) {
	official, err := s.repo.GetOfficialByID(id)
	if err != nil {
		return nil, errors.New("official tidak ditemukan")
	}
	if req.NamaLengkap != "" { official.NamaLengkap = req.NamaLengkap }
	if req.JenisKelamin != "" { official.JenisKelamin = req.JenisKelamin }
	if req.TanggalLahir != "" { official.TanggalLahir = req.TanggalLahir }
	if req.TempatLahir != "" { official.TempatLahir = req.TempatLahir }
	if req.NIK != "" { official.NIK = req.NIK }
	if req.SekolahAsal != "" { official.SekolahAsal = req.SekolahAsal }
	if req.Jabatan != "" { official.Jabatan = req.Jabatan }
	if req.Alamat != "" { official.Alamat = req.Alamat }
	if req.KabupatenKota != "" { official.KabupatenKota = req.KabupatenKota }
	if req.NoHP != "" { official.NoHP = req.NoHP }
	if req.Email != "" { official.Email = req.Email }
	if req.Catatan != "" { official.Catatan = req.Catatan }
	if err := s.repo.UpdateOfficial(official); err != nil {
		return nil, errors.New("gagal mengupdate official")
	}
	return official, nil
}

func (s *Service) DeleteOfficial(id uint) error {
	if _, err := s.repo.GetOfficialByID(id); err != nil {
		return errors.New("official tidak ditemukan")
	}
	return s.repo.DeleteOfficial(id)
}

// ===== TRX PENDAFTARAN =====

func (s *Service) CreateTrxAtlet(req *CreateTrxAtletRequest) (*TrxPendaftaranAtlet, error) {
	trx := &TrxPendaftaranAtlet{
		AtletID: req.AtletID,
		CaborID: req.CaborID,
		NomorID: req.NomorID,
	}
	if err := s.repo.CreateTrxAtlet(trx); err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan atlet: %w", err)
	}
	return trx, nil
}

func (s *Service) DeleteTrxAtlet(id uint) error {
	return s.repo.DeleteTrxAtlet(id)
}

func (s *Service) CreateTrxPelatih(req *CreateTrxPelatihRequest) (*TrxPendaftaranPelatih, error) {
	trx := &TrxPendaftaranPelatih{
		PelatihID: req.PelatihID,
		CaborID:   req.CaborID,
	}
	if err := s.repo.CreateTrxPelatih(trx); err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan pelatih: %w", err)
	}
	return trx, nil
}

func (s *Service) DeleteTrxPelatih(id uint) error {
	return s.repo.DeleteTrxPelatih(id)
}

func (s *Service) CreateTrxOfficial(req *CreateTrxOfficialRequest) (*TrxPendaftaranOfficial, error) {
	trx := &TrxPendaftaranOfficial{
		OfficialID: req.OfficialID,
	}
	if err := s.repo.CreateTrxOfficial(trx); err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan official: %w", err)
	}
	return trx, nil
}

func (s *Service) DeleteTrxOfficial(id uint) error {
	return s.repo.DeleteTrxOfficial(id)
}

// ===== SUBMIT TAHAP 3 =====

// Submit kunci tahap 3:
// 1. Cek minimal ada 1 atlet
// 2. Bulk insert semua atlet → trx_pendaftaran_atlet (yang belum ada)
// 3. Bulk insert semua pelatih → trx_pendaftaran_pelatih (yang belum ada)
// 4. Bulk insert semua official → trx_pendaftaran_official (yang belum ada)
// 5. Update status master menjadi "terdaftar"
// 6. Set tahap3_status = SUBMITTED di tabel kontingen
func (s *Service) Submit(kontingenID uint) error {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return errors.New("kontingen tidak ditemukan")
	}
	if kontingen.Tahap3Status == "SUBMITTED" {
		return errors.New("tahap 3 sudah disubmit")
	}

	// Validasi minimal ada 1 atlet
	atlets, err := s.repo.GetAtlets(kontingenID)
	if err != nil {
		return err
	}
	if len(atlets) == 0 {
		return errors.New("minimal harus ada 1 atlet sebelum submit")
	}

	// Bulk insert trx — skip yang sudah ada
	if err := s.repo.BulkInsertTrxAtlets(kontingenID); err != nil {
		return fmt.Errorf("gagal mendaftarkan atlet: %w", err)
	}
	if err := s.repo.BulkInsertTrxPelatihs(kontingenID); err != nil {
		return fmt.Errorf("gagal mendaftarkan pelatih: %w", err)
	}
	if err := s.repo.BulkInsertTrxOfficials(kontingenID); err != nil {
		return fmt.Errorf("gagal mendaftarkan official: %w", err)
	}

	// Set tahap3_status = SUBMITTED
	return s.repo.SetTahap3Submitted(kontingenID)
}

// ===== EXPORT =====

// GetExportData menyiapkan semua data untuk export tahap 3 (atlet + pelatih + official)
func (s *Service) GetExportData(kontingenID uint) (*ExportData, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, errors.New("kontingen tidak ditemukan")
	}

	atlets, err := s.repo.GetAtlets(kontingenID)
	if err != nil {
		return nil, err
	}
	pelatihs, err := s.repo.GetPelatihs(kontingenID)
	if err != nil {
		return nil, err
	}
	officials, err := s.repo.GetOfficials(kontingenID)
	if err != nil {
		return nil, err
	}

	return &ExportData{
		Kontingen: kontingen,
		Atlets:    atlets,
		Pelatihs:  pelatihs,
		Officials: officials,
	}, nil
}

// GetExportAtlets menyiapkan data atlet saja untuk export
func (s *Service) GetExportAtlets(kontingenID uint) (*Kontingen, []MasterAtlet, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, nil, errors.New("kontingen tidak ditemukan")
	}
	atlets, err := s.repo.GetAtlets(kontingenID)
	if err != nil {
		return nil, nil, err
	}
	return kontingen, atlets, nil
}

// GetExportPelatihs menyiapkan data pelatih saja untuk export
func (s *Service) GetExportPelatihs(kontingenID uint) (*Kontingen, []MasterPelatih, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, nil, errors.New("kontingen tidak ditemukan")
	}
	pelatihs, err := s.repo.GetPelatihs(kontingenID)
	if err != nil {
		return nil, nil, err
	}
	return kontingen, pelatihs, nil
}

// GetExportOfficials menyiapkan data official saja untuk export
func (s *Service) GetExportOfficials(kontingenID uint) (*Kontingen, []MasterOfficial, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, nil, errors.New("kontingen tidak ditemukan")
	}
	officials, err := s.repo.GetOfficials(kontingenID)
	if err != nil {
		return nil, nil, err
	}
	return kontingen, officials, nil
}
