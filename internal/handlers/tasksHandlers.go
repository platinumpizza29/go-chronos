package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/platinumpizza29/go-chronos/internal/models"
	"github.com/platinumpizza29/go-chronos/internal/services"
)

type TasksHandler struct {
	TaskService *services.TaskService
}

func NewTasksHandler(taskService *services.TaskService) *TasksHandler {
	return &TasksHandler{
		TaskService: taskService,
	}
}

func (h *TasksHandler) CreateTask(ctx *gin.Context) {
	var req models.TaskCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error decoding request body: %v", err)
		ctx.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	// add error handling for the task
	if err := h.TaskService.CreateTask(ctx, &req); err != nil {
		log.Printf("error creating task: %v", err)
		ctx.JSON(500, gin.H{"error": "failed to create task"})
		return
	}

	ctx.JSON(201, gin.H{"message": "task created successfully"})
}

func (h *TasksHandler) GetTasks(ctx *gin.Context) {
	// add error handling for the task
	tasks, err := h.TaskService.GetTasks(ctx)
	if err != nil {
		log.Printf("error getting tasks: %v", err)
		ctx.JSON(500, gin.H{"error": "failed to get tasks"})
		return
	}

	ctx.JSON(200, tasks)
}

func (h *TasksHandler) OptimiseTasks(ctx *gin.Context) {
	// add error handling for the task
	priorityResults, err := h.TaskService.OptimiseTasks(ctx)
	if err != nil {
		log.Printf("error optimising tasks: %v", err)
		ctx.JSON(500, gin.H{"error": "failed to optimise tasks"})
		return
	}

	ctx.JSON(200, priorityResults)
}
