package main

import (
	"os"
	Controllers "task_manager/Delivery/Controllers"
	routers "task_manager/Delivery/routers"
	Infrastructure "task_manager/Infrastructure"
	Repositories "task_manager/Repositories"
	Usecases "task_manager/Usecases"
)

func main() {
	// MongoDB connection
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // default fallback
	}
	client := Infrastructure.ConnectMongo(uri)
	db := client.Database("task_manager")

	// Repositories
	userRepo := &Repositories.MongoUserRepository{Collection: db.Collection("users")}
	taskRepo := &Repositories.MongoTaskRepository{Collection: db.Collection("tasks")}

	// Services
	passwordService := Infrastructure.NewPasswordService()
	jwtSecret := []byte("your_jwt_secret") // TODO: move to config/env
	jwtService := Infrastructure.NewJWTService(string(jwtSecret))

	// Usecases
	userUsecase := &Usecases.UserUsecase{
		Repo:            userRepo,
		PasswordService: passwordService,
		JWTService:      jwtService,
	}
	taskUsecase := &Usecases.TaskUsecase{
		Repo: taskRepo,
	}

	// Controllers
	userController := &Controllers.UserController{UserUsecase: userUsecase}
	taskController := &Controllers.TaskController{TaskUsecase: taskUsecase}

	// Router
	router := routers.SetupRouter(userController, taskController, jwtSecret)
	router.Run()
}

