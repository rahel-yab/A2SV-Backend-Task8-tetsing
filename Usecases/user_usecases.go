package Usecases

import (
	"errors"
	"task_manager/Domain"
	"task_manager/Infrastructure"
	"task_manager/Repositories"
)

type UserUsecase struct {
	Repo            Repositories.UserRepository
	PasswordService Infrastructure.PasswordService
	JWTService      Infrastructure.JWTService
}

func (u *UserUsecase) PromoteUserToAdmin(identifier string) any {
	panic("unimplemented")
}

func (u *UserUsecase) RegisterUser(username, email, password string) (string, error) {
	if exists, _ := u.Repo.UserExistsByEmail(email); exists {
		return "", errors.New("email already registered")
	}
	if exists, _ := u.Repo.UserExistsByUsername(username); exists {
		return "", errors.New("username already taken")
	}
	isEmpty, _ := u.Repo.IsUsersCollectionEmpty()
	role := "user"
	if isEmpty {
		role = "admin"
	}
	hashed, err := u.PasswordService.HashPassword(password)
	if err != nil {
		return "", err
	}
	user := Domain.User{
		Username: username,
		Email:    email,
		Password: hashed,
		Role:     role,
	}
	if err := u.Repo.AddUser(user); err != nil {
		return "", err
	}
	return role, nil
}

func (u *UserUsecase) LoginUser(usernameOrEmail, password string) (string, string, error) {
	var user *Domain.User
	var err error
	if user, err = u.Repo.GetUserByEmail(usernameOrEmail); err != nil || user == nil {
		user, err = u.Repo.GetUserByUsername(usernameOrEmail)
		if err != nil || user == nil {
			return "", "", errors.New("invalid email/username or password")
		}
	}
	if !u.PasswordService.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid email/username or password")
	}
	token, err := u.JWTService.GenerateToken(user)
	if err != nil {
		return "", "", err
	}
	return token, user.Role, nil
}
