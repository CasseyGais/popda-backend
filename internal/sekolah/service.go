package sekolah

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Sekolah, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*Sekolah, error) {
	sekolah, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("sekolah tidak ditemukan")
	}
	return sekolah, nil
}

func (s *Service) GetByNPSN(npsn string) (*Sekolah, error) {
	sekolah, err := s.repository.GetByNPSN(npsn)
	if err != nil {
		return nil, errors.New("sekolah dengan NPSN tersebut tidak ditemukan")
	}
	return sekolah, nil
}

func (s *Service) Search(keyword string) ([]Sekolah, error) {
	if keyword == "" {
		return s.repository.GetAll()
	}
	return s.repository.Search(keyword)
}

func (s *Service) Create(request *CreateSekolahRequest) (*Sekolah, error) {
	// Check if NPSN already exists
	existingSekolah, _ := s.repository.GetByNPSN(request.NPSN)
	if existingSekolah != nil {
		return nil, errors.New("NPSN sudah terdaftar")
	}

	sekolah := &Sekolah{
		Name:      request.Name,
		NPSN:      request.NPSN,
		Alamat:    request.Alamat,
		Kabupaten: request.Kabupaten,
	}

	err := s.repository.Create(sekolah)
	if err != nil {
		return nil, errors.New("gagal membuat sekolah")
	}

	return sekolah, nil
}

func (s *Service) Update(id uint, request *UpdateSekolahRequest) (*Sekolah, error) {
	sekolah, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("sekolah tidak ditemukan")
	}

	if request.Name != "" {
		sekolah.Name = request.Name
	}
	if request.NPSN != "" {
		// Check if NPSN already exists for other sekolah
		existingSekolah, _ := s.repository.GetByNPSN(request.NPSN)
		if existingSekolah != nil && existingSekolah.ID != id {
			return nil, errors.New("NPSN sudah digunakan oleh sekolah lain")
		}
		sekolah.NPSN = request.NPSN
	}
	if request.Alamat != "" {
		sekolah.Alamat = request.Alamat
	}
	if request.Kabupaten != "" {
		sekolah.Kabupaten = request.Kabupaten
	}

	err = s.repository.Update(sekolah)
	if err != nil {
		return nil, errors.New("gagal mengupdate sekolah")
	}

	return sekolah, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus sekolah")
	}
	return nil
}

type CreateSekolahRequest struct {
	Name      string `json:"nama" binding:"required"`
	NPSN      string `json:"npsn" binding:"required"`
	Alamat    string `json:"alamat"`
	Kabupaten string `json:"kabupaten"`
}

type UpdateSekolahRequest struct {
	Name      string `json:"nama"`
	NPSN      string `json:"npsn"`
	Alamat    string `json:"alamat"`
	Kabupaten string `json:"kabupaten"`
}
