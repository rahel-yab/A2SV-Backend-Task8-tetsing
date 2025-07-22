package Domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

type Task struct {
	ID          string `json:"id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Status      string `json:"status" bson:"status"`
}

type TaskRepository interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id string) (*Task, error)
	AddTask(task Task) error
	UpdateTask(task Task) error
	DeleteTask(id string) error
	FetchByTaskID(ctx context.Context, taskID primitive.ObjectID) (Task, error)
}

type UserRepository interface {
	AddUser(user User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	IsUsersCollectionEmpty() (bool, error)
	UserExistsByEmail(email string) (bool, error)
	UserExistsByUsername(username string) (bool, error)
	PromoteUserToAdmin(identifier string) error
}