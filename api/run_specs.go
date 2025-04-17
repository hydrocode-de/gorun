package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/tool"
	"github.com/hydrocode-de/gorun/internal/toolSpec"
	"github.com/spf13/viper"
)

type ListToolSpecResponse struct {
	Count int                 `json:"count"`
	Tools []toolSpec.ToolSpec `json:"tools"`
}

type CreateRunPayload struct {
	ToolName    string                 `json:"name"`
	DockerImage string                 `json:"docker_image"`
	Parameters  map[string]interface{} `json:"parameters"`
	DataPaths   map[string]string      `json:"data"`
}

func RunMiddleware(handler func(http.ResponseWriter, *http.Request, tool.Tool)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id := r.Header.Get("X-User-ID")
		if user_id == "" {
			RespondWithError(w, http.StatusUnauthorized, "User ID is required")
			return
		}
		DB := viper.Get("db").(*db.Queries)

		idPath := r.PathValue("id")
		if idPath == "" {
			RespondWithError(w, http.StatusBadRequest, "missing run id")
			return
		}
		id, err := strconv.ParseInt(idPath, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("the passed run id is not a valid integer: %v", err))
		}

		run, err := DB.GetRun(r.Context(), db.GetRunParams{
			ID:     id,
			UserID: user_id,
		})
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		tool, err := tool.FromDBRun(run)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}

		handler(w, r, tool)
	}
}

func GetToolSpec(w http.ResponseWriter, r *http.Request) {
	toolName := r.PathValue("toolname")
	if toolName == "" {
		RespondWithError(w, http.StatusNotFound, "missing tool name")
	}

	Cache := viper.Get("cache").(*cache.Cache)
	spec, wasFound := Cache.GetToolSpec(toolName)
	if !wasFound {
		RespondWithError(w, http.StatusNotFound, "tool not found")
	}
	ResondWithJSON(w, http.StatusOK, spec)
}

func ListToolSpecs(w http.ResponseWriter, r *http.Request) {
	Cache := viper.Get("cache").(*cache.Cache)
	specs := Cache.ListToolSpecs()

	ResondWithJSON(w, http.StatusOK, ListToolSpecResponse{
		Count: len(specs),
		Tools: specs,
	})
}

func CreateRun(w http.ResponseWriter, r *http.Request) {
	user_id := r.Header.Get("X-User-ID")
	if user_id == "" {
		RespondWithError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	var payload CreateRunPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// create the mount paths with random strategy
	opts := tool.CreateRunOptions{
		Name:       payload.ToolName,
		Image:      payload.DockerImage,
		Parameters: payload.Parameters,
		Datasets:   payload.DataPaths,
	}
	runData, err := tool.CreateToolRun(r.Context(), "_random", opts, user_id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ResondWithJSON(w, http.StatusCreated, runData)
}
