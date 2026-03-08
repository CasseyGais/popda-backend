package roles

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Role, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*Role, error) {
	role, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("role tidak ditemukan")
	}
	return role, nil
}

func (s *Service) GetByName(name string) (*Role, error) {
	role, err := s.repository.GetByName(name)
	if err != nil {
		return nil, errors.New("role tidak ditemukan")
	}
	return role, nil
}

func (s *Service) GetByUserID(userID uint) ([]Role, error) {
	return s.repository.GetByUserID(userID)
}

func (s *Service) Create(request *CreateRoleRequest) (*Role, error) {
	// Check if role already exists
	existing, _ := s.repository.GetByName(request.Name)
	if existing != nil {
		return nil, errors.New("role dengan nama tersebut sudah ada")
	}

	role := &Role{
		Name:        request.Name,
		Description: request.Description,
	}

	err := s.repository.Create(role)
	if err != nil {
		return nil, errors.New("gagal membuat role")
	}

	return role, nil
}

func (s *Service) Update(id uint, request *UpdateRoleRequest) (*Role, error) {
	role, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("role tidak ditemukan")
	}

	if request.Name != "" {
		// Check if role name already exists (excluding current)
		existing, _ := s.repository.GetByName(request.Name)
		if existing != nil && existing.ID != id {
			return nil, errors.New("role dengan nama tersebut sudah ada")
		}
		role.Name = request.Name
	}
	if request.Description != "" {
		role.Description = request.Description
	}

	err = s.repository.Update(role)
	if err != nil {
		return nil, errors.New("gagal mengupdate role")
	}

	return role, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus role")
	}
	return nil
}

func (s *Service) AssignPermission(roleID, permissionID uint) error {
	return s.repository.AssignPermission(roleID, permissionID)
}

func (s *Service) RemovePermission(roleID, permissionID uint) error {
	return s.repository.RemovePermission(roleID, permissionID)
}

func (s *Service) GetRolePermissions(roleID uint) ([]uint, error) {
	return s.repository.GetRolePermissions(roleID)
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
