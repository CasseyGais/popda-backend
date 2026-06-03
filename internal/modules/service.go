package modules

import "errors"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Module, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*Module, error) {
	module, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("module tidak ditemukan")
	}
	return module, nil
}

func (s *Service) Create(req *CreateModuleRequest) (*Module, error) {
	// Cek duplikat code
	existing, _ := s.repository.GetByCode(req.Code)
	if existing != nil {
		return nil, errors.New("module dengan code tersebut sudah ada")
	}

	module := &Module{
		Name:  req.Name,
		Label: req.Label,
		Code:  req.Code,
		URL:   req.URL,
	}

	if err := s.repository.Create(module); err != nil {
		return nil, errors.New("gagal membuat module")
	}
	return module, nil
}

func (s *Service) Update(id uint, req *UpdateModuleRequest) (*Module, error) {
	module, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("module tidak ditemukan")
	}

	if req.Name != "" {
		module.Name = req.Name
	}
	if req.Label != "" {
		module.Label = req.Label
	}
	if req.Code != "" {
		// Cek duplikat code (kecuali diri sendiri)
		existing, _ := s.repository.GetByCode(req.Code)
		if existing != nil && existing.ID != id {
			return nil, errors.New("module dengan code tersebut sudah ada")
		}
		module.Code = req.Code
	}
	if req.URL != "" {
		module.URL = req.URL
	}

	if err := s.repository.Update(module); err != nil {
		return nil, errors.New("gagal mengupdate module")
	}
	return module, nil
}

func (s *Service) Delete(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return errors.New("gagal menghapus module")
	}
	return nil
}
