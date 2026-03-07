package tahap1

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetData(kontingenID uint) (map[string]interface{}, error) {

	meta, err := s.repo.GetMeta(kontingenID)
	if err != nil {
		return nil, err
	}

	cabor, _ := s.repo.GetCabor(kontingenID)

	return map[string]interface{}{
		"caborList":       cabor,
		"jumlahAtlet":     meta.JumlahAtlet,
		"jumlahPelatih":   meta.JumlahPelatih,
		"jumlahOfficial":  meta.JumlahOfficial,
		"tahap1Submitted": meta.Tahap1Submitted,
		"submittedAt":     meta.SubmittedAt,
	}, nil
}

func (s *Service) Update(kontingenID uint, cabor []uint, a, p, o int) error {
	if err := s.repo.SaveCabor(kontingenID, cabor); err != nil {
		return err
	}
	return s.repo.SaveMeta(kontingenID, a, p, o)
}

func (s *Service) Submit(kontingenID uint) error {

	meta, err := s.repo.GetMeta(kontingenID)
	if err != nil {
		return err
	}

	if meta.Tahap1Submitted {
		return errors.New("Tahap 1 sudah disubmit")
	}

	return s.repo.SetSubmitted(kontingenID)
}
