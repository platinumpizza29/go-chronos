package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platinumpizza29/go-chronos/internal/models"
)

type TaskDB struct {
	Pool *pgxpool.Pool
}

func NewTaskDB(pool *pgxpool.Pool) *TaskDB {
	return &TaskDB{
		Pool: pool,
	}
}

// give me error handled CRUD functions for Task
func (t *TaskDB) Create(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (id, title, description, urgency, severity, priority_score, horizon, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := t.Pool.Exec(ctx, query, task.ID, task.Title, task.Description, task.Urgency, task.Severity, task.PriorityScore, task.Horizon, task.Status, task.CreatedAt, task.UpdatedAt)
	return err
}

// get all tasks
func (t *TaskDB) GetAll(ctx context.Context) ([]*models.Task, error) {
	query := `SELECT id, title, description, urgency, severity, priority_score, horizon, status, created_at, updated_at FROM tasks`
	rows, err := t.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Urgency, &task.Severity, &task.PriorityScore, &task.Horizon, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}

// get by id
func (t *TaskDB) GetByID(ctx context.Context, id string) (*models.Task, error) {
	query := `SELECT id, title, description, urgency, severity, priority_score, horizon, status, created_at, updated_at FROM tasks WHERE id = $1`
	var task models.Task
	err := t.Pool.QueryRow(ctx, query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Urgency, &task.Severity, &task.PriorityScore, &task.Horizon, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// update
func (t *TaskDB) Update(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, urgency = $3, severity = $4, priority_score = $5, horizon = $6, status = $7, updated_at = NOW() WHERE id = $8 RETURNING updated_at`
	return t.Pool.QueryRow(ctx, query, task.Title, task.Description, task.Urgency, task.Severity, task.PriorityScore, task.Horizon, task.Status, task.ID).Scan(&task.UpdatedAt)
}

// delete
func (t *TaskDB) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := t.Pool.Exec(ctx, query, id)
	return err
}

func (t *TaskDB) UpdatePriority(ctx context.Context, id string, score int, horizon int) error {
	query := `
		UPDATE tasks 
		SET priority_score = $1, 
		    horizon = $2, 
		    updated_at = $3 
		WHERE id = $4`

	_, err := t.Pool.Exec(ctx, query, score, horizon, time.Now(), id)
	return err
}
