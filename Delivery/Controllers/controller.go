package Controllers

import (
	"net/http"
	"task_manager/Domain"
	"task_manager/Usecases"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

func toDomainUser(dto *UserDTO) *Domain.User {
	return &Domain.User{
		ID:       dto.ID,
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		Role:     dto.Role,
	}
}

func toUserDTO(user *Domain.User) *UserDTO {
	return &UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}

type TaskDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

func toDomainTask(dto *TaskDTO) *Domain.Task {
	return &Domain.Task{
		ID:          dto.ID,
		Title:       dto.Title,
		Description: dto.Description,
		DueDate:     dto.DueDate,
		Status:      dto.Status,
	}
}

func toTaskDTO(task *Domain.Task) *TaskDTO {
	return &TaskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}
}

// --- UserController ---

type UserController struct {
	userUsecase *Usecases.UserUsecase
}

func NewUserController(userUsecase *Usecases.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

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

// --- TaskController ---

type TaskController struct {
	taskUsecase *Usecases.TaskUsecase
}

func NewTaskController(taskUsecase *Usecases.TaskUsecase) *TaskController {
	return &TaskController{taskUsecase: taskUsecase}
}

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

func (ctrl *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := ctrl.taskUsecase.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, toTaskDTO(task))
}

func (ctrl *TaskController) AddTask(c *gin.Context) {
	var dto TaskDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := toDomainTask(&dto)
	if err := ctrl.taskUsecase.Create(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func (ctrl *TaskController) UpdateTask(c *gin.Context) {
	var dto TaskDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := toDomainTask(&dto)
	if err := ctrl.taskUsecase.UpdateTask(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (ctrl *TaskController) RemoveTask(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.taskUsecase.DeleteTask(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task removed"})
}


