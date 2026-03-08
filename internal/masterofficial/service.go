package masterofficial

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]MasterOfficial, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*MasterOfficial, error) {
	official, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("official tidak ditemukan")
	}
	return official, nil
}

func (s *Service) GetByKontingenID(kontingenID uint) ([]MasterOfficial, error) {
	return s.repository.GetByKontingenID(kontingenID)
}

func (s *Service) Create(request *CreateMasterOfficialRequest) (*MasterOfficial, error) {
	official := &MasterOfficial{
		KontingenID: request.KontingenID,
		Nama:        request.Nama,
		Jabatan:     request.Jabatan,
		NoHP:        request.NoHP,
	}

	err := s.repository.Create(official)
	if err != nil {
		return nil, errors.New("gagal membuat official")
	}

	return official, nil
}

func (s *Service) Update(id uint, request *UpdateMasterOfficialRequest) (*MasterOfficial, error) {
	official, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("official tidak ditemukan")
	}

	if request.Nama != "" {
		official.Nama = request.Nama
	}
	if request.Jabatan != "" {
		official.Jabatan = request.Jabatan
	}
	if request.NoHP != "" {
		official.NoHP = request.NoHP
	}

	err = s.repository.Update(official)
	if err != nil {
		return nil, errors.New("gagal mengupdate official")
	}

	return official, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus official")
	}
	return nil
}

type CreateMasterOfficialRequest struct {
	KontingenID uint   `json:"kontingen_id" binding:"required"`
	Nama        string `json:"nama" binding:"required"`
	Jabatan     string `json:"jabatan" binding:"required"`
	NoHP        string `json:"no_hp"`
}

type UpdateMasterOfficialRequest struct {
	Nama    string `json:"nama"`
	Jabatan string `json:"jabatan"`
	NoHP    string `json:"no_hp"`
}
