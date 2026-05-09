package models

import "time"

type Task struct {
	ID            string    `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	Urgency       int       `json:"urgency" db:"urgency"`
	Severity      int       `json:"severity" db:"severity"`
	PriorityScore int       `json:"priority_score" db:"priority_score"`
	Horizon       int       `json:"horizon" db:"horizon"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type TaskCreateRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Urgency       int    `json:"urgency"`
	Severity      int    `json:"severity"`
	PriorityScore int    `json:"priority_score"`
}

type TaskPriorityResult struct {
	ID            string `json:"id"`
	PriorityScore int    `json:"priority_score"`
	Horizon       int    `json:"horizon"`
}
