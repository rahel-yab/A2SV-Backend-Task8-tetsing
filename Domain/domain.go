package domain

import (
	"context"
	"time"
)

type User struct {
	ID       string 
	Username string 
	Email    string
	Password string 
	Role     string 
}

type Task struct {
	ID          string    
	Title       string    
	Description string    
	DueDate     time.Time 
	Status      string    
}

type ITaskRepository interface {
	AddTask(ctx context.Context, task *Task) error
	GetAllTasks(ctx context.Context) ([]Task, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id string) error
}

type IUserRepository interface {
	AddUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	IsUsersCollectionEmpty(ctx context.Context) (bool, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	UserExistsByUsername(ctx context.Context, username string) (bool, error)
	PromoteUserToAdmin(ctx context.Context, identifier string) error
}

type IPasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type IJWTService interface {
	GenerateToken(user *User) (string, error)
}