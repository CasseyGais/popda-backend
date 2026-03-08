package users

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]User, error) {
	var users []User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *Repository) GetByID(id uint) (*User, error) {
	var user User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	var user User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Create(user *User) error {
	return r.DB.Create(user).Error
}

func (r *Repository) Update(user *User) error {
	return r.DB.Save(user).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&User{}, id).Error
}

func (r *Repository) UpdateAvatar(id uint, avatar string) error {
	return r.DB.Model(&User{}).Where("id = ?", id).Update("avatar", avatar).Error
}

func (r *Repository) UpdatePassword(id uint, password string) error {
	return r.DB.Model(&User{}).Where("id = ?", id).Update("password", password).Error
}

func (r *Repository) UpdateStatus(id uint, isActive bool) error {
	return r.DB.Model(&User{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *Repository) GetRoles(userID uint) ([]uint, error) {
	var roleIDs []uint
	err := r.DB.Table("user_roles").
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

func (r *Repository) AssignRole(userID, roleID uint) error {
	assignment := map[string]interface{}{
		"user_id": userID,
		"role_id": roleID,
	}
	return r.DB.Table("user_roles").Create(assignment).Error
}

func (r *Repository) RemoveRole(userID, roleID uint) error {
	return r.DB.Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&map[string]interface{}{}).Error
}

func (r *Repository) GetTerritories(userID uint) ([]uint, error) {
	var territoryIDs []uint
	err := r.DB.Table("user_territories").
		Where("user_id = ?", userID).
		Pluck("territory_id", &territoryIDs).Error
	return territoryIDs, err
}

func (r *Repository) AssignTerritory(userID, territoryID uint) error {
	assignment := map[string]interface{}{
		"user_id":       userID,
		"territory_id": territoryID,
	}
	return r.DB.Table("user_territories").Create(assignment).Error
}

func (r *Repository) RemoveTerritory(userID, territoryID uint) error {
	return r.DB.Table("user_territories").
		Where("user_id = ? AND territory_id = ?", userID, territoryID).
		Delete(&map[string]interface{}{}).Error
}
