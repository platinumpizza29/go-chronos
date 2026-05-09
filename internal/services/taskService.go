package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/platinumpizza29/go-chronos/internal/db"
	"github.com/platinumpizza29/go-chronos/internal/models"
	"github.com/platinumpizza29/go-chronos/internal/utils"
	"google.golang.org/genai"
)

type TaskService struct {
	TaskDB *db.TaskDB
}

func NewTaskService(taskDB *db.TaskDB) *TaskService {
	return &TaskService{
		TaskDB: taskDB,
	}
}

func (t *TaskService) CreateTask(ctx context.Context, task *models.TaskCreateRequest) error {
	// add error handling for the task

	// create a new task
	newTask := &models.Task{
		ID:            uuid.New().String(),
		Title:         task.Title,
		Description:   task.Description,
		Urgency:       task.Urgency,
		Severity:      task.Severity,
		PriorityScore: task.PriorityScore,
		Horizon:       0,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err := t.TaskDB.Create(ctx, newTask)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) GetTasks(ctx context.Context) ([]*models.Task, error) {
	// add error handling for the task
	tasks, err := t.TaskDB.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// optimise the task service
// call google cloud bigquery to get the tasks
func (t *TaskService) OptimiseTasks(ctx context.Context) ([]*models.TaskPriorityResult, error) {
	gemini_api_key := os.Getenv("GEMINI_API_KEY")

	// 1. Fetch tasks from DB
	tasks, err := t.TaskDB.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// 2. Initialize the client
	// Note: It's better to pass the API key via environment variable (GOOGLE_API_KEY)
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  gemini_api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	// 3. Generate the prompt
	prompt, err := utils.GeneratePrompt(tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to generate prompt: %w", err)
	}

	// 4. Define the schema and config
	taskSchema := &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"id":             {Type: genai.TypeString},
				"priority_score": {Type: genai.TypeInteger},
				"horizon":        {Type: genai.TypeInteger},
			},
			Required: []string{"id", "priority_score", "horizon"},
		},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   taskSchema,
	}

	// 5. Call the model
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview", // Use 2.0 Flash for better JSON performance
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("ai generation failed: %w", err)
	}

	// 6. Parse the JSON result
	var priorityResults []*models.TaskPriorityResult

	// Use result.Text() helper to get the generated string
	responseText := result.Text()
	if err := json.Unmarshal([]byte(responseText), &priorityResults); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return priorityResults, nil
}
