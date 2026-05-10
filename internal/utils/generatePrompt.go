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

	today := time.Now().UTC()
	dayNames := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	todayWeekday := (int(today.Weekday()) + 6) % 7

	prompt := fmt.Sprintf(`
You are a task scheduling engine.
Today is %s (%s).

Given the following tasks with urgency (1-5) and severity (1-5), return a JSON array.
Each item must have exactly these fields:
- id (string)
- priority_score (int, 0-100)
- horizon (int: 0=Monday, 1=Tuesday, 2=Wednesday, 3=Thursday, 4=Friday, 5=Saturday, 6=Sunday)

Scheduling rules:
1. Calculate priority_score based on urgency and severity (both 1-5) — high urgency + high severity = higher score
2. Sort tasks by priority_score descending
3. Assign the highest-priority task to today (horizon=%d)
4. Assign each subsequent task to the next weekday, overflowing to the following day when a day reaches capacity
5. Estimate capacity per day: roughly N tasks / 5 weekdays, rounded up
6. Only use Saturday(5) or Sunday(6) if all weekdays (0-4) are full
7. Return ONLY a valid JSON array, no explanation, no markdown

Tasks: %s
`, today.Format("2006-01-02"), dayNames[todayWeekday], todayWeekday, string(tasksJSON))

	return prompt, nil
}
