package main

import (
	"context"
	"log"
	"os"
	"time"

	controllers "task_manager/delivery/controllers"
	routers "task_manager/delivery/routers"
	infrastructure "task_manager/infrastructure"
	repositories "task_manager/repositories"
	usecases "task_manager/usecases"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	_ = godotenv.Load()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
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

	// Repositories (pass only client)
	userRepo := repositories.NewUserRepository(client)
	taskRepo := repositories.NewTaskRepository(client)

	// Services
	passwordService := infrastructure.NewPasswordService()
	 jwtSecret := []byte(os.Getenv("JWT_SECRET"))
   if len(jwtSecret) == 0 {
       log.Fatal("JWT_SECRET is not set in environment")
   }
	jwtService := infrastructure.NewJWTService(string(jwtSecret))

	// Usecases
	userUsecase := usecases.NewUserUsecase(userRepo, passwordService, jwtService, 5*time.Second)
	taskUsecase := usecases.NewTaskUsecase(taskRepo, 5*time.Second)

	// Controllers
	userController := controllers.NewUserController(userUsecase)
	taskController := controllers.NewTaskController(taskUsecase)

	// Router
	router := routers.SetupRouter(userController, taskController, jwtSecret)
	router.Run()
}