package tahap3

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetData(kontingenID uint) (*Response, error) {

	submitted, submittedAt, err := s.repo.GetMeta(kontingenID)
	if err != nil {
		return nil, err
	}

	atlets, _ := s.repo.GetAtlets(kontingenID)
	pelatihs, _ := s.repo.GetPelatihs(kontingenID)
	officials, _ := s.repo.GetOfficials(kontingenID)

	return &Response{
		Tahap3Submitted: submitted,
		SubmittedAt:     submittedAt,
		Atlets:          atlets,
		Pelatihs:        pelatihs,
		Officials:       officials,
	}, nil
}

func (s *Service) Submit(kontingenID uint) error {

	atlets, _ := s.repo.GetAtlets(kontingenID)
	pelatihs, _ := s.repo.GetPelatihs(kontingenID)
	officials, _ := s.repo.GetOfficials(kontingenID)

	if len(atlets)+len(pelatihs)+len(officials) == 0 {
		return errors.New("minimal harus ada 1 data")
	}

	return s.repo.SetSubmitted(kontingenID)
}
