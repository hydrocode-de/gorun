package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/tool"
)

type RunsResponse struct {
	Count  int         `json:"count"`
	Status string      `json:"status"`
	Runs   []tool.Tool `json:"runs"`
}

func GetAllRuns(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {
	filter := r.URL.Query().Get("status")

	user_id := r.Header.Get("X-User-ID")
	if user_id == "" {
		RespondWithError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	var runs []db.Run
	var err error
	switch filter {
	case "pending":
		runs, err = c.GetDB().GetIdleRuns(r.Context(), db.GetIdleRunsParams{
			UserID: user_id,
		})
	case "running":
		runs, err = c.GetDB().GetRunning(r.Context(), db.GetRunningParams{
			UserID: user_id,
		})
	case "finished":
		runs, err = c.GetDB().GetFinishedRuns(r.Context(), db.GetFinishedRunsParams{
			UserID: user_id,
		})
	case "errored":
		runs, err = c.GetDB().GetErroredRuns(r.Context(), db.GetErroredRunsParams{
			UserID: user_id,
		})
	default:
		runs, err = c.GetDB().GetAllRuns(r.Context(), db.GetAllRunsParams{
			UserID: user_id,
		})
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var toolRuns []tool.Tool
	for _, run := range runs {
		toolRun, err := tool.FromDBRun(run)
		if err != nil {
			log.Printf("Error while loading tool run: %s", err)
			continue
		}
		toolRuns = append(toolRuns, toolRun)
	}

	ResondWithJSON(w, http.StatusOK, RunsResponse{
		Count:  len(runs),
		Status: filter,
		Runs:   toolRuns,
	})
}

func DeleteRun(w http.ResponseWriter, r *http.Request, tool tool.Tool, c *config.APIConfig) {
	user_id := r.Header.Get("X-User-ID")
	if user_id == "" {
		RespondWithError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	// the tool may have a saved mount point, so we delete it first
	_, ok := tool.Mounts["/in"]
	if ok {
		parent := filepath.Dir(tool.Mounts["/in"])
		err := os.RemoveAll(parent)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}

	}

	err := c.GetDB().DeleteRun(r.Context(), db.DeleteRunParams{
		ID:     tool.ID,
		UserID: user_id,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	ResondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Run deleted",
	})

}

func GetRunStatus(w http.ResponseWriter, r *http.Request, run tool.Tool, c *config.APIConfig) {
	ResondWithJSON(w, http.StatusOK, run)
}

func HandleRunStart(w http.ResponseWriter, r *http.Request, run tool.Tool, c *config.APIConfig) {
	user_id := r.Header.Get("X-User-ID")
	if user_id == "" {
		RespondWithError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	opt := tool.RunToolOptions{
		DB:   (*c).GetDB(),
		Tool: run,
		Env:  []string{},
		// Cmd:  []string{},
	}

	go tool.RunTool(context.Background(), (*c).GetDockerClient(), opt, user_id)

	// wait a few miliseconds to make sure the container is started
	time.Sleep(time.Millisecond * 100)
	started, err := (*c).GetDB().GetRun(r.Context(), db.GetRunParams{
		ID:     run.ID,
		UserID: user_id,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	ResondWithJSON(w, http.StatusProcessing, started)
}
