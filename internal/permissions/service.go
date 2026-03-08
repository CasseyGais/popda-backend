package permissions

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Permission, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*Permission, error) {
	permission, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("permission tidak ditemukan")
	}
	return permission, nil
}

func (s *Service) GetByName(name string) (*Permission, error) {
	permission, err := s.repository.GetByName(name)
	if err != nil {
		return nil, errors.New("permission tidak ditemukan")
	}
	return permission, nil
}

func (s *Service) GetByRoleID(roleID uint) ([]Permission, error) {
	return s.repository.GetByRoleID(roleID)
}

func (s *Service) Create(request *CreatePermissionRequest) (*Permission, error) {
	// Check if permission already exists
	existing, _ := s.repository.GetByName(request.Name)
	if existing != nil {
		return nil, errors.New("permission dengan nama tersebut sudah ada")
	}

	permission := &Permission{
		Name:        request.Name,
		Description: request.Description,
	}

	err := s.repository.Create(permission)
	if err != nil {
		return nil, errors.New("gagal membuat permission")
	}

	return permission, nil
}

func (s *Service) Update(id uint, request *UpdatePermissionRequest) (*Permission, error) {
	permission, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("permission tidak ditemukan")
	}

	if request.Name != "" {
		// Check if permission name already exists (excluding current)
		existing, _ := s.repository.GetByName(request.Name)
		if existing != nil && existing.ID != id {
			return nil, errors.New("permission dengan nama tersebut sudah ada")
		}
		permission.Name = request.Name
	}
	if request.Description != "" {
		permission.Description = request.Description
	}

	err = s.repository.Update(permission)
	if err != nil {
		return nil, errors.New("gagal mengupdate permission")
	}

	return permission, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus permission")
	}
	return nil
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
