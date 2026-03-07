package kontingen

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ================= IDENTITAS KONTINGEN =================
func (s *Service) GetIdentitas(kontingenID uint) (*IdentitasKontingen, error) {
	return s.repo.GetByKontingenID(kontingenID)
}

func (s *Service) UpdateIdentitas(kontingenID uint, input *IdentitasKontingen) error {
	input.KontingenID = kontingenID

	existing, err := s.repo.GetByKontingenID(kontingenID)
	if err != nil {
		return err
	}

	// Kalau belum ada → INSERT
	if existing == nil {
		// Cek apakah kontingen ada di database
		kontingenExists, err := s.repo.KontingenExists(kontingenID)
		if err != nil {
			return err
		}

		if !kontingenExists {
			// Buat kontingen baru dengan territory_id yang sama
			err = s.repo.CreateKontingenForID(kontingenID)
			if err != nil {
				return err
			}
		}

		return s.repo.Create(input)
	}

	// Kalau sudah ada → UPDATE
	return s.repo.Update(kontingenID, input)
}

// ================= KONTINGEN =================
func (s *Service) GetKontingenByID(id uint) (*Kontingen, error) {
	return s.repo.GetKontingenByID(id)
}

func (s *Service) GetKontingenByTerritoryID(territoryID uint) (*Kontingen, error) {
	return s.repo.GetKontingenByTerritoryID(territoryID)
}

func (s *Service) CreateKontingen(kontingen *Kontingen) error {
	return s.repo.CreateKontingen(kontingen)
}

func (s *Service) UpdateKontingen(id uint, kontingen *Kontingen) error {
	return s.repo.UpdateKontingen(id, kontingen)
}
