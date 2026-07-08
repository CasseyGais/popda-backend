package tahap1

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

// GetKontingenIDByTerritory cari kontingen_id dari territory_id.
// Dipakai oleh handler untuk resolve kontingen superadmin.
func (s *Service) GetKontingenIDByTerritory(territoryID uint) (uint, error) {
	return s.repo.GetKontingenIDByTerritory(territoryID)
}

// GetData ambil status tahap1 + daftar cabor yang sudah dipilih kontingen
// Response menyertakan kontingen_id dan territory_id agar frontend bisa verifikasi
func (s *Service) GetData(kontingenID uint) (map[string]interface{}, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, errors.New("kontingen tidak ditemukan")
	}

	caborList, err := s.repo.GetCabor(kontingenID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"kontingen_id":            kontingen.ID,
		"territory_id":            kontingen.TerritoryID,
		"nama_kontingen":          kontingen.NamaKontingen,
		"tahap1_status":           kontingen.Tahap1Status,
		"tahap1_submitted_at":     kontingen.Tahap1SubmittedAt,
		"tahap1_validasi_status":  kontingen.Tahap1ValidasiStatus,
		"tahap1_validasi_catatan": kontingen.Tahap1ValidasiCatatan,
		"cabor_list":              caborList,
	}, nil
}

// UpsertCabor tambah atau update satu cabor di tahap 1.
// Validasi kuota dilakukan di service agar berlaku baik saat insert maupun update.
// (DB trigger hanya aktif saat INSERT)
func (s *Service) UpsertCabor(kontingenID, caborID uint, putra, putri, pelatih int) error {
	// Cek status — tidak boleh edit setelah SUBMITTED
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return errors.New("kontingen tidak ditemukan")
	}
	if kontingen.Tahap1Status == "SUBMITTED" {
		return errors.New("Tahap 1 sudah disubmit, tidak dapat diubah")
	}

	// Validasi kuota dari master_cabor (berlaku untuk INSERT dan UPDATE)
	maxPutra, maxPutri, maxPelatih, err := s.repo.GetMasterCaborKuota(caborID)
	if err != nil {
		return errors.New("cabor tidak ditemukan")
	}
	if putra > maxPutra {
		return fmt.Errorf("jumlah atlet putra melebihi kuota (%d)", maxPutra)
	}
	if putri > maxPutri {
		return fmt.Errorf("jumlah atlet putri melebihi kuota (%d)", maxPutri)
	}
	if pelatih > maxPelatih {
		return fmt.Errorf("jumlah pelatih melebihi kuota (%d)", maxPelatih)
	}

	totalAtlet := putra + putri
	totalPersonel := totalAtlet + pelatih

	trx := &TrxKontingenCabor{
		KontingenID:   kontingenID,
		CaborID:       caborID,
		Putra:         putra,
		Putri:         putri,
		Pelatih:       pelatih,
		TotalAtlet:    totalAtlet,
		TotalPersonel: totalPersonel,
	}

	return s.repo.UpsertCabor(trx)
}

// DeleteCabor hapus satu cabor dari daftar tahap 1
func (s *Service) DeleteCabor(kontingenID, caborID uint) error {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return errors.New("kontingen tidak ditemukan")
	}
	if kontingen.Tahap1Status == "SUBMITTED" {
		return errors.New("Tahap 1 sudah disubmit, tidak dapat diubah")
	}

	return s.repo.DeleteCabor(kontingenID, caborID)
}

// Submit kunci tahap 1 — ubah tahap1_status ke SUBMITTED di tabel kontingen
func (s *Service) Submit(kontingenID uint) error {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return errors.New("kontingen tidak ditemukan")
	}
	if kontingen.Tahap1Status == "SUBMITTED" {
		return errors.New("Tahap 1 sudah disubmit")
	}

	// Pastikan minimal ada satu cabor
	caborList, err := s.repo.GetCabor(kontingenID)
	if err != nil {
		return err
	}
	if len(caborList) == 0 {
		return errors.New("pilih minimal satu cabang olahraga sebelum submit")
	}

	return s.repo.SetTahap1Submitted(kontingenID)
}

// ResetTahap1 kembalikan tahap1_status ke DRAFT agar kontingen bisa edit & submit ulang.
// Hanya superadmin yang boleh memanggil ini (guard di handler).
// Juga reset validasi status ke NULL karena data akan diubah.
func (s *Service) ResetTahap1(kontingenID uint) error {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return errors.New("kontingen tidak ditemukan")
	}
	if kontingen.Tahap1Status != "SUBMITTED" {
		return errors.New("Tahap 1 belum disubmit, tidak perlu di-reset")
	}
	return s.repo.ResetTahap1(kontingenID)
}
func (s *Service) GetExportData(kontingenID uint) (*ExportData, error) {
	kontingen, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, errors.New("kontingen tidak ditemukan")
	}

	rows, err := s.repo.GetCaborWithNama(kontingenID)
	if err != nil {
		return nil, err
	}

	var totalPutra, totalPutri, totalPelatih, totalAtlet, totalPersonel int
	for _, r := range rows {
		totalPutra += r.Putra
		totalPutri += r.Putri
		totalPelatih += r.Pelatih
		totalAtlet += r.TotalAtlet
		totalPersonel += r.TotalPersonel
	}

	return &ExportData{
		Kontingen:     kontingen,
		CaborList:     rows,
		TotalPutra:    totalPutra,
		TotalPutri:    totalPutri,
		TotalPelatih:  totalPelatih,
		TotalAtlet:    totalAtlet,
		TotalPersonel: totalPersonel,
	}, nil
}
