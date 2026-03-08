package userterritories

import (
	"gorm.io/gorm"
)

type UserTerritory struct {
	UserID      uint `gorm:"primaryKey;column:user_id"`
	TerritoryID uint `gorm:"primaryKey;column:territory_id"`
}

func (UserTerritory) TableName() string {
	return "user_territories"
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]UserTerritory, error) {
	var userTerritories []UserTerritory
	err := r.DB.Find(&userTerritories).Error
	return userTerritories, err
}

func (r *Repository) Create(userTerritory *UserTerritory) error {
	return r.DB.Create(userTerritory).Error
}

func (r *Repository) Delete(userID, territoryID uint) error {
	return r.DB.Where("user_id = ? AND territory_id = ?", userID, territoryID).Delete(&UserTerritory{}).Error
}

func (r *Repository) DeleteByUserID(userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&UserTerritory{}).Error
}

func (r *Repository) DeleteByTerritoryID(territoryID uint) error {
	return r.DB.Where("territory_id = ?", territoryID).Delete(&UserTerritory{}).Error
}

func (r *Repository) GetByUserID(userID uint) ([]UserTerritory, error) {
	var userTerritories []UserTerritory
	err := r.DB.Where("user_id = ?", userID).Find(&userTerritories).Error
	return userTerritories, err
}

func (r *Repository) GetByTerritoryID(territoryID uint) ([]UserTerritory, error) {
	var userTerritories []UserTerritory
	err := r.DB.Where("territory_id = ?", territoryID).Find(&userTerritories).Error
	return userTerritories, err
}

func (r *Repository) GetTerritoryIDsByUserID(userID uint) ([]uint, error) {
	var territoryIDs []uint
	err := r.DB.Table("user_territories").
		Where("user_id = ?", userID).
		Pluck("territory_id", &territoryIDs).Error
	return territoryIDs, err
}

func (r *Repository) GetUserIDsByTerritoryID(territoryID uint) ([]uint, error) {
	var userIDs []uint
	err := r.DB.Table("user_territories").
		Where("territory_id = ?", territoryID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
