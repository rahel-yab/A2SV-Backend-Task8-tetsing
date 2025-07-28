package Infrastructure

import (
	"task_manager/Domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtService struct {
	secretKey []byte
}

func NewJWTService(secret string) Domain.IJWTService {
	return &jwtService{secretKey: []byte(secret)}
}

func (j *jwtService) GenerateToken(user *Domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}