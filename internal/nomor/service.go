package nomor

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Nomor, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*Nomor, error) {
	nomor, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("nomor pertandingan tidak ditemukan")
	}
	return nomor, nil
}

func (s *Service) GetByCaborID(caborID uint) ([]Nomor, error) {
	return s.repository.GetByCaborID(caborID)
}

func (s *Service) Create(request *CreateNomorRequest) (*Nomor, error) {
	nomor := &Nomor{
		Nama:         request.Nama,
		CaborID:      request.CaborID,
		JenisKelamin: request.JenisKelamin,
		Tipe:         request.Tipe,
		IsActive:     true,
	}

	err := s.repository.Create(nomor)
	if err != nil {
		return nil, errors.New("gagal membuat nomor pertandingan")
	}

	return nomor, nil
}

func (s *Service) Update(id uint, request *UpdateNomorRequest) (*Nomor, error) {
	nomor, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("nomor pertandingan tidak ditemukan")
	}

	if request.Nama != "" {
		nomor.Nama = request.Nama
	}
	if request.CaborID != 0 {
		nomor.CaborID = request.CaborID
	}
	if request.JenisKelamin != "" {
		nomor.JenisKelamin = request.JenisKelamin
	}
	if request.Tipe != "" {
		nomor.Tipe = request.Tipe
	}

	err = s.repository.Update(nomor)
	if err != nil {
		return nil, errors.New("gagal mengupdate nomor pertandingan")
	}

	return nomor, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus nomor pertandingan")
	}
	return nil
}

type CreateNomorRequest struct {
	Nama         string `json:"nama" binding:"required"`
	CaborID      uint   `json:"cabor_id" binding:"required"`
	JenisKelamin string `json:"jenis_kelamin" binding:"required,oneof=PUTRA PUTRI CAMPURAN"`
	Tipe         string `json:"tipe" binding:"required,oneof=INDIVIDU BEREGU"`
}

type UpdateNomorRequest struct {
	Nama         string `json:"nama"`
	CaborID      uint   `json:"cabor_id"`
	JenisKelamin string `json:"jenis_kelamin" binding:"omitempty,oneof=PUTRA PUTRI CAMPURAN"`
	Tipe         string `json:"tipe" binding:"omitempty,oneof=INDIVIDU BEREGU"`
}
