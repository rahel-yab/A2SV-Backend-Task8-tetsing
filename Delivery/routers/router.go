package routers

import (
	"task_manager/delivery/controllers"
	"task_manager/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *controllers.UserController, taskController *controllers.TaskController, jwtSecret []byte) *gin.Engine {
	router := gin.Default()

	taskGroup := router.Group("/tasks", infrastructure.AuthMiddleware(jwtSecret))
	{
		taskGroup.GET("", taskController.GetTasks)
		taskGroup.GET(":id", taskController.GetTask)
		taskGroup.DELETE(":id", infrastructure.AdminOnly(), taskController.RemoveTask)
		taskGroup.PUT(":id", infrastructure.AdminOnly(), taskController.UpdateTask)
		taskGroup.POST("", infrastructure.AdminOnly(), taskController.AddTask)
	}

	router.POST("/register", userController.RegisterUser)
	router.POST("/login", userController.LoginUser)

	// Protected route for promoting users
	router.POST("/promote", infrastructure.AuthMiddleware(jwtSecret), infrastructure.AdminOnly(), userController.PromoteUser)

	return router
}