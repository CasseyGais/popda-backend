package atlet

import (
	"errors"
	"time"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Atlet, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id uint) (*Atlet, error) {
	atlet, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("atlet tidak ditemukan")
	}
	return atlet, nil
}

func (s *Service) GetByKontingenID(kontingenID uint) ([]Atlet, error) {
	return s.repository.GetByKontingenID(kontingenID)
}

func (s *Service) GetBySekolahID(sekolahID uint) ([]Atlet, error) {
	return s.repository.GetBySekolahID(sekolahID)
}

func (s *Service) GetByStatus(status string) ([]Atlet, error) {
	return s.repository.GetByStatus(status)
}

func (s *Service) Create(request *CreateAtletRequest) (*Atlet, error) {
	atlet := &Atlet{
		KontingenID:     request.KontingenID,
		SekolahID:       request.SekolahID,
		NISN:            request.NISN,
		Nama:            request.Nama,
		JenisKelamin:    request.JenisKelamin,
		TanggalLahir:    request.TanggalLahir,
		Kelas:           request.Kelas,
		Tinggi:          request.Tinggi,
		Berat:           request.Berat,
		Foto:            request.Foto,
		StatusVerifikasi: "PENDING",
	}

	err := s.repository.Create(atlet)
	if err != nil {
		return nil, errors.New("gagal membuat atlet")
	}

	return atlet, nil
}

func (s *Service) Update(id uint, request *UpdateAtletRequest) (*Atlet, error) {
	atlet, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("atlet tidak ditemukan")
	}

	if request.NISN != "" {
		atlet.NISN = request.NISN
	}
	if request.Nama != "" {
		atlet.Nama = request.Nama
	}
	if request.JenisKelamin != "" {
		atlet.JenisKelamin = request.JenisKelamin
	}
	if request.TanggalLahir != nil {
		atlet.TanggalLahir = request.TanggalLahir
	}
	if request.Kelas != "" {
		atlet.Kelas = request.Kelas
	}
	if request.Tinggi != nil {
		atlet.Tinggi = request.Tinggi
	}
	if request.Berat != nil {
		atlet.Berat = request.Berat
	}
	if request.Foto != "" {
		atlet.Foto = request.Foto
	}

	err = s.repository.Update(atlet)
	if err != nil {
		return nil, errors.New("gagal mengupdate atlet")
	}

	return atlet, nil
}

func (s *Service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus atlet")
	}
	return nil
}

func (s *Service) UpdateStatus(id uint, status string) error {
	if status != "PENDING" && status != "VALID" && status != "DITOLAK" {
		return errors.New("status verifikasi tidak valid")
	}

	err := s.repository.UpdateStatus(id, status)
	if err != nil {
		return errors.New("gagal mengupdate status verifikasi")
	}
	return nil
}

func (s *Service) UpdateFoto(id uint, foto string) error {
	err := s.repository.UpdateFoto(id, foto)
	if err != nil {
		return errors.New("gagal mengupdate foto atlet")
	}
	return nil
}

type CreateAtletRequest struct {
	KontingenID  uint       `json:"kontingen_id" binding:"required"`
	SekolahID    uint       `json:"sekolah_id" binding:"required"`
	NISN         string     `json:"nisn"`
	Nama         string     `json:"nama" binding:"required"`
	JenisKelamin string     `json:"jenis_kelamin" binding:"required,oneof=PUTRA PUTRI"`
	TanggalLahir *time.Time `json:"tanggal_lahir"`
	Kelas        string     `json:"kelas"`
	Tinggi       *int       `json:"tinggi"`
	Berat        *float64   `json:"berat"`
	Foto         string     `json:"foto"`
}

type UpdateAtletRequest struct {
	NISN         string     `json:"nisn"`
	Nama         string     `json:"nama"`
	JenisKelamin string     `json:"jenis_kelamin" binding:"omitempty,oneof=PUTRA PUTRI"`
	TanggalLahir *time.Time `json:"tanggal_lahir"`
	Kelas        string     `json:"kelas"`
	Tinggi       *int       `json:"tinggi"`
	Berat        *float64   `json:"berat"`
	Foto         string     `json:"foto"`
}
