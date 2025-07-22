package Controllers
import (
	"net/http"
    "task_manager/Usecases"
	"task_manager/Domain"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase *Usecases.UserUsecase
}

type TaskController struct {
	TaskUsecase *Usecases.TaskUsecase
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
	role, err := ctrl.UserUsecase.RegisterUser(req.Username, req.Email, req.Password)
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
	token, role, err := ctrl.UserUsecase.LoginUser(usernameOrEmail, req.Password)
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
	err := ctrl.UserUsecase.PromoteUserToAdmin(req.Identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.(error).Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

func (ctrl *TaskController) GetTasks(c *gin.Context) {
	tasks, err := ctrl.TaskUsecase.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (ctrl *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := ctrl.TaskUsecase.GetTaskByID(id)
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
	if err := ctrl.TaskUsecase.AddTask(newTask); err != nil {
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
	if err := ctrl.TaskUsecase.UpdateTask(updatedTask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (ctrl *TaskController) RemoveTask(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.TaskUsecase.DeleteTask(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task removed"})
}
