package controllers

import (
	"net/http"
	"task_manager/usecases"
	"task_manager/domain"
	"time"

	"github.com/gin-gonic/gin"
)

// UserDTO is a data transfer object for user information.
type UserDTO struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password,omitempty"`
    Role     string `json:"role"`
}

// todomainUser converts a UserDTO to a domain.User.
func todomainUser(dto *UserDTO) *domain.User {
    return &domain.User{
        ID:       dto.ID,
        Username: dto.Username,
        Email:    dto.Email,
        Password: dto.Password,
        Role:     dto.Role,
    }
}

// toUserDTO converts a domain.User to a UserDTO.
func toUserDTO(user *domain.User) *UserDTO {
    return &UserDTO{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
        Role:     user.Role,
    }
}

// TaskDTO is a data transfer object for task information.
type TaskDTO struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    DueDate     time.Time `json:"due_date"`
    Status      string    `json:"status"`
}

// todomainTask converts a TaskDTO to a domain.Task.
func todomainTask(dto *TaskDTO) *domain.Task {
    return &domain.Task{
        ID:          dto.ID,
        Title:       dto.Title,
        Description: dto.Description,
        DueDate:     dto.DueDate,
        Status:      dto.Status,
    }
}

// toTaskDTO converts a domain.Task to a TaskDTO.
func toTaskDTO(task *domain.Task) *TaskDTO {
    return &TaskDTO{
        ID:          task.ID,
        Title:       task.Title,
        Description: task.Description,
        DueDate:     task.DueDate,
        Status:      task.Status,
    }
}

// UserController handles user-related HTTP requests.
type UserController struct {
    userUsecase *usecases.UserUsecase
}

// NewUserController creates a new UserController.
func NewUserController(userUsecase *usecases.UserUsecase) *UserController {
    return &UserController{userUsecase: userUsecase}
}

// RegisterUser handles user registration requests.
func (ctrl *UserController) RegisterUser(c *gin.Context) {
    var req UserDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
    role, err := ctrl.userUsecase.RegisterUser(c.Request.Context(), req.Username, req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "role": role})
}

// LoginUser handles user login requests.
func (ctrl *UserController) LoginUser(c *gin.Context) {
    var req UserDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
    usernameOrEmail := req.Email
    if usernameOrEmail == "" {
        usernameOrEmail = req.Username
    }
    token, role, err := ctrl.userUsecase.LoginUser(c.Request.Context(), usernameOrEmail, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token, "role": role})
}

// PromoteUser promotes a user to admin.
func (ctrl *UserController) PromoteUser(c *gin.Context) {
    var req struct {
        Identifier string `json:"identifier"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.Identifier == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email is required"})
        return
    }
    err := ctrl.userUsecase.PromoteUserToAdmin(c.Request.Context(), req.Identifier)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

// TaskController handles task-related HTTP requests.
type TaskController struct {
    taskUsecase *usecases.TaskUsecase
}

// NewTaskController creates a new TaskController.
func NewTaskController(taskUsecase *usecases.TaskUsecase) *TaskController {
    return &TaskController{taskUsecase: taskUsecase}
}

// GetTasks returns all tasks.
func (ctrl *TaskController) GetTasks(c *gin.Context) {
    tasks, err := ctrl.taskUsecase.GetAllTasks(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    var dtos []TaskDTO
    for _, t := range tasks {
        dtos = append(dtos, *toTaskDTO(&t))
    }
    c.JSON(http.StatusOK, dtos)
}

// GetTask returns a task by ID.
func (ctrl *TaskController) GetTask(c *gin.Context) {
    id := c.Param("id")
    task, err := ctrl.taskUsecase.GetTaskByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
        return
    }
    c.JSON(http.StatusOK, toTaskDTO(task))
}

// AddTask adds a new task.
func (ctrl *TaskController) AddTask(c *gin.Context) {
    var dto TaskDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    task := todomainTask(&dto)
    if err := ctrl.taskUsecase.Create(c.Request.Context(), task); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

// UpdateTask updates an existing task.
func (ctrl *TaskController) UpdateTask(c *gin.Context) {
    var dto TaskDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    task := todomainTask(&dto)
    if err := ctrl.taskUsecase.UpdateTask(c.Request.Context(), task); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// RemoveTask deletes a task by ID.
func (ctrl *TaskController) RemoveTask(c *gin.Context) {
    id := c.Param("id")
    if err := ctrl.taskUsecase.DeleteTask(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Task removed"})
}

