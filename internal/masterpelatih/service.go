package masterpelatih

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

// GetKontingenIDByTerritory untuk handler resolveKontingenID
func (s *Service) GetKontingenIDByTerritory(territoryID uint) (uint, error) {
	return s.repo.GetKontingenIDByTerritory(territoryID)
}

// GetByKontingenID ambil semua pelatih milik kontingen
func (s *Service) GetByKontingenID(kontingenID uint) ([]MasterPelatih, error) {
	return s.repo.GetByKontingenID(kontingenID)
}

// GetByID ambil satu pelatih
func (s *Service) GetByID(id uint) (*MasterPelatih, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}
	return p, nil
}

// Create buat pelatih baru — kontingen_id dari JWT/territory, bukan dari body
func (s *Service) Create(kontingenID uint, req *CreateMasterPelatihRequest) (*MasterPelatih, error) {
	p := &MasterPelatih{
		KontingenID:        kontingenID,
		NamaLengkap:        req.NamaLengkap,
		JenisKelamin:       req.JenisKelamin,
		TanggalLahir:       req.TanggalLahir,
		TempatLahir:        req.TempatLahir,
		NIK:                req.NIK,
		SekolahAsal:        req.SekolahAsal,
		Profesi:            req.Profesi,
		Jabatan:            req.Jabatan,
		Alamat:             req.Alamat,
		KabupatenKota:      req.KabupatenKota,
		NoHP:               req.NoHP,
		Email:              req.Email,
		NamaIstriSuami:     req.NamaIstriSuami,
		PrestasiSebelumnya: req.PrestasiSebelumnya,
		Catatan:            req.Catatan,
		Status:             "draft",
	}
	if err := s.repo.Create(p); err != nil {
		return nil, fmt.Errorf("gagal membuat pelatih: %w", err)
	}
	return p, nil
}

// Update ubah data pelatih — validasi kepemilikan kontingen
func (s *Service) Update(id uint, kontingenID uint, req *UpdateMasterPelatihRequest) (*MasterPelatih, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}
	// Pastikan pelatih ini milik kontingen yang sedang request
	if p.KontingenID != kontingenID {
		return nil, errors.New("tidak diizinkan mengubah data pelatih kontingen lain")
	}
	if req.NamaLengkap != "" {
		p.NamaLengkap = req.NamaLengkap
	}
	if req.JenisKelamin != "" {
		p.JenisKelamin = req.JenisKelamin
	}
	if req.TanggalLahir != "" {
		p.TanggalLahir = req.TanggalLahir
	}
	if req.TempatLahir != "" {
		p.TempatLahir = req.TempatLahir
	}
	if req.NIK != "" {
		p.NIK = req.NIK
	}
	if req.SekolahAsal != "" {
		p.SekolahAsal = req.SekolahAsal
	}
	if req.Profesi != "" {
		p.Profesi = req.Profesi
	}
	if req.Jabatan != "" {
		p.Jabatan = req.Jabatan
	}
	if req.Alamat != "" {
		p.Alamat = req.Alamat
	}
	if req.KabupatenKota != "" {
		p.KabupatenKota = req.KabupatenKota
	}
	if req.NoHP != "" {
		p.NoHP = req.NoHP
	}
	if req.Email != "" {
		p.Email = req.Email
	}
	if req.NamaIstriSuami != "" {
		p.NamaIstriSuami = req.NamaIstriSuami
	}
	if req.PrestasiSebelumnya != "" {
		p.PrestasiSebelumnya = req.PrestasiSebelumnya
	}
	if req.Catatan != "" {
		p.Catatan = req.Catatan
	}
	if err := s.repo.Update(p); err != nil {
		return nil, errors.New("gagal mengupdate pelatih")
	}
	return p, nil
}

// Delete hapus pelatih beserta trx-nya — validasi kepemilikan kontingen
func (s *Service) Delete(id uint, kontingenID uint) error {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("pelatih tidak ditemukan")
	}
	if p.KontingenID != kontingenID {
		return errors.New("tidak diizinkan menghapus pelatih kontingen lain")
	}
	if err := s.repo.DeleteTrxByPelatih(id); err != nil {
		return fmt.Errorf("gagal menghapus transaksi pelatih: %w", err)
	}
	return s.repo.Delete(id)
}

// UpdateFile update path file/foto — validasi kepemilikan
func (s *Service) UpdateFile(id uint, kontingenID uint, column, path string) error {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("pelatih tidak ditemukan")
	}
	if p.KontingenID != kontingenID {
		return errors.New("tidak diizinkan mengubah file pelatih kontingen lain")
	}
	return s.repo.UpdateFile(id, column, path)
}

// ===== TRX PENDAFTARAN =====

// GetTrxByKontingen ambil daftar trx pendaftaran pelatih milik kontingen
func (s *Service) GetTrxByKontingen(kontingenID uint) ([]TrxPendaftaranPelatih, error) {
	return s.repo.GetTrxByKontingen(kontingenID)
}

// CreateTrx daftarkan pelatih ke cabor — validasi pelatih milik kontingen
func (s *Service) CreateTrx(kontingenID, pelatihID, caborID uint) (*TrxPendaftaranPelatih, error) {
	p, err := s.repo.GetByID(pelatihID)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}
	if p.KontingenID != kontingenID {
		return nil, errors.New("tidak diizinkan mendaftarkan pelatih kontingen lain")
	}
	trx := &TrxPendaftaranPelatih{
		PelatihID: pelatihID,
		CaborID:   caborID,
	}
	if err := s.repo.CreateTrx(trx); err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan pelatih: %w", err)
	}
	return trx, nil
}

// DeleteTrx batalkan pendaftaran pelatih
func (s *Service) DeleteTrx(id uint) error {
	return s.repo.DeleteTrx(id)
}
