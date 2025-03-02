package tool

import (
	"encoding/json"
	"time"

	"github.com/hydrocode-de/gorun/internal/db"
)

type Tool struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Image       string                 `json:"image"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Data        map[string]string      `json:"data,omitempty"`
	Mounts      map[string]string      `json:"mounts,omitempty"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   time.Time              `json:"started_at,omitempty"`
	FinishedAt  time.Time              `json:"finished_at,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

func FromDBRun(run db.Run) (Tool, error) {
	tool := Tool{
		ID:          run.ID,
		Name:        run.Name,
		Title:       run.Title,
		Image:       run.DockerImage,
		Description: run.Description,
		Parameters:  nil,
		Data:        nil,
		Mounts:      nil,
		Status:      run.Status,
		CreatedAt:   run.CreatedAt,
		StartedAt:   run.StartedAt.Time,
		FinishedAt:  run.FinishedAt.Time,
		Error:       run.ErrorMessage.String,
	}
	err := json.Unmarshal([]byte(run.Parameters), &tool.Parameters)
	if err != nil {
		return Tool{}, err
	}
	err = json.Unmarshal([]byte(run.Data), &tool.Data)
	if err != nil {
		return Tool{}, err
	}
	err = json.Unmarshal([]byte(run.Mounts), &tool.Mounts)
	if err != nil {
		return Tool{}, err
	}

	return tool, nil
}
