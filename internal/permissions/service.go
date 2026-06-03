package permissions

import "errors"

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

func (s *Service) GetByModuleID(moduleID uint) ([]Permission, error) {
	return s.repository.GetByModuleID(moduleID)
}

func (s *Service) GetByRoleID(roleID uint) ([]Permission, error) {
	return s.repository.GetByRoleID(roleID)
}

func (s *Service) Create(req *CreatePermissionRequest) (*Permission, error) {
	// Cek duplikat name
	existing, _ := s.repository.GetByName(req.Name)
	if existing != nil {
		return nil, errors.New("permission dengan nama tersebut sudah ada")
	}

	permission := &Permission{
		ModuleID:    req.ModuleID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repository.Create(permission); err != nil {
		return nil, errors.New("gagal membuat permission")
	}
	return permission, nil
}

func (s *Service) Update(id uint, req *UpdatePermissionRequest) (*Permission, error) {
	permission, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("permission tidak ditemukan")
	}

	if req.Name != "" {
		// Cek duplikat name (kecuali diri sendiri)
		existing, _ := s.repository.GetByName(req.Name)
		if existing != nil && existing.ID != id {
			return nil, errors.New("permission dengan nama tersebut sudah ada")
		}
		permission.Name = req.Name
	}
	if req.ModuleID != nil {
		permission.ModuleID = req.ModuleID
	}
	if req.Description != "" {
		permission.Description = req.Description
	}

	if err := s.repository.Update(permission); err != nil {
		return nil, errors.New("gagal mengupdate permission")
	}
	return permission, nil
}

func (s *Service) Delete(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return errors.New("gagal menghapus permission")
	}
	return nil
}
