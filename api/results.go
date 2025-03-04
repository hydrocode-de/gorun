package api

import (
	"fmt"
	"net/http"

	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/files"
	"github.com/hydrocode-de/gorun/internal/tool"
)

type ListRunResultsResponse struct {
	Count int                `json:"count"`
	Files []files.ResultFile `json:"files"`
}

func ListRunResults(w http.ResponseWriter, r *http.Request, tool tool.Tool, c *config.APIConfig) {
	results, err := tool.ListResults()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ResondWithJSON(w, http.StatusOK, ListRunResultsResponse{
		Count: len(results),
		Files: results,
	})
}

func GetResultFile(w http.ResponseWriter, r *http.Request, tool tool.Tool, c *config.APIConfig) {
	//filename := r.URL.Query().Get("filename")
	filename := r.PathValue("filename")

	info, err := tool.WriteResultFile(filename, w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", info.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", info.Filename))
}
