package api

import (
	"net/http"

	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/db"
)

type RunsResponse struct {
	Count  int      `json:"count"`
	Status string   `json:"status"`
	Runs   []db.Run `json:"runs"`
}

func GetAllRuns(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {
	filter := r.URL.Query().Get("status")

	var runs []db.Run
	var err error
	switch filter {
	case "idle":
		runs, err = c.GetDB().GetIdleRuns(r.Context())
	case "running":
		runs, err = c.GetDB().GetRunning(r.Context())
	case "finished":
		runs, err = c.GetDB().GetFinishedRuns(r.Context())
	case "erored":
		runs, err = c.GetDB().GetErroredRuns(r.Context())
	default:
		runs, err = c.GetDB().GetAllRuns(r.Context())
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ResondWithJSON(w, http.StatusOK, RunsResponse{
		Count:  len(runs),
		Status: filter,
		Runs:   runs,
	})
}
