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

		return nil, err

	}

	if !user.IsActive {

		return nil, errors.New("akun tidak aktif")

	}

	if user.Password != hashPassword(password) {

		return nil, errors.New("password salah")

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

	token, err := jwt.GenerateFullToken(

		user.ID,

		kontingenID,

		userRole,

		user.Email,

		user.Name,

		territoryName, // Gunakan nama territory dari database

		user.Avatar,
	)

	if err != nil {

		return nil, err

	}

	user.Password = "" // jangan pernah kirim password

	return &LoginResponse{

		Token: token,

		User: *user,

		Role: userRole,
	}, nil

}
