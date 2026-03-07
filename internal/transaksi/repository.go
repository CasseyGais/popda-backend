package transaksi

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// TrxKontingenCabor operations
func (r *Repository) CreateTrxKontingenCabor(trx *TrxKontingenCabor) error {
	return r.db.Create(trx).Error
}

func (r *Repository) GetTrxKontingenCabor(kontingenID uint) ([]TrxKontingenCabor, error) {
	var trx []TrxKontingenCabor
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&trx).Error
	return trx, err
}

func (r *Repository) UpdateTrxKontingenCabor(kontingenID, caborID uint, trx *TrxKontingenCabor) error {
	return r.db.Where("kontingen_id = ? AND cabor_id = ?", kontingenID, caborID).Updates(trx).Error
}

// TrxKontingenNomor operations
func (r *Repository) CreateTrxKontingenNomor(trx *TrxKontingenNomor) error {
	return r.db.Create(trx).Error
}

func (r *Repository) GetTrxKontingenNomor(kontingenID uint) ([]TrxKontingenNomor, error) {
	var trx []TrxKontingenNomor
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&trx).Error
	return trx, err
}

func (r *Repository) DeleteTrxKontingenNomor(kontingenID, nomorID uint) error {
	return r.db.Where("kontingen_id = ? AND nomor_id = ?", kontingenID, nomorID).Delete(&TrxKontingenNomor{}).Error
}

// TrxPendaftaranAtlet operations
func (r *Repository) CreateTrxPendaftaranAtlet(trx *TrxPendaftaranAtlet) error {
	return r.db.Create(trx).Error
}

func (r *Repository) GetTrxPendaftaranAtlet(kontingenID uint) ([]TrxPendaftaranAtlet, error) {
	var trx []TrxPendaftaranAtlet
	err := r.db.Where("kontingen_id = ?", kontingenID).Find(&trx).Error
	return trx, err
}

func (r *Repository) UpdateTrxPendaftaranAtlet(atletID, nomorID uint, trx *TrxPendaftaranAtlet) error {
	return r.db.Where("atlet_id = ? AND nomor_id = ?", atletID, nomorID).Updates(trx).Error
}
