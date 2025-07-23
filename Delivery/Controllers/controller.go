package Controllers

import (
	"net/http"
	"task_manager/Domain"
	"task_manager/Usecases"
	"github.com/gin-gonic/gin"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userUsecase *Usecases.UserUsecase
}

// NewUserController creates a new UserController
func NewUserController(userUsecase *Usecases.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (ctrl *UserController) RegisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
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
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
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

// TaskController handles task-related HTTP requests
type TaskController struct {
	taskUsecase *Usecases.TaskUsecase
}

// NewTaskController creates a new TaskController
func NewTaskController(taskUsecase *Usecases.TaskUsecase) *TaskController {
	return &TaskController{taskUsecase: taskUsecase}
}

func (ctrl *TaskController) GetTasks(c *gin.Context) {
	tasks, err := ctrl.taskUsecase.GetAllTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (ctrl *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := ctrl.taskUsecase.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (ctrl *TaskController) AddTask(c *gin.Context) {
	var newTask Domain.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.taskUsecase.Create(c.Request.Context(), &newTask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func (ctrl *TaskController) UpdateTask(c *gin.Context) {
	var updatedTask Domain.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.taskUsecase.UpdateTask(c.Request.Context(), &updatedTask); err != nil {
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