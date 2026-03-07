package tahap3

import "time"

type Atlet struct {
	ID          uint      `gorm:"primaryKey"`
	KontingenID uint      `gorm:"index"`
	CaborID     uint      `gorm:"column:cabor_id"`
	Nama        string    `gorm:"column:nama"`
	Gender      string    `gorm:"column:gender"`
	TglLahir    time.Time `gorm:"column:tgl_lahir"`
	NISN        string    `gorm:"column:nisn"`
	SekolahID   uint      `gorm:"column:sekolah_id"`
	BPJS        string    `gorm:"column:bpjs"`
	Kelas       string    `gorm:"column:kelas"`
	Domisili    string    `gorm:"column:domisili"`
	IsBinaan    bool      `gorm:"column:is_binaan"`

	AkteKelahiran   string `gorm:"column:akte_kelahiran"`
	RaporTerakhir   string `gorm:"column:rapor_terakhir"`
	IjazahSttb      string `gorm:"column:ijazah_sttb"`
	BuktiNisn       string `gorm:"column:bukti_nisn"`
	SuratSekolah    string `gorm:"column:surat_sekolah"`
	SuratSehat      string `gorm:"column:surat_sehat"`
	PasFoto         string `gorm:"column:pas_foto"`
	SuratPernyataan string `gorm:"column:surat_pernyataan"`
	SkBinaan        string `gorm:"column:sk_binaan"`
	PaktaIntegritas string `gorm:"column:pakta_integritas"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Atlet) TableName() string { return "master_atlet" }

type Pelatih struct {
	ID          string    `gorm:"primaryKey"`
	KontingenID uint      `gorm:"index"`
	CaborID     uint      `gorm:"column:cabor_id"`
	Nama        string    `gorm:"column:nama"`
	TglLahir    time.Time `gorm:"column:tgl_lahir"`
	Gender      string    `gorm:"column:gender"`
	NIK         string    `gorm:"column:nik"`
	Kontak      string    `gorm:"column:kontak"`
	Provinsi    string    `gorm:"column:provinsi"`
	Pendidikan  string    `gorm:"column:pendidikan"`
	FokusTim    string    `gorm:"column:fokus_tim"`
	Email       string    `gorm:"column:email"`
	KTP         string    `gorm:"column:ktp"`
	Foto        string    `gorm:"column:foto"`
	SuratSehat  string    `gorm:"column:surat_sehat"`
	SkKontingen string    `gorm:"column:sk_kontingen"`
	SuratTugas  string    `gorm:"column:surat_tugas"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Pelatih) TableName() string { return "master_pelatih" }

type Official struct {
	ID          string    `gorm:"primaryKey"`
	KontingenID uint      `gorm:"index"`
	Nama        string    `gorm:"column:nama"`
	Jabatan     string    `gorm:"column:jabatan"`
	TglLahir    time.Time `gorm:"column:tgl_lahir"`
	Gender      string    `gorm:"column:gender"`
	NIK         string    `gorm:"column:nik"`
	Kontak      string    `gorm:"column:kontak"`
	Provinsi    string    `gorm:"column:provinsi"`
	Pendidikan  string    `gorm:"column:pendidikan"`
	Email       string    `gorm:"column:email"`
	KTP         string    `gorm:"column:ktp"`
	SuratSehat  string    `gorm:"column:surat_sehat"`
	Foto        string    `gorm:"column:foto"`
	SkKontingen string    `gorm:"column:sk_kontingen"`
	SuratTugas  string    `gorm:"column:surat_tugas"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Official) TableName() string { return "master_official" }

type TrxPendaftaranAtlet struct {
	ID        uint      `gorm:"primaryKey"`
	AtletID   uint      `gorm:"column:atlet_id"`
	NomorID   uint      `gorm:"column:nomor_id"`
	KelasID   *uint     `gorm:"column:kelas_id"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (TrxPendaftaranAtlet) TableName() string {
	return "trx_pendaftaran_atlet"
}

type Response struct {
	Tahap3Submitted bool       `json:"tahap3_submitted"`
	SubmittedAt     *time.Time `json:"submitted_at,omitempty"`
	Atlets          []Atlet    `json:"atlets"`
	Pelatihs        []Pelatih  `json:"pelatihs"`
	Officials       []Official `json:"officials"`
}
