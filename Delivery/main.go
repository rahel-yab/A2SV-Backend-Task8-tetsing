package main

import (
	"context"
	"log"
	"os"
	"time"

	Controllers "task_manager/Delivery/controllers"
	routers "task_manager/Delivery/routers"
	Infrastructure "task_manager/Infrastructure"
	Repositories "task_manager/Repositories"
	Usecases "task_manager/Usecases"

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
	userRepo := Repositories.NewUserRepository(client)
	taskRepo := Repositories.NewTaskRepository(client)

	// Services
	passwordService := Infrastructure.NewPasswordService()
	 jwtSecret := []byte(os.Getenv("JWT_SECRET"))
   if len(jwtSecret) == 0 {
       log.Fatal("JWT_SECRET is not set in environment")
   }
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