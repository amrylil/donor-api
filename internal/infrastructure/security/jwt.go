package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService menangani pembuatan dan validasi token
type JWTService struct {
	secretKey       string
	expirationHours int64
}

// NewJWTService membuat instance baru dari JWTService
func NewJWTService(secretKey string, expirationHours int64) *JWTService {
	return &JWTService{
		secretKey:       secretKey,
		expirationHours: expirationHours,
	}
}

// GenerateToken membuat token JWT baru untuk ID pengguna
func (s *JWTService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(), // <-- Simpan sebagai string
		"exp": time.Now().Add(time.Hour * time.Duration(s.expirationHours)).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken memvalidasi token dan mengembalikan klaim jika valid
func (s *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
