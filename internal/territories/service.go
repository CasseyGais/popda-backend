package territories

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Territory, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*Territory, error) {
	territory, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("territory tidak ditemukan")
	}
	return territory, nil
}

func (s *Service) GetByType(territoryType string) ([]Territory, error) {
	if territoryType != "PROVINSI" && territoryType != "KABUPATEN" && territoryType != "KOTA" {
		return nil, errors.New("type territory tidak valid")
	}
	return s.repository.GetByType(territoryType)
}

func (s *Service) GetProvinces() ([]Territory, error) {
	return s.repository.GetProvinces()
}

func (s *Service) GetKabupatens() ([]Territory, error) {
	return s.repository.GetKabupatens()
}

func (s *Service) GetKotas() ([]Territory, error) {
	return s.repository.GetKotas()
}

func (s *Service) GetByUserID(userID uint) ([]Territory, error) {
	return s.repository.GetByUserID(userID)
}

func (s *Service) Create(request *CreateTerritoryRequest) (*Territory, error) {
	// Validate type
	if request.Type != "PROVINSI" && request.Type != "KABUPATEN" && request.Type != "KOTA" {
		return nil, errors.New("type territory tidak valid")
	}

	territory := &Territory{
		Name: request.Name,
		Type: request.Type,
	}

	err := s.repository.Create(territory)
	if err != nil {
		return nil, errors.New("gagal membuat territory")
	}

	return territory, nil
}

func (s *Service) Update(id uint, request *UpdateTerritoryRequest) (*Territory, error) {
	territory, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("territory tidak ditemukan")
	}

	if request.Name != "" {
		territory.Name = request.Name
	}
	if request.Type != "" {
		// Validate type
		if request.Type != "PROVINSI" && request.Type != "KABUPATEN" && request.Type != "KOTA" {
			return nil, errors.New("type territory tidak valid")
		}
		territory.Type = request.Type
	}

	err = s.repository.Update(territory)
	if err != nil {
		return nil, errors.New("gagal mengupdate territory")
	}

	return territory, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus territory")
	}
	return nil
}

func (s *Service) AssignToUser(userID, territoryID uint) error {
	return s.repository.AssignToUser(userID, territoryID)
}

func (s *Service) RemoveFromUser(userID, territoryID uint) error {
	return s.repository.RemoveFromUser(userID, territoryID)
}

type CreateTerritoryRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required,oneof=PROVINSI KABUPATEN KOTA"`
}

type UpdateTerritoryRequest struct {
	Name string `json:"name"`
	Type string `json:"type" binding:"omitempty,oneof=PROVINSI KABUPATEN KOTA"`
}
