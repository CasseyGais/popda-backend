package transaksi

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// TrxKontingenCabor operations
func (s *Service) CreateTrxKontingenCabor(kontingenID, caborID uint, putra, putri, pelatih int) error {
	totalAtlet := putra + putri
	totalPersonel := totalAtlet + pelatih

	trx := &TrxKontingenCabor{
		KontingenID:   kontingenID,
		CaborID:       caborID,
		Putra:         putra,
		Putri:         putri,
		Pelatih:       pelatih,
		TotalAtlet:    totalAtlet,
		TotalPersonel: totalPersonel,
	}

	return s.repo.CreateTrxKontingenCabor(trx)
}

func (s *Service) GetTrxKontingenCabor(kontingenID uint) ([]TrxKontingenCabor, error) {
	return s.repo.GetTrxKontingenCabor(kontingenID)
}

func (s *Service) UpdateTrxKontingenCabor(kontingenID, caborID uint, putra, putri, pelatih int) error {
	totalAtlet := putra + putri
	totalPersonel := totalAtlet + pelatih

	trx := &TrxKontingenCabor{
		Putra:         putra,
		Putri:         putri,
		Pelatih:       pelatih,
		TotalAtlet:    totalAtlet,
		TotalPersonel: totalPersonel,
	}

	return s.repo.UpdateTrxKontingenCabor(kontingenID, caborID, trx)
}

// TrxKontingenNomor operations
func (s *Service) CreateTrxKontingenNomor(kontingenID, nomorID uint) error {
	trx := &TrxKontingenNomor{
		KontingenID: kontingenID,
		NomorID:     nomorID,
	}

	return s.repo.CreateTrxKontingenNomor(trx)
}

func (s *Service) GetTrxKontingenNomor(kontingenID uint) ([]TrxKontingenNomor, error) {
	return s.repo.GetTrxKontingenNomor(kontingenID)
}

func (s *Service) DeleteTrxKontingenNomor(kontingenID, nomorID uint) error {
	return s.repo.DeleteTrxKontingenNomor(kontingenID, nomorID)
}

// TrxPendaftaranAtlet operations
func (s *Service) CreateTrxPendaftaranAtlet(atletID, nomorID, kelasID uint) error {
	trx := &TrxPendaftaranAtlet{
		AtletID: atletID,
		NomorID: nomorID,
		KelasID: &kelasID,
		Status:  "PENDING",
	}

	return s.repo.CreateTrxPendaftaranAtlet(trx)
}

func (s *Service) GetTrxPendaftaranAtlet(kontingenID uint) ([]TrxPendaftaranAtlet, error) {
	return s.repo.GetTrxPendaftaranAtlet(kontingenID)
}

func (s *Service) UpdateTrxPendaftaranAtlet(atletID, nomorID uint, status string) error {
	trx := &TrxPendaftaranAtlet{
		Status: status,
	}

	return s.repo.UpdateTrxPendaftaranAtlet(atletID, nomorID, trx)
}
