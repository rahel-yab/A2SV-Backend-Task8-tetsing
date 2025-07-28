package usecases

import (
	"context"
	"errors"
	"strings"
	"task_manager/domain"
	"time"
)

type UserUsecase struct {
	userRepository domain.IUserRepository
	passwordService domain.IPasswordService
	jwtService domain.IJWTService
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository domain.IUserRepository, passwordService domain.IPasswordService, jwtService domain.IJWTService, timeout time.Duration) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
		passwordService: passwordService,
		jwtService: jwtService,
		contextTimeout: timeout,
	}
}

func (uu *UserUsecase) RegisterUser(ctx context.Context, username, email, password string) (string, error) {
	// Validate input parameters
	if username == "" {
		return "", errors.New("username is required")
	}
	if email == "" {
		return "", errors.New("email is required")
	}
	if password == "" {
		return "", errors.New("password is required")
	}
	if len(password) < 6 {
		return "", errors.New("password must be at least 6 characters long")
	}
	
	// Basic email validation
	if !strings.Contains(email, "@") {
		return "", errors.New("invalid email format")
	}
	
	c, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()
	if exists, _ := uu.userRepository.UserExistsByEmail(c, email); exists {
		return "", errors.New("email already registered")
	}
	if exists, _ := uu.userRepository.UserExistsByUsername(c, username); exists {
		return "", errors.New("username already taken")
	}
	isEmpty, _ := uu.userRepository.IsUsersCollectionEmpty(c)
	role := "user"
	if isEmpty {
		role = "admin"
	}
	hashed, err := uu.passwordService.HashPassword(password)
	if err != nil {
		return "", err
	}
	user := &domain.User{
		Username: username,
		Email:    email,
		Password: hashed,
		Role:     role,
	}
	if err := uu.userRepository.AddUser(c, user); err != nil {
		return "", err
	}
	return role, nil
}

func (uu *UserUsecase) LoginUser(ctx context.Context, usernameOrEmail, password string) (string, string, error) {
	// Validate input parameters
	if usernameOrEmail == "" {
		return "", "", errors.New("username or email is required")
	}
	if password == "" {
		return "", "", errors.New("password is required")
	}
	
	c, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()
	var user *domain.User
	var err error
	if user, err = uu.userRepository.GetUserByEmail(c, usernameOrEmail); err != nil || user == nil {
		user, err = uu.userRepository.GetUserByUsername(c, usernameOrEmail)
		if err != nil || user == nil {
			return "", "", errors.New("invalid email/username or password")
		}
	}
	if !uu.passwordService.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid email/username or password")
	}
	token, err := uu.jwtService.GenerateToken(user)
	if err != nil {
		return "", "", err
	}
	return token, user.Role, nil
}

func (uu *UserUsecase) PromoteUserToAdmin(ctx context.Context, identifier string) error {
	// Validate input parameters
	if identifier == "" {
		return errors.New("identifier is required")
	}
	
	c, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.PromoteUserToAdmin(c, identifier)
}
