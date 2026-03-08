package cabor

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Cabor, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*Cabor, error) {
	cabor, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("cabang olahraga tidak ditemukan")
	}
	return cabor, nil
}

func (s *Service) Create(request *CreateCaborRequest) (*Cabor, error) {
	cabor := &Cabor{
		Nama:       request.Nama,
		MaxPutra:   request.MaxPutra,
		MaxPutri:   request.MaxPutri,
		MaxPelatih: request.MaxPelatih,
		IsActive:   true,
	}

	err := s.repository.Create(cabor)
	if err != nil {
		return nil, errors.New("gagal membuat cabang olahraga")
	}

	return cabor, nil
}

func (s *Service) Update(id uint, request *UpdateCaborRequest) (*Cabor, error) {
	cabor, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("cabang olahraga tidak ditemukan")
	}

	if request.Nama != "" {
		cabor.Nama = request.Nama
	}
	if request.MaxPutra != 0 {
		cabor.MaxPutra = request.MaxPutra
	}
	if request.MaxPutri != 0 {
		cabor.MaxPutri = request.MaxPutri
	}
	if request.MaxPelatih != 0 {
		cabor.MaxPelatih = request.MaxPelatih
	}

	err = s.repository.Update(cabor)
	if err != nil {
		return nil, errors.New("gagal mengupdate cabang olahraga")
	}

	return cabor, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus cabang olahraga")
	}
	return nil
}

type CreateCaborRequest struct {
	Nama       string `json:"nama" binding:"required"`
	MaxPutra   int    `json:"max_putra"`
	MaxPutri   int    `json:"max_putri"`
	MaxPelatih int    `json:"max_pelatih"`
}

type UpdateCaborRequest struct {
	Nama       string `json:"nama"`
	MaxPutra   int    `json:"max_putra"`
	MaxPutri   int    `json:"max_putri"`
	MaxPelatih int    `json:"max_pelatih"`
}
