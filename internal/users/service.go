package users

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]User, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*User, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}

func (s *Service) GetByEmail(email string) (*User, error) {
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}

func (s *Service) Create(request *CreateUserRequest) (*User, error) {
	// Check if email already exists
	existing, _ := s.repository.GetByEmail(request.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	user := &User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		IsActive:  true,
		Avatar:    request.Avatar,
	}

	err := s.repository.Create(user)
	if err != nil {
		return nil, errors.New("gagal membuat user")
	}

	return user, nil
}

func (s *Service) Update(id uint, request *UpdateUserRequest) (*User, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if request.Name != "" {
		user.Name = request.Name
	}
	if request.Email != "" {
		// Check if email already exists (excluding current)
		existing, _ := s.repository.GetByEmail(request.Email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("email sudah terdaftar")
		}
		user.Email = request.Email
	}
	if request.Password != "" {
		user.Password = request.Password
	}
	if request.Avatar != "" {
		user.Avatar = request.Avatar
	}

	err = s.repository.Update(user)
	if err != nil {
		return nil, errors.New("gagal mengupdate user")
	}

	return user, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus user")
	}
	return nil
}

func (s *Service) UpdateAvatar(id uint, avatar string) error {
	err := s.repository.UpdateAvatar(id, avatar)
	if err != nil {
		return errors.New("gagal mengupdate avatar user")
	}
	return nil
}

func (s *Service) UpdatePassword(id uint, password string) error {
	err := s.repository.UpdatePassword(id, password)
	if err != nil {
		return errors.New("gagal mengupdate password user")
	}
	return nil
}

func (s *Service) UpdateStatus(id uint, isActive bool) error {
	err := s.repository.UpdateStatus(id, isActive)
	if err != nil {
		return errors.New("gagal mengupdate status user")
	}
	return nil
}

func (s *Service) GetRoles(userID uint) ([]uint, error) {
	return s.repository.GetRoles(userID)
}

func (s *Service) AssignRole(userID, roleID uint) error {
	return s.repository.AssignRole(userID, roleID)
}

func (s *Service) RemoveRole(userID, roleID uint) error {
	return s.repository.RemoveRole(userID, roleID)
}

func (s *Service) GetTerritories(userID uint) ([]uint, error) {
	return s.repository.GetTerritories(userID)
}

func (s *Service) AssignTerritory(userID, territoryID uint) error {
	return s.repository.AssignTerritory(userID, territoryID)
}

func (s *Service) RemoveTerritory(userID, territoryID uint) error {
	return s.repository.RemoveTerritory(userID, territoryID)
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Avatar   string `json:"avatar"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Avatar   string `json:"avatar"`
}
