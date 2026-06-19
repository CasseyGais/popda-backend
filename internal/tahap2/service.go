package tahap2

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetKontingenIDByTerritory untuk superadmin resolve territory → kontingen
func (s *Service) GetKontingenIDByTerritory(territoryID uint) (uint, error) {
	return s.repo.GetKontingenIDByTerritory(territoryID)
}

// GetData ambil status tahap2 + daftar nomor dari cabor yang dipilih di tahap 1
// Response menyertakan kontingen_id dan territory_id agar frontend bisa verifikasi
func (s *Service) GetData(kontingenID uint) (map[string]interface{}, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, errors.New("kontingen tidak ditemukan")
	}

	// Cek tahap 1 sudah submit
	if kontingen.Tahap1Status != "SUBMITTED" {
		return nil, errors.New("tahap 1 belum disubmit")
	}

	nomorList, err := s.repo.GetNomorByCabor(kontingenID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"kontingen_id":        kontingen.ID,
		"territory_id":            kontingen.TerritoryID,
		"nama_kontingen":          kontingen.NamaKontingen,
		"tahap2_status":           kontingen.Tahap2Status,
		"tahap2_submitted_at":     kontingen.Tahap2SubmittedAt,
		"tahap2_validasi_status":  kontingen.Tahap2ValidasiStatus,
		"tahap2_validasi_catatan": kontingen.Tahap2ValidasiCatatan,
		"nomor_list":              nomorList,
	}, nil
}

// DaftarNomor tambah satu nomor ke pendaftaran kontingen
func (s *Service) DaftarNomor(kontingenID, nomorID uint) error {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return errors.New("kontingen tidak ditemukan")
	}
	if kontingen.Tahap2Status == "SUBMITTED" {
		return errors.New("tahap 2 sudah disubmit, tidak dapat diubah")
	}

	// Validasi: nomor harus dari cabor yang dipilih di tahap 1
	valid, err := s.repo.IsNomorDariCaborKontingen(kontingenID, nomorID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("nomor tidak termasuk dalam cabor yang didaftarkan di tahap 1")
	}

	return s.repo.DaftarNomor(kontingenID, nomorID)
}

// BatalNomor hapus satu nomor dari pendaftaran
func (s *Service) BatalNomor(kontingenID, nomorID uint) error {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return errors.New("kontingen tidak ditemukan")
	}
	if kontingen.Tahap2Status == "SUBMITTED" {
		return errors.New("tahap 2 sudah disubmit, tidak dapat diubah")
	}

	return s.repo.BatalNomor(kontingenID, nomorID)
}

// Submit kunci tahap 2 — ubah tahap2_status ke SUBMITTED di tabel kontingen
func (s *Service) Submit(kontingenID uint) error {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return errors.New("kontingen tidak ditemukan")
	}
	if kontingen.Tahap2Status == "SUBMITTED" {
		return errors.New("tahap 2 sudah disubmit")
	}

	// Pastikan minimal ada satu nomor yang dipilih
	terdaftar, err := s.repo.GetTerdaftar(kontingenID)
	if err != nil {
		return err
	}
	if len(terdaftar) == 0 {
		return errors.New("pilih minimal satu nomor pertandingan sebelum submit")
	}

	return s.repo.SetTahap2Submitted(kontingenID)
}

// GetExportData menyiapkan data lengkap untuk export PDF/Excel tahap 2
func (s *Service) GetExportData(kontingenID uint) (*ExportData, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, errors.New("kontingen tidak ditemukan")
	}

	// Syarat export: tahap 1 harus SUBMITTED
	if kontingen.Tahap1Status != "SUBMITTED" {
		return nil, errors.New("tahap 1 belum disubmit")
	}

	rows, err := s.repo.GetNomorTerdaftarForExport(kontingenID)
	if err != nil {
		return nil, err
	}

	return &ExportData{
		Kontingen: kontingen,
		NomorList: rows,
	}, nil
}
