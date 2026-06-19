package sertifikat

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

// GetAll ambil sertifikat dengan filter opsional.
// Jika kontingenID > 0, filter hanya penerima milik kontingen tersebut
// dengan cara cek atlet/pelatih/official yang kontingen_id-nya cocok.
// Filter tambahan: tipe, atlet_id, pelatih_id, official_id.
func (s *Service) GetAll(filter map[string]interface{}) ([]Sertifikat, error) {
	return s.repo.GetAll(filter)
}

// GetByID ambil satu sertifikat.
func (s *Service) GetByID(id uint) (*Sertifikat, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("sertifikat tidak ditemukan")
	}
	return data, nil
}

// Create buat sertifikat baru.
// nama_penerima diisi otomatis dari nama_lengkap tabel master sesuai tipe.
func (s *Service) Create(req *CreateSertifikatRequest) (*Sertifikat, error) {
	// Validasi: pastikan FK yang sesuai tipe_penerima terisi
	var namaPenerima string
	var err error

	switch req.TipePenerima {
	case "ATLET":
		if req.AtletID == nil {
			return nil, errors.New("atlet_id wajib diisi jika tipe_penerima ATLET")
		}
		namaPenerima, err = s.repo.GetNamaAtlet(*req.AtletID)
		if err != nil {
			return nil, errors.New("atlet tidak ditemukan")
		}

	case "PELATIH":
		if req.PelatihID == nil {
			return nil, errors.New("pelatih_id wajib diisi jika tipe_penerima PELATIH")
		}
		namaPenerima, err = s.repo.GetNamaPelatih(*req.PelatihID)
		if err != nil {
			return nil, errors.New("pelatih tidak ditemukan")
		}

	case "OFFICIAL":
		if req.OfficialID == nil {
			return nil, errors.New("official_id wajib diisi jika tipe_penerima OFFICIAL")
		}
		namaPenerima, err = s.repo.GetNamaOfficial(*req.OfficialID)
		if err != nil {
			return nil, errors.New("official tidak ditemukan")
		}

	default:
		return nil, errors.New("tipe_penerima harus ATLET, PELATIH, atau OFFICIAL")
	}

	data := &Sertifikat{
		TipePenerima:    req.TipePenerima,
		AtletID:         req.AtletID,
		PelatihID:       req.PelatihID,
		OfficialID:      req.OfficialID,
		NamaPenerima:    namaPenerima, // otomatis dari tabel master
		Judul:           req.Judul,
		NomorSertifikat: req.NomorSertifikat,
		TanggalTerbit:   req.TanggalTerbit,
		Catatan:         req.Catatan,
	}

	if err := s.repo.Create(data); err != nil {
		return nil, fmt.Errorf("gagal membuat sertifikat: %w", err)
	}
	return data, nil
}

// Update ubah data sertifikat — partial, hanya field yang dikirim.
// NamaPenerima tidak bisa diubah via Update.
func (s *Service) Update(id uint, req *UpdateSertifikatRequest) (*Sertifikat, error) {
	// Pastikan sertifikat ada
	if _, err := s.repo.GetByID(id); err != nil {
		return nil, errors.New("sertifikat tidak ditemukan")
	}

	updates := make(map[string]interface{})
	if req.Judul != "" {
		updates["judul"] = req.Judul
	}
	if req.NomorSertifikat != nil {
		updates["nomor_sertifikat"] = req.NomorSertifikat
	}
	if req.TanggalTerbit != "" {
		updates["tanggal_terbit"] = req.TanggalTerbit
	}
	if req.Catatan != nil {
		updates["catatan"] = req.Catatan
	}

	if len(updates) == 0 {
		return s.repo.GetByID(id)
	}

	result, err := s.repo.Update(id, updates)
	if err != nil {
		return nil, fmt.Errorf("gagal mengupdate sertifikat: %w", err)
	}
	return result, nil
}

// Delete hapus sertifikat (hard delete).
func (s *Service) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.New("sertifikat tidak ditemukan")
	}
	return s.repo.Delete(id)
}

// UpdateFile update path file_sertifikat.
func (s *Service) UpdateFile(id uint, path string) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.New("sertifikat tidak ditemukan")
	}
	return s.repo.UpdateFile(id, path)
}

// ===== DROPDOWN PENERIMA =====

func (s *Service) GetAtletDropdown() ([]DropdownItem, error) {
	return s.repo.GetAtletDropdown()
}

func (s *Service) GetPelatihDropdown() ([]DropdownItem, error) {
	return s.repo.GetPelatihDropdown()
}

func (s *Service) GetOfficialDropdown() ([]DropdownItem, error) {
	return s.repo.GetOfficialDropdown()
}

// ===== DROPDOWN DATA =====
// Dipakai form sertifikat — return semua data tanpa filter kontingen.

func (s *Service) GetAllAtletSingkat() ([]PenerimaSingkat, error) {
	return s.repo.GetAllAtletSingkat()
}

func (s *Service) GetAllPelatihSingkat() ([]PenerimaSingkat, error) {
	return s.repo.GetAllPelatihSingkat()
}

func (s *Service) GetAllOfficialSingkat() ([]PenerimaSingkat, error) {
	return s.repo.GetAllOfficialSingkat()
}
