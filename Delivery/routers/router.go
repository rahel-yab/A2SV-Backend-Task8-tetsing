package routers

import (
	"task_manager/Delivery/Controllers"
	"task_manager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *Controllers.UserController, taskController *Controllers.TaskController, jwtSecret []byte) *gin.Engine {
	router := gin.Default()

	taskGroup := router.Group("/tasks", Infrastructure.AuthMiddleware(jwtSecret))
	{
		taskGroup.GET("", taskController.GetTasks)
		taskGroup.GET(":id", taskController.GetTask)
		taskGroup.DELETE(":id", Infrastructure.AdminOnly(), taskController.RemoveTask)
		taskGroup.PUT(":id", Infrastructure.AdminOnly(), taskController.UpdateTask)
		taskGroup.POST("", Infrastructure.AdminOnly(), taskController.AddTask)
	}

	router.POST("/register", userController.RegisterUser)
	router.POST("/login", userController.LoginUser)

	// Protected route for promoting users
	router.POST("/promote", Infrastructure.AuthMiddleware(jwtSecret), Infrastructure.AdminOnly(), userController.PromoteUser)

	return router
}