package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/tool"
	toolspec "github.com/hydrocode-de/tool-spec-go"
	"github.com/hydrocode-de/tool-spec-go/validate"
	"github.com/spf13/viper"
)

type ListToolSpecResponse struct {
	Count int                 `json:"count"`
	Tools []toolspec.ToolSpec `json:"tools"`
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
	RespondWithJSON(w, http.StatusOK, spec)
}

func ListToolSpecs(w http.ResponseWriter, r *http.Request) {
	Cache := viper.Get("cache").(*cache.Cache)
	specs := Cache.ListToolSpecs()

	RespondWithJSON(w, http.StatusOK, ListToolSpecResponse{
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

	Cache := viper.Get("cache").(*cache.Cache)
	toolSlug := fmt.Sprintf("%s::%s", payload.DockerImage, payload.ToolName)
	toolSpec, wasFound := Cache.GetToolSpec(toolSlug)
	if !wasFound {
		RespondWithError(w, http.StatusNotFound, fmt.Sprintf("a tool %s was not found in the cache", toolSlug))
		return
	}
	hasErrors, errs := validate.ValidateInputs(*toolSpec, toolspec.ToolInput{
		Parameters: payload.Parameters,
		Datasets:   payload.DataPaths,
	})
	if hasErrors {
		RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
			"message": fmt.Sprintf("the provided payload is invalid for the tool %s", toolSlug),
			"errors":  errs,
		})
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

	RespondWithJSON(w, http.StatusCreated, runData)
}
