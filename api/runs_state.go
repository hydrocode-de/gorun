package api

import (
	"log"
	"net/http"

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

	var runs []db.Run
	var err error
	switch filter {
	case "pending":
		runs, err = c.GetDB().GetIdleRuns(r.Context())
	case "running":
		runs, err = c.GetDB().GetRunning(r.Context())
	case "finished":
		runs, err = c.GetDB().GetFinishedRuns(r.Context())
	case "errored":
		runs, err = c.GetDB().GetErroredRuns(r.Context())
	default:
		runs, err = c.GetDB().GetAllRuns(r.Context())
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
