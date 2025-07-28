package infrastructure

import (
	"fmt"
	"task_manager/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtService struct {
	secretKey []byte
}

func NewJWTService(secret string) domain.IJWTService {
	return &jwtService{secretKey: []byte(secret)}
}

func (j *jwtService) GenerateToken(user *domain.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"iat":      now.Unix(),
		"exp":      now.Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}