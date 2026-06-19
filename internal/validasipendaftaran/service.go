package validasipendaftaran

import (
	"errors"
	"strings"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetAll list semua kontingen beserta status validasi.
// Hanya untuk superadmin — guard di handler.
func (s *Service) GetAll(filterStatus string, filterTahap int, filterTerritoryID uint) ([]map[string]interface{}, error) {
	rows, err := s.repo.GetAll(filterStatus, filterTahap, filterTerritoryID)
	if err != nil {
		return nil, err
	}

	// Transform ke format response yang diinginkan frontend
	result := make([]map[string]interface{}, 0, len(rows))
	for _, k := range rows {
		result = append(result, map[string]interface{}{
			"kontingen_id":   k.KontingenID,
			"territory_id":   k.TerritoryID,
			"nama_kontingen": k.NamaKontingen,
			"tahap1": map[string]interface{}{
				"submit_status":    k.Tahap1SubmitStatus,
				"submitted_at":     k.Tahap1SubmittedAt,
				"validasi_status":  k.Tahap1ValidasiStatus,
				"validasi_catatan": k.Tahap1ValidasiCatatan,
				"validasi_at":      k.Tahap1ValidasiAt,
			},
			"tahap2": map[string]interface{}{
				"submit_status":    k.Tahap2SubmitStatus,
				"submitted_at":     k.Tahap2SubmittedAt,
				"validasi_status":  k.Tahap2ValidasiStatus,
				"validasi_catatan": k.Tahap2ValidasiCatatan,
				"validasi_at":      k.Tahap2ValidasiAt,
			},
			"tahap3": map[string]interface{}{
				"submit_status":    k.Tahap3SubmitStatus,
				"submitted_at":     k.Tahap3SubmittedAt,
				"validasi_status":  k.Tahap3ValidasiStatus,
				"validasi_catatan": k.Tahap3ValidasiCatatan,
				"validasi_at":      k.Tahap3ValidasiAt,
			},
		})
	}
	return result, nil
}

// GetStatus ambil status validasi kontingen untuk widget dashboard.
func (s *Service) GetStatus(kontingenID uint) (map[string]interface{}, error) {
	data, err := s.repo.GetStatusByKontingen(kontingenID)
	if err != nil {
		return nil, errors.New("kontingen tidak ditemukan")
	}

	return map[string]interface{}{
		"kontingen_id":   data.KontingenID,
		"nama_kontingen": data.NamaKontingen,
		"tahap1": map[string]interface{}{
			"validasi_status":  data.Tahap1ValidasiStatus,
			"validasi_catatan": data.Tahap1ValidasiCatatan,
		},
		"tahap2": map[string]interface{}{
			"validasi_status":  data.Tahap2ValidasiStatus,
			"validasi_catatan": data.Tahap2ValidasiCatatan,
		},
		"tahap3": map[string]interface{}{
			"validasi_status":  data.Tahap3ValidasiStatus,
			"validasi_catatan": data.Tahap3ValidasiCatatan,
		},
	}, nil
}

// GetKontingenIDByTerritory resolve kontingen_id dari territory_id.
func (s *Service) GetKontingenIDByTerritory(territoryID uint) (uint, error) {
	return s.repo.GetKontingenIDByTerritory(territoryID)
}

// GetRekap ambil semua data pendaftaran kontingen dalam satu response.
// Dipakai halaman Rekap Pendaftaran (admin) dan detail validasi (superadmin).
func (s *Service) GetRekap(kontingenID uint) (map[string]interface{}, error) {
	return s.repo.GetRekap(kontingenID)
}

// SetValidasi simpan status VALID atau REVISI untuk satu tahap satu kontingen.
// Hanya superadmin — guard di handler.
func (s *Service) SetValidasi(kontingenID uint, tahap int, req SetValidasiRequest) (*ValidasiResult, error) {
	if tahap < 1 || tahap > 3 {
		return nil, errors.New("tahap harus 1, 2, atau 3")
	}

	status := strings.ToUpper(req.Status)
	if status != "VALID" && status != "REVISI" {
		return nil, errors.New("Status harus VALID atau REVISI")
	}

	// Catatan wajib jika REVISI
	if status == "REVISI" {
		if req.Catatan == nil || strings.TrimSpace(*req.Catatan) == "" {
			return nil, errors.New("Catatan wajib diisi jika status REVISI")
		}
	}

	// Cek apakah tahap sudah disubmit
	k, err := s.repo.GetKontingen(kontingenID)
	if err != nil {
		return nil, errors.New("kontingen tidak ditemukan")
	}

	var submitStatus string
	switch tahap {
	case 1:
		submitStatus = k.Tahap1Status
	case 2:
		submitStatus = k.Tahap2Status
	case 3:
		submitStatus = k.Tahap3Status
	}

	if submitStatus != "SUBMITTED" {
		return nil, errors.New("kontingen belum submit tahap ini, tidak ada yang bisa divalidasi")
	}

	// Jika VALID, kosongkan catatan
	var catatan *string
	if status == "REVISI" {
		catatan = req.Catatan
	}

	if err := s.repo.SetValidasi(kontingenID, tahap, status, catatan); err != nil {
		return nil, err
	}

	now := time.Now()
	return &ValidasiResult{
		KontingenID:   kontingenID,
		NamaKontingen: k.NamaKontingen,
		Tahap:         tahap,
		Status:        status,
		Catatan:       catatan,
		ValidasiAt:    &now,
	}, nil
}
