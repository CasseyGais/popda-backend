package auth

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {

	return &Repository{db: db}

}

// ================= FIND USER =================

func (r *Repository) FindByEmail(email string) (*User, error) {

	var user User

	err := r.db.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, errors.New("user tidak ditemukan")

		}

		return nil, err

	}

	return &user, nil

}

// ================= GET USER ROLE =================

func (r *Repository) GetUserRole(userID uint) (string, error) {

	var roleName string

	err := r.db.
		Table("user_roles ur").
		Select("r.name").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Limit(1).
		Scan(&roleName).Error

	if err != nil {

		return "", err

	}

	if roleName == "" {

		return "user", nil

	}

	return roleName, nil

}

// ================= GET KONTINGEN ID =================

// Alur:

// users

// → user_territories

// → territories

// → kontingen

func (r *Repository) GetKontingenIDByUser(userID uint) (uint, error) {

	var kontingenID uint

	err := r.db.
		Table("user_territories ut").
		Select("k.id").
		Joins("JOIN territories t ON t.id = ut.territory_id").
		Joins("JOIN kontingen k ON k.territory_id = t.id").
		Where("ut.user_id = ?", userID).
		Limit(1).
		Scan(&kontingenID).Error

	if err != nil {

		return 0, err

	}

	if kontingenID == 0 {

		return 0, errors.New("kontingen tidak ditemukan untuk user ini")

	}

	return kontingenID, nil

}

// ================= GET TERRITORY NAME =================

// users → user_territories → territories

func (r *Repository) GetTerritoryNameByUser(userID uint) (string, error) {

	var territoryName string

	err := r.db.
		Table("user_territories ut").
		Select("t.name").
		Joins("JOIN territories t ON t.id = ut.territory_id").
		Where("ut.user_id = ?", userID).
		Limit(1).
		Scan(&territoryName).Error

	if err != nil {

		return "", err

	}

	return territoryName, nil

}

// ================= DEBUG METHODS =================

func (r *Repository) DebugUserTerritory(userID uint) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	// Cek user di user_territories

	var territoryID uint

	err := r.db.Table("user_territories").
		Select("territory_id").
		Where("user_id = ?", userID).
		First(&territoryID).Error

	if err != nil {

		result["user_territories"] = "NOT FOUND: " + err.Error()

	} else {

		result["user_territories"] = map[string]interface{}{

			"user_id": userID,

			"territory_id": territoryID,
		}

	}

	// Cek territory

	if territoryID > 0 {

		var territory struct {
			ID uint

			Name string
		}

		err := r.db.Table("territories").
			Select("id, name").
			Where("id = ?", territoryID).
			First(&territory).Error

		if err != nil {

			result["territories"] = "NOT FOUND: " + err.Error()

		} else {

			result["territories"] = territory

		}

		// Cek kontingen

		var kontingen struct {
			ID uint

			TerritoryID uint

			NamaKontingen string
		}

		err = r.db.Table("kontingen").
			Select("id, territory_id, nama_kontingen").
			Where("territory_id = ?", territoryID).
			First(&kontingen).Error

		if err != nil {

			result["kontingen"] = "NOT FOUND: " + err.Error()

		} else {

			result["kontingen"] = kontingen

		}

	}

	return result, nil

}
