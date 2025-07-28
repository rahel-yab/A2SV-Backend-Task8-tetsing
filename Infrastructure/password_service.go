package Infrastructure

import (
	"task_manager/Domain"

	"golang.org/x/crypto/bcrypt"
)

type passwordService struct{}

func NewPasswordService() Domain.IPasswordService {
	return &passwordService{}
}

func (p *passwordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (p *passwordService) CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}