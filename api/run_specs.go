package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/tool"
	"github.com/hydrocode-de/gorun/internal/toolSpec"
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

func RunMiddleware(handler func(http.ResponseWriter, *http.Request, tool.Tool, *config.APIConfig), c *config.APIConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idPath := r.PathValue("id")
		if idPath == "" {
			RespondWithError(w, http.StatusBadRequest, "missing run id")
			return
		}
		id, err := strconv.ParseInt(idPath, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("the passed run id is not a valid integer: %v", err))
		}

		run, err := c.GetDB().GetRun(r.Context(), id)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		tool, err := tool.FromDBRun(run)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}

		handler(w, r, tool, c)
	}
}

func GetToolSpec(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {
	toolName := r.PathValue("toolname")
	if toolName == "" {
		RespondWithError(w, http.StatusNotFound, "missing tool name")
	}

	spec, wasFound := c.Cache.GetToolSpec(toolName)
	if !wasFound {
		RespondWithError(w, http.StatusNotFound, "tool not found")
	}
	ResondWithJSON(w, http.StatusOK, spec)
}

func ListToolSpecs(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {
	specs := c.Cache.ListToolSpecs()

	ResondWithJSON(w, http.StatusOK, ListToolSpecResponse{
		Count: len(specs),
		Tools: specs,
	})
}

func CreateRun(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {
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
	runData, err := tool.CreateToolRun("_random", opts, c, r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ResondWithJSON(w, http.StatusCreated, runData)
}

func GetRunStatus(w http.ResponseWriter, r *http.Request, run tool.Tool, c *config.APIConfig) {
	ResondWithJSON(w, http.StatusOK, run)
}

func HandleRunStart(w http.ResponseWriter, r *http.Request, run tool.Tool, c *config.APIConfig) {
	opt := tool.RunToolOptions{
		DB:   (*c).GetDB(),
		Tool: run,
		Env:  []string{},
		// Cmd:  []string{},
	}

	go tool.RunTool(context.Background(), (*c).GetDockerClient(), opt)

	// wait a few miliseconds to make sure the container is started
	time.Sleep(time.Millisecond * 50)
	started, err := (*c).GetDB().GetRun(r.Context(), run.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	ResondWithJSON(w, http.StatusProcessing, started)
}
