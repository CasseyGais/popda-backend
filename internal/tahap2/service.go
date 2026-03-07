package tahap2

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

	events, err := s.repo.GetEventsWithStatus(kontingenID)
	if err != nil {
		return nil, err
	}

	return &Response{
		Tahap2Submitted: submitted,
		SubmittedAt:     submittedAt,
		Events:          events,
	}, nil
}

func (s *Service) Update(kontingenID uint, nomorIDs []uint) error {
	return s.repo.SaveSelectedEvents(kontingenID, nomorIDs)
}

func (s *Service) Submit(kontingenID uint) error {

	submitted, _, err := s.repo.GetMeta(kontingenID)
	if err != nil {
		return err
	}

	if submitted {
		return errors.New("Tahap 2 sudah disubmit")
	}

	events, err := s.repo.GetEventsWithStatus(kontingenID)
	if err != nil {
		return err
	}

	count := 0
	for _, e := range events {
		if e.Ikut {
			count++
		}
	}

	if count == 0 {
		return errors.New("minimal 1 event harus dicentang")
	}

	return s.repo.SetSubmitted(kontingenID)
}
