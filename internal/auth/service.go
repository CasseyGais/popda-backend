package auth

import (
	"crypto/sha256"

	"encoding/hex"

	"errors"

	"popda_bulutangkis/pkg/jwt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {

	return &Service{repo: repo}

}

func hashPassword(password string) string {

	hash := sha256.Sum256([]byte(password))

	return hex.EncodeToString(hash[:])

}

type LoginResponse struct {
	Token string `json:"token"`

	User User `json:"user"`

	Role string `json:"role"`
}

func (s *Service) Login(identifier, password string) (*LoginResponse, error) {

	user, err := s.repo.FindByEmail(identifier)

	if err != nil {
		// Jangan bocorkan apakah email ditemukan atau tidak
		return nil, errors.New("Email atau password anda salah")
	}

	if !user.IsActive {

		return nil, errors.New("akun tidak aktif")

	}

	if user.Password != hashPassword(password) {
		// Pesan sama untuk email tidak ditemukan vs password salah
		return nil, errors.New("Email atau password anda salah")

	}

	kontingenID, err := s.repo.GetKontingenIDByUser(user.ID)

	if err != nil {

		return nil, err

	}

	userRole, err := s.repo.GetUserRole(user.ID)

	if err != nil {

		return nil, err

	}

	territoryName, err := s.repo.GetTerritoryNameByUser(user.ID)

	if err != nil {

		return nil, err

	}

	userTerritories, err := s.repo.GetTerritoriesByUser(user.ID)

	if err != nil {

		return nil, err

	}

	token, err := jwt.GenerateFullToken(

		user.ID,

		kontingenID,

		userRole,

		user.Email,

		user.Name,

		territoryName,

		user.Avatar,
	)

	if err != nil {

		return nil, err

	}

	user.Password = "" // jangan pernah kirim password

	user.Territories = userTerritories

	return &LoginResponse{

		Token: token,

		User: *user,

		Role: userRole,
	}, nil

}

func (s *Service) Logout(userID interface{}) error {
	// For stateless JWT, logout is handled client-side by deleting the token
	// This method can be extended to implement token blacklisting if needed
	return nil
}
