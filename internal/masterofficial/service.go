package masterofficial

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

// GetByKontingenID ambil semua official milik kontingen
func (s *Service) GetByKontingenID(kontingenID uint) ([]MasterOfficial, error) {
	return s.repo.GetByKontingenID(kontingenID)
}

// GetByID ambil satu official
func (s *Service) GetByID(id uint) (*MasterOfficial, error) {
	o, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("official tidak ditemukan")
	}
	return o, nil
}

// Create buat official baru — kontingen_id dari JWT/territory, bukan dari body
func (s *Service) Create(kontingenID uint, req *CreateMasterOfficialRequest) (*MasterOfficial, error) {
	o := &MasterOfficial{
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
	if err := s.repo.Create(o); err != nil {
		return nil, fmt.Errorf("gagal membuat official: %w", err)
	}
	return o, nil
}

// Update ubah data official — validasi kepemilikan kontingen
func (s *Service) Update(id uint, kontingenID uint, req *UpdateMasterOfficialRequest) (*MasterOfficial, error) {
	o, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("official tidak ditemukan")
	}
	if o.KontingenID != kontingenID {
		return nil, errors.New("tidak diizinkan mengubah data official kontingen lain")
	}
	if req.NamaLengkap != "" {
		o.NamaLengkap = req.NamaLengkap
	}
	if req.JenisKelamin != "" {
		o.JenisKelamin = req.JenisKelamin
	}
	if req.TanggalLahir != "" {
		o.TanggalLahir = req.TanggalLahir
	}
	if req.TempatLahir != "" {
		o.TempatLahir = req.TempatLahir
	}
	if req.NIK != "" {
		o.NIK = req.NIK
	}
	if req.SekolahAsal != "" {
		o.SekolahAsal = req.SekolahAsal
	}
	if req.Jabatan != "" {
		o.Jabatan = req.Jabatan
	}
	if req.Alamat != "" {
		o.Alamat = req.Alamat
	}
	if req.KabupatenKota != "" {
		o.KabupatenKota = req.KabupatenKota
	}
	if req.NoHP != "" {
		o.NoHP = req.NoHP
	}
	if req.Email != "" {
		o.Email = req.Email
	}
	if req.Catatan != "" {
		o.Catatan = req.Catatan
	}
	if err := s.repo.Update(o); err != nil {
		return nil, errors.New("gagal mengupdate official")
	}
	return o, nil
}

// Delete hapus official beserta trx-nya — validasi kepemilikan kontingen
func (s *Service) Delete(id uint, kontingenID uint) error {
	o, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("official tidak ditemukan")
	}
	if o.KontingenID != kontingenID {
		return errors.New("tidak diizinkan menghapus official kontingen lain")
	}
	if err := s.repo.DeleteTrxByOfficial(id); err != nil {
		return fmt.Errorf("gagal menghapus transaksi official: %w", err)
	}
	return s.repo.Delete(id)
}

// UpdateFile update path file/foto — validasi kepemilikan
func (s *Service) UpdateFile(id uint, kontingenID uint, column, path string) error {
	o, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("official tidak ditemukan")
	}
	if o.KontingenID != kontingenID {
		return errors.New("tidak diizinkan mengubah file official kontingen lain")
	}
	return s.repo.UpdateFile(id, column, path)
}

// ===== TRX PENDAFTARAN =====

// GetTrxByKontingen ambil daftar trx pendaftaran official milik kontingen
func (s *Service) GetTrxByKontingen(kontingenID uint) ([]TrxPendaftaranOfficial, error) {
	return s.repo.GetTrxByKontingen(kontingenID)
}

// CreateTrx daftarkan official — validasi official milik kontingen, otomatis isi trx_pendaftaran_official
func (s *Service) CreateTrx(kontingenID, officialID uint) (*TrxPendaftaranOfficial, error) {
	o, err := s.repo.GetByID(officialID)
	if err != nil {
		return nil, errors.New("official tidak ditemukan")
	}
	if o.KontingenID != kontingenID {
		return nil, errors.New("tidak diizinkan mendaftarkan official kontingen lain")
	}
	trx := &TrxPendaftaranOfficial{
		OfficialID: officialID,
	}
	if err := s.repo.CreateTrx(trx); err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan official: %w", err)
	}
	return trx, nil
}

// DeleteTrx batalkan pendaftaran official
func (s *Service) DeleteTrx(id uint) error {
	return s.repo.DeleteTrx(id)
}