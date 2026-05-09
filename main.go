package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/platinumpizza29/go-chronos/internal/db"
	"github.com/platinumpizza29/go-chronos/internal/handlers"
	"github.com/platinumpizza29/go-chronos/internal/services"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	dbUrl := os.Getenv("DATABASE_URL")

	if err := db.Connect(context.Background(), dbUrl); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	taskDB := db.NewTaskDB(db.Pool)
	taskService := services.NewTaskService(taskDB)
	taskHandler := handlers.NewTasksHandler(taskService)

	router := gin.Default()

	// For development, allow all origins.
	// For production, you should restrict this to your frontend's domain.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	//define routes here
	tasksRoutes := router.Group("/api/tasks")
	{
		tasksRoutes.POST("/", taskHandler.CreateTask)
		tasksRoutes.GET("/", taskHandler.GetTasks)
		tasksRoutes.POST("/optimise", taskHandler.OptimiseTasks)
	}

	// Define routes
	router.Run()
}
