package masterpelatih

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]MasterPelatih, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*MasterPelatih, error) {
	pelatih, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}
	return pelatih, nil
}

func (s *Service) GetByKontingenID(kontingenID uint) ([]MasterPelatih, error) {
	return s.repository.GetByKontingenID(kontingenID)
}

func (s *Service) Create(request *CreateMasterPelatihRequest) (*MasterPelatih, error) {
	pelatih := &MasterPelatih{
		KontingenID: request.KontingenID,
		Nama:       request.Nama,
		NoHP:       request.NoHP,
		Sertifikat:  request.Sertifikat,
		Foto:       request.Foto,
	}

	err := s.repository.Create(pelatih)
	if err != nil {
		return nil, errors.New("gagal membuat pelatih")
	}

	return pelatih, nil
}

func (s *Service) Update(id uint, request *UpdateMasterPelatihRequest) (*MasterPelatih, error) {
	pelatih, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}

	if request.Nama != "" {
		pelatih.Nama = request.Nama
	}
	if request.NoHP != "" {
		pelatih.NoHP = request.NoHP
	}
	if request.Sertifikat != "" {
		pelatih.Sertifikat = request.Sertifikat
	}
	if request.Foto != "" {
		pelatih.Foto = request.Foto
	}

	err = s.repository.Update(pelatih)
	if err != nil {
		return nil, errors.New("gagal mengupdate pelatih")
	}

	return pelatih, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus pelatih")
	}
	return nil
}

func (s *Service) UpdateFoto(id uint, foto string) error {
	err := s.repository.UpdateFoto(id, foto)
	if err != nil {
		return errors.New("gagal mengupdate foto pelatih")
	}
	return nil
}

func (s *Service) UpdateSertifikat(id uint, sertifikat string) error {
	err := s.repository.UpdateSertifikat(id, sertifikat)
	if err != nil {
		return errors.New("gagal mengupdate sertifikat pelatih")
	}
	return nil
}

type CreateMasterPelatihRequest struct {
	KontingenID uint   `json:"kontingen_id" binding:"required"`
	Nama       string `json:"nama" binding:"required"`
	NoHP       string `json:"no_hp"`
	Sertifikat  string `json:"sertifikat"`
	Foto       string `json:"foto"`
}

type UpdateMasterPelatihRequest struct {
	Nama      string `json:"nama"`
	NoHP      string `json:"no_hp"`
	Sertifikat string `json:"sertifikat"`
	Foto      string `json:"foto"`
}
