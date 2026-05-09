package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

func GeneratePrompt(tasks interface{}) (string, error) {
	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`
You are a task prioritisation engine.
Today is %s.

Given the following tasks with urgency (1-5) and severity (1-5), return a JSON array.
Each item must have exactly these fields:
- id (string)
- priority_score (int, 0-100)
- horizon (int: 0 = today, 1 = this week, 2 = this month)

Rules:
- High urgency + high severity = high score, horizon 0
- Low urgency + high severity = medium score, horizon 1
- Low urgency + low severity = low score, horizon 2
- Return ONLY a valid JSON array, no explanation, no markdown

Tasks:
%s
`, time.Now().UTC().Format("2006-01-02"), string(tasksJSON))

	return prompt, nil
}
