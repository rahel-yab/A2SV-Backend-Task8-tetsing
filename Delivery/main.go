package main

import (
	"context"
	"log"
	"os"
	Controllers "task_manager/Delivery/Controllers"
	routers "task_manager/Delivery/routers"
	Infrastructure "task_manager/Infrastructure"
	Repositories "task_manager/Repositories"
	Usecases "task_manager/Usecases"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables from .env if present
	_ = godotenv.Load()

	// MongoDB connection
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // default fallback
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	db := client.Database("task_manager")

	// Repositories
	userRepo := Repositories.NewUserRepository(db.Collection("users"))
	taskRepo := Repositories.NewTaskRepository(db.Collection("tasks"))

	// Services
	passwordService := Infrastructure.NewPasswordService()
	jwtSecret := []byte("your_jwt_secret") // TODO: move to config/env
	jwtService := Infrastructure.NewJWTService(string(jwtSecret))

	// Usecases
	userUsecase := Usecases.NewUserUsecase(userRepo, passwordService, jwtService, 5*time.Second)
	taskUsecase := Usecases.NewTaskUsecase(taskRepo, 5*time.Second)

	// Controllers
	userController := Controllers.NewUserController(userUsecase)
	taskController := Controllers.NewTaskController(taskUsecase)

	// Router
	router := routers.SetupRouter(userController, taskController, jwtSecret)
	router.Run()
}