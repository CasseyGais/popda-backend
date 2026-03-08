package territories

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]Territory, error) {
	var territories []Territory
	err := r.DB.Find(&territories).Error
	return territories, err
}

func (r *Repository) GetByID(id uint) (*Territory, error) {
	var territory Territory
	err := r.DB.First(&territory, id).Error
	if err != nil {
		return nil, err
	}
	return &territory, nil
}

func (r *Repository) GetByType(territoryType string) ([]Territory, error) {
	var territories []Territory
	err := r.DB.Where("type = ?", territoryType).Find(&territories).Error
	return territories, err
}

func (r *Repository) GetProvinces() ([]Territory, error) {
	return r.GetByType("PROVINSI")
}

func (r *Repository) GetKabupatens() ([]Territory, error) {
	return r.GetByType("KABUPATEN")
}

func (r *Repository) GetKotas() ([]Territory, error) {
	return r.GetByType("KOTA")
}

func (r *Repository) Create(territory *Territory) error {
	return r.DB.Create(territory).Error
}

func (r *Repository) Update(territory *Territory) error {
	return r.DB.Save(territory).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&Territory{}, id).Error
}

func (r *Repository) GetByUserID(userID uint) ([]Territory, error) {
	var territories []Territory
	err := r.DB.Table("territories").
		Joins("INNER JOIN user_territories ON territories.id = user_territories.territory_id").
		Where("user_territories.user_id = ?", userID).
		Find(&territories).Error
	return territories, err
}

func (r *Repository) AssignToUser(userID, territoryID uint) error {
	assignment := map[string]interface{}{
		"user_id":       userID,
		"territory_id": territoryID,
	}
	return r.DB.Table("user_territories").Create(assignment).Error
}

func (r *Repository) RemoveFromUser(userID, territoryID uint) error {
	return r.DB.Table("user_territories").
		Where("user_id = ? AND territory_id = ?", userID, territoryID).
		Delete(&map[string]interface{}{}).Error
}
