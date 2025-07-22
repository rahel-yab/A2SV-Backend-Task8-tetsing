package Usecases

import (
	"errors"
	"task_manager/Domain"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

type mockUserRepo struct {
	users map[string]Domain.User
	promoted map[string]bool
}

func (m *mockUserRepo) AddUser(user Domain.User) error {
	m.users[user.Email] = user
	return nil
}
func (m *mockUserRepo) GetUserByEmail(email string) (*Domain.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return &user, nil
}
func (m *mockUserRepo) GetUserByUsername(username string) (*Domain.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockUserRepo) IsUsersCollectionEmpty() (bool, error) {
	return len(m.users) == 0, nil
}
func (m *mockUserRepo) UserExistsByEmail(email string) (bool, error) {
	_, ok := m.users[email]
	return ok, nil
}
func (m *mockUserRepo) UserExistsByUsername(username string) (bool, error) {
	for _, user := range m.users {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}
func (m *mockUserRepo) PromoteUserToAdmin(identifier string) error {
	m.promoted[identifier] = true
	return nil
}

type mockPasswordService struct{}
func (m *mockPasswordService) HashPassword(password string) (string, error) { return password, nil }
func (m *mockPasswordService) CheckPasswordHash(password, hash string) bool { return password == hash }

type mockJWTService struct{}
func (m *mockJWTService) GenerateToken(user *Domain.User) (string, error) { return "token", nil }
func (m *mockJWTService) ValidateToken(tokenString string) (*jwt.Token, error) { return nil, nil }

func TestRegisterUser(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]Domain.User), promoted: make(map[string]bool)}
	uc := &UserUsecase{Repo: repo, PasswordService: &mockPasswordService{}, JWTService: &mockJWTService{}}
	role, err := uc.RegisterUser("testuser", "test@example.com", "pass")
	if err != nil || role != "admin" {
		t.Errorf("expected admin, got %v, err %v", role, err)
	}
	// Register again, should get user
	role, err = uc.RegisterUser("testuser2", "test2@example.com", "pass")
	if err != nil || role != "user" {
		t.Errorf("expected user, got %v, err %v", role, err)
	}
}

func TestLoginUser(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]Domain.User), promoted: make(map[string]bool)}
	uc := &UserUsecase{Repo: repo, PasswordService: &mockPasswordService{}, JWTService: &mockJWTService{}}
	uc.RegisterUser("testuser", "test@example.com", "pass")
	_, _, err := uc.LoginUser("test@example.com", "pass")
	if err != nil {
		t.Errorf("expected login success, got err %v", err)
	}
	_, _, err = uc.LoginUser("wrong@example.com", "pass")
	if err == nil {
		t.Errorf("expected login fail for wrong email")
	}
}

func TestPromoteUserToAdmin(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]Domain.User), promoted: make(map[string]bool)}
	uc := &UserUsecase{Repo: repo, PasswordService: &mockPasswordService{}, JWTService: &mockJWTService{}}
	uc.RegisterUser("testuser", "test@example.com", "pass")
	err := repo.PromoteUserToAdmin("testuser")
	if err != nil {
		t.Errorf("expected promote success, got err %v", err)
	}
	if !repo.promoted["testuser"] {
		t.Errorf("expected user to be promoted")
	}
} 