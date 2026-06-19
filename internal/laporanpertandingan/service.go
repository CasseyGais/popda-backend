package laporanpertandingan

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetAll ambil semua laporan dengan filter opsional
func (s *Service) GetAll(filter FilterLaporan) ([]LaporanDetail, error) {
	rows, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}
	// Populate atlet per row
	for i := range rows {
		atletA, _ := s.repo.GetAtletBySisi(rows[i].ID, "A")
		atletB, _ := s.repo.GetAtletBySisi(rows[i].ID, "B")
		if atletA == nil {
			atletA = []AtletSisiItem{}
		}
		if atletB == nil {
			atletB = []AtletSisiItem{}
		}
		rows[i].AtletA = atletA
		rows[i].AtletB = atletB
	}
	return rows, nil
}

// GetByID ambil satu laporan lengkap dengan atlet
func (s *Service) GetByID(id uint) (*LaporanDetail, error) {
	laporan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("laporan pertandingan tidak ditemukan")
	}

	atletA, _ := s.repo.GetAtletBySisi(id, "A")
	atletB, _ := s.repo.GetAtletBySisi(id, "B")
	if atletA == nil {
		atletA = []AtletSisiItem{}
	}
	if atletB == nil {
		atletB = []AtletSisiItem{}
	}
	laporan.AtletA = atletA
	laporan.AtletB = atletB

	return laporan, nil
}

// Create buat laporan baru + insert atlet per sisi
func (s *Service) Create(req *CreateLaporanRequest, userID uint) (*LaporanDetail, error) {
	// Normalisasi ke uppercase agar toleran terhadap input frontend
	req.Babak = strings.ToUpper(req.Babak)
	req.Pemenang = strings.ToUpper(req.Pemenang)

	if err := validateBabak(req.Babak); err != nil {
		return nil, err
	}
	if err := validatePemenang(req.Pemenang); err != nil {
		return nil, err
	}

	var createdBy *uint
	if userID > 0 {
		createdBy = &userID
	}

	laporan := &LaporanPertandingan{
		TanggalPertandingan: TanggalDate{mustParseDate(req.TanggalPertandingan)},
		WaktuPertandingan:   req.WaktuPertandingan,
		Venue:               req.Venue,
		CaborID:             req.CaborID,
		NomorID:             req.NomorID,
		Babak:               req.Babak,
		KontingenAID:        req.KontingenAID,
		KontingenBID:        req.KontingenBID,
		HasilPertandingan:   req.HasilPertandingan,
		Pemenang:            req.Pemenang,
		JuaraKe:             req.JuaraKe,
		Wasit:               req.Wasit,
		CatatanKhusus:       req.CatatanKhusus,
		CreatedBy:           createdBy,
	}

	if err := s.repo.Create(laporan); err != nil {
		return nil, fmt.Errorf("gagal membuat laporan: %w", err)
	}

	// Insert atlet sisi A dan B
	if len(req.AtletA) > 0 {
		if err := s.repo.ReplaceAtletSisi(laporan.ID, "A", req.AtletA); err != nil {
			return nil, fmt.Errorf("gagal menyimpan atlet sisi A: %w", err)
		}
	}
	if len(req.AtletB) > 0 {
		if err := s.repo.ReplaceAtletSisi(laporan.ID, "B", req.AtletB); err != nil {
			return nil, fmt.Errorf("gagal menyimpan atlet sisi B: %w", err)
		}
	}

	return s.GetByID(laporan.ID)
}

// Update ubah data laporan — partial, hanya field yang dikirim
func (s *Service) Update(id uint, req *UpdateLaporanRequest) (*LaporanDetail, error) {
	// Pastikan laporan ada
	if _, err := s.repo.GetByID(id); err != nil {
		return nil, errors.New("laporan pertandingan tidak ditemukan")
	}

	updates := make(map[string]interface{})
	if req.TanggalPertandingan != "" {
		updates["tanggal_pertandingan"] = mustParseDate(req.TanggalPertandingan)
	}
	if req.WaktuPertandingan != "" {
		updates["waktu_pertandingan"] = req.WaktuPertandingan
	}
	if req.Venue != "" {
		updates["venue"] = req.Venue
	}
	if req.CaborID > 0 {
		updates["cabor_id"] = req.CaborID
	}
	if req.NomorID > 0 {
		updates["nomor_id"] = req.NomorID
	}
	if req.Babak != "" {
		req.Babak = strings.ToUpper(req.Babak)
		if err := validateBabak(req.Babak); err != nil {
			return nil, err
		}
		updates["babak"] = req.Babak
	}
	if req.KontingenAID > 0 {
		updates["kontingen_a_id"] = req.KontingenAID
	}
	if req.KontingenBID != nil {
		updates["kontingen_b_id"] = req.KontingenBID
	}
	if req.HasilPertandingan != "" {
		updates["hasil_pertandingan"] = req.HasilPertandingan
	}
	if req.Pemenang != "" {
		req.Pemenang = strings.ToUpper(req.Pemenang)
		if err := validatePemenang(req.Pemenang); err != nil {
			return nil, err
		}
		updates["pemenang"] = req.Pemenang
	}
	if req.JuaraKe != nil {
		updates["juara_ke"] = req.JuaraKe
	}
	if req.Wasit != "" {
		updates["wasit"] = req.Wasit
	}
	if req.CatatanKhusus != nil {
		updates["catatan_khusus"] = req.CatatanKhusus
	}

	if len(updates) > 0 {
		if err := s.repo.Update(id, updates); err != nil {
			return nil, fmt.Errorf("gagal mengupdate laporan: %w", err)
		}
	}

	// Update atlet jika dikirim (replace)
	if req.AtletA != nil {
		if err := s.repo.ReplaceAtletSisi(id, "A", req.AtletA); err != nil {
			return nil, fmt.Errorf("gagal update atlet sisi A: %w", err)
		}
	}
	if req.AtletB != nil {
		if err := s.repo.ReplaceAtletSisi(id, "B", req.AtletB); err != nil {
			return nil, fmt.Errorf("gagal update atlet sisi B: %w", err)
		}
	}

	return s.GetByID(id)
}

// Delete hapus laporan (atlet ikut terhapus via ON DELETE CASCADE)
func (s *Service) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.New("laporan pertandingan tidak ditemukan")
	}
	return s.repo.Delete(id)
}

// UpdateFile update path foto_bukti atau video_bukti
func (s *Service) UpdateFile(id uint, column, path string) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.New("laporan pertandingan tidak ditemukan")
	}
	return s.repo.UpdateFile(id, column, path)
}

// mustParseDate parse "YYYY-MM-DD" ke time.Time. Return zero jika gagal.
func mustParseDate(s string) time.Time {
	for _, layout := range []string{"2006-01-02", time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// ===== VALIDASI ENUM =====

var validBabak = map[string]bool{
	"PENYISIHAN": true, "8_BESAR": true, "PEREMPAT_FINAL": true,
	"SEMIFINAL": true, "FINAL": true, "PEREBUTAN_TEMPAT_3": true, "LAINNYA": true,
}

var validPemenang = map[string]bool{
	"TIM_A": true, "TIM_B": true, "DRAW": true,
}

func validateBabak(b string) error {
	if !validBabak[b] {
		return errors.New("babak tidak valid (PENYISIHAN/8_BESAR/PEREMPAT_FINAL/SEMIFINAL/FINAL/PEREBUTAN_TEMPAT_3/LAINNYA)")
	}
	return nil
}

func validatePemenang(p string) error {
	if !validPemenang[p] {
		return errors.New("pemenang tidak valid (TIM_A/TIM_B/DRAW)")
	}
	return nil
}

// ===== DROPDOWN =====

// GetKontingenDropdown ambil semua kontingen untuk dropdown Tim A/B
func (s *Service) GetKontingenDropdown() ([]KontingenDropdownItem, error) {
	return s.repo.GetKontingenDropdown()
}

// GetCaborDropdown ambil cabor aktif untuk dropdown
func (s *Service) GetCaborDropdown() ([]CaborDropdownItem, error) {
	return s.repo.GetCaborDropdown()
}

// GetNomorDropdown ambil nomor aktif, difilter by cabor_id jika ada
func (s *Service) GetNomorDropdown(caborID uint) ([]NomorDropdownItem, error) {
	return s.repo.GetNomorDropdown(caborID)
}

// GetAtletTerdaftarDropdown ambil atlet dari trx_pendaftaran_atlet
// Filter: kontingen_id, cabor_id, nomor_id — semua opsional
func (s *Service) GetAtletTerdaftarDropdown(kontingenID, caborID, nomorID uint) ([]AtletTerdaftarDropdownItem, error) {
	return s.repo.GetAtletTerdaftarDropdown(kontingenID, caborID, nomorID)
}
