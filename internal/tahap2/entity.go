package tahap2

import "time"

type EntryNumber struct {
	ID          uint `gorm:"primaryKey"`
	KontingenID uint `gorm:"index;uniqueIndex:idx_kontingen_nomor"`
	NomorID     uint `gorm:"index;uniqueIndex:idx_kontingen_nomor"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (EntryNumber) TableName() string {
	return "trx_kontingen_nomor"
}

type EventWithStatus struct {
	EventID      uint   `json:"event_id"`
	Cabor        string `json:"cabor"`
	NamaEvent    string `json:"nama_event"`
	JenisKelamin string `json:"jenis_kelamin"`
	Ikut         bool   `json:"ikut"`
}

type Response struct {
	Tahap2Submitted bool              `json:"tahap2_submitted"`
	SubmittedAt     *time.Time        `json:"submitted_at,omitempty"`
	AvailableCabor  []string          `json:"available_cabor"`
	Events          []EventWithStatus `json:"events"`
}
