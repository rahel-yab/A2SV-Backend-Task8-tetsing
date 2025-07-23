package Usecases

import (
	"context"
	"errors"
	"task_manager/Domain"
	"task_manager/Infrastructure"
	"time"
)

type UserUsecase struct {
	userRepository Domain.UserRepository
	passwordService Infrastructure.PasswordService
	jwtService Infrastructure.JWTService
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository Domain.UserRepository, passwordService Infrastructure.PasswordService, jwtService Infrastructure.JWTService, timeout time.Duration) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
		passwordService: passwordService,
		jwtService: jwtService,
		contextTimeout: timeout,
	}
}

func (uu *UserUsecase) RegisterUser(ctx context.Context, username, email, password string) (string, error) {
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
	user := &Domain.User{
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
	c, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()
	var user *Domain.User
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
	c, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.PromoteUserToAdmin(c, identifier)
}