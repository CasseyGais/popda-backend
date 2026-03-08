package kontingenidentitas

import (
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]KontingenIdentitas, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*KontingenIdentitas, error) {
	identitas, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("identitas kontingen tidak ditemukan")
	}
	return identitas, nil
}

func (s *Service) GetByKontingenID(kontingenID uint) (*KontingenIdentitas, error) {
	identitas, err := s.repository.GetByKontingenID(kontingenID)
	if err != nil {
		return nil, errors.New("identitas kontingen tidak ditemukan")
	}
	return identitas, nil
}

func (s *Service) Create(request *CreateKontingenIdentitasRequest) (*KontingenIdentitas, error) {
	identitas := &KontingenIdentitas{
		KontingenID:   request.KontingenID,
		KepalaNama:    request.KepalaNama,
		KepalaJabatan: request.KepalaJabatan,
		KepalaNIP:     request.KepalaNIP,
		KepalaTelepon: request.KepalaTelepon,
		KepalaFoto:    request.KepalaFoto,
		PICNama:       request.PICNama,
		PICJabatan:    request.PICJabatan,
		PICTelepon:    request.PICTelepon,
		PICFoto:       request.PICFoto,
		Alamat:        request.Alamat,
		EmailInstansi: request.EmailInstansi,
		PhoneInstansi: request.PhoneInstansi,
	}

	err := s.repository.Create(identitas)
	if err != nil {
		return nil, errors.New("gagal membuat identitas kontingen")
	}

	return identitas, nil
}

func (s *Service) Update(id uint, request *UpdateKontingenIdentitasRequest) (*KontingenIdentitas, error) {
	identitas, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("identitas kontingen tidak ditemukan")
	}

	if request.KepalaNama != "" {
		identitas.KepalaNama = request.KepalaNama
	}
	if request.KepalaJabatan != "" {
		identitas.KepalaJabatan = request.KepalaJabatan
	}
	if request.KepalaNIP != "" {
		identitas.KepalaNIP = request.KepalaNIP
	}
	if request.KepalaTelepon != "" {
		identitas.KepalaTelepon = request.KepalaTelepon
	}
	if request.KepalaFoto != "" {
		identitas.KepalaFoto = request.KepalaFoto
	}
	if request.PICNama != "" {
		identitas.PICNama = request.PICNama
	}
	if request.PICJabatan != "" {
		identitas.PICJabatan = request.PICJabatan
	}
	if request.PICTelepon != "" {
		identitas.PICTelepon = request.PICTelepon
	}
	if request.PICFoto != "" {
		identitas.PICFoto = request.PICFoto
	}
	if request.Alamat != "" {
		identitas.Alamat = request.Alamat
	}
	if request.EmailInstansi != "" {
		identitas.EmailInstansi = request.EmailInstansi
	}
	if request.PhoneInstansi != "" {
		identitas.PhoneInstansi = request.PhoneInstansi
	}

	err = s.repository.Update(identitas)
	if err != nil {
		return nil, errors.New("gagal mengupdate identitas kontingen")
	}

	return identitas, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus identitas kontingen")
	}
	return nil
}

func (s *Service) UpdateKepalaFoto(id uint, foto string) error {
	err := s.repository.UpdateKepalaFoto(id, foto)
	if err != nil {
		return errors.New("gagal mengupdate foto kepala")
	}
	return nil
}

func (s *Service) UpdatePICFoto(id uint, foto string) error {
	err := s.repository.UpdatePICFoto(id, foto)
	if err != nil {
		return errors.New("gagal mengupdate foto PIC")
	}
	return nil
}

type CreateKontingenIdentitasRequest struct {
	KontingenID   uint   `json:"kontingen_id" binding:"required"`
	KepalaNama    string `json:"kepala_nama"`
	KepalaJabatan string `json:"kepala_jabatan"`
	KepalaNIP     string `json:"kepala_nip"`
	KepalaTelepon string `json:"kepala_telepon"`
	KepalaFoto    string `json:"kepala_foto"`
	PICNama       string `json:"pic_nama"`
	PICJabatan    string `json:"pic_jabatan"`
	PICTelepon    string `json:"pic_telepon"`
	PICFoto       string `json:"pic_foto"`
	Alamat        string `json:"alamat"`
	EmailInstansi string `json:"email_instansi"`
	PhoneInstansi string `json:"phone_instansi"`
}

type UpdateKontingenIdentitasRequest struct {
	KepalaNama    string `json:"kepala_nama"`
	KepalaJabatan string `json:"kepala_jabatan"`
	KepalaNIP     string `json:"kepala_nip"`
	KepalaTelepon string `json:"kepala_telepon"`
	KepalaFoto    string `json:"kepala_foto"`
	PICNama       string `json:"pic_nama"`
	PICJabatan    string `json:"pic_jabatan"`
	PICTelepon    string `json:"pic_telepon"`
	PICFoto       string `json:"pic_foto"`
	Alamat        string `json:"alamat"`
	EmailInstansi string `json:"email_instansi"`
	PhoneInstansi string `json:"phone_instansi"`
}
