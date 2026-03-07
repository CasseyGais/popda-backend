package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      uint   `json:"user_id"`
	KontingenID uint   `json:"kontingen_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	KabKota     string `json:"kab_kota"`
	FotoProfil  string `json:"foto_profil"`
	jwt.RegisteredClaims
}

type UserClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

// Get JWT Secret from environment
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "rahasiasuperkuat1234567890" // fallback
	}
	return []byte(secret)
}

// Simplified token generation for basic auth
func GenerateToken(claims UserClaims) (string, error) {
	tokenClaims := Claims{
		UserID: claims.UserID,
		Role:   claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString(getJWTSecret())
}

// Full token generation with all user details
func GenerateFullToken(
	userID uint,
	kontingenID uint,
	role string,
	email string,
	name string,
	kabKota string,
	fotoProfil string,
) (string, error) {

	claims := Claims{
		UserID:      userID,
		KontingenID: kontingenID,
		Role:        role,
		Email:       email,
		Username:    email, // username sementara pakai email
		Name:        name,
		KabKota:     kabKota,
		FotoProfil:  fotoProfil,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
