package pengaturantahap

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

// GetAll kembalikan status semua tahap.
// Dipakai GET /admin/pengaturan-tahap (semua role).
func (s *Service) GetAll() ([]PengaturanTahap, error) {
	return s.repo.GetAll()
}

// IsOpen cek apakah tahap tertentu sedang dibuka.
// Dipakai middleware checkTahapOpen.
func (s *Service) IsOpen(tahap uint) (bool, error) {
	return s.repo.IsOpen(tahap)
}

// Update ubah pengaturan satu tahap.
// Hanya superadmin yang boleh memanggil ini (guard di handler).
//
// Aturan urutan:
//   - Tahap 2 tidak bisa dibuka sebelum tahap 1 pernah dibuka
//   - Tahap 3 tidak bisa dibuka sebelum tahap 2 pernah dibuka
func (s *Service) Update(tahap uint, req UpdatePengaturanRequest) (*PengaturanTahap, error) {
	if tahap < 1 || tahap > 3 {
		return nil, errors.New("tahap harus 1, 2, atau 3")
	}

	// Validasi urutan: jika mau membuka (is_open = true), cek prasyarat tahap sebelumnya
	if req.IsOpen != nil && *req.IsOpen {
		if tahap == 2 {
			pernah, err := s.repo.PernahDibuka(1)
			if err != nil {
				return nil, err
			}
			if !pernah {
				return nil, errors.New("Tahap 2 tidak bisa dibuka sebelum Tahap 1 pernah dibuka")
			}
		}
		if tahap == 3 {
			pernah, err := s.repo.PernahDibuka(2)
			if err != nil {
				return nil, err
			}
			if !pernah {
				return nil, errors.New("Tahap 3 tidak bisa dibuka sebelum Tahap 2 pernah dibuka")
			}
		}
	}

	// Bangun map update — hanya field yang dikirim
	updates := make(map[string]interface{})
	if req.IsOpen != nil {
		updates["is_open"] = *req.IsOpen
		// Jika baru saja dibuka dan tanggal_buka belum diset, biarkan frontend set via tanggal_buka
	}
	if req.TanggalBuka != nil {
		if *req.TanggalBuka == "" {
			updates["tanggal_buka"] = nil
		} else {
			updates["tanggal_buka"] = *req.TanggalBuka
		}
	}
	if req.TanggalTutup != nil {
		if *req.TanggalTutup == "" {
			updates["tanggal_tutup"] = nil
		} else {
			updates["tanggal_tutup"] = *req.TanggalTutup
		}
	}

	if len(updates) == 0 {
		// Tidak ada yang diubah, kembalikan data saat ini
		return s.repo.GetByTahap(tahap)
	}

	result, err := s.repo.Update(tahap, updates)
	if err != nil {
		return nil, fmt.Errorf("gagal update pengaturan tahap %d: %w", tahap, err)
	}
	return result, nil
}
