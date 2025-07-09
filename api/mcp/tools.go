package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/files"
	"github.com/hydrocode-de/gorun/internal/tool"
	"github.com/spf13/viper"
)

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

type McpTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
	Required    []string    `json:"required"`
}

type ToolListResult struct {
	Tools      []McpTool `json:"tools"`
	NextCursor string    `json:"nextCursor,omitempty"`
}

type ToolListResponse struct {
	McpResponse
	Result ToolListResult `json:"result"`
}

type CallParams struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"arguments"`
}
type ToolCallRequest struct {
	McpRequest
	CallParams CallParams `json:"params"`
}

func ListTools(ctx context.Context, id int) (ToolListResponse, error) {
	cache := viper.Get("cache").(*cache.Cache)

	if !cache.Ready {
		select {
		case <-cache.ReadyCh:
		case <-time.After(3 * time.Second):
			if !cache.Ready {
				return ToolListResponse{}, fmt.Errorf("failed to list tools: cache not ready")
			}
		}
	}

	toolSpecs := cache.ListToolSpecs()
	tools := make([]McpTool, len(toolSpecs))
	for _, spec := range toolSpecs {

		properties := make(map[string]Property, len(spec.Parameters))
		required := make([]string, 0)
		for name, param := range spec.Parameters {
			properties[name] = Property{
				Type:        param.ToolType,
				Description: param.Description,
			}
			if param.Optional {
				required = append(required, name)
			}
		}
		for name, dataspec := range spec.Data {
			fullName := fmt.Sprintf("dataset:%s", name)
			description := fmt.Sprintf("A dataset should be a file:// protocol resource uploaded to the server.\nMime: %s\nData description: %s", strings.Join(dataspec.Extensions, ", "), dataspec.Description)
			properties[fullName] = Property{
				Type:        "string",
				Description: description,
			}
		}

		tools = append(tools, McpTool{
			Name:        spec.ID,
			Description: spec.Description,
			InputSchema: InputSchema{
				Type:       "object",
				Properties: properties,
			},
			Required: required,
		})
	}

	return ToolListResponse{
		McpResponse: McpResponse{
			Jsonrpc: JSONRPC_VERSION,
			Id:      id,
		},
		Result: ToolListResult{
			Tools: tools,
		},
	}, nil
}

func RunTool(ctx context.Context, toolName string, params CallParams) ([]files.ResultFile, error) {
	cache := viper.Get("cache").(*cache.Cache)
	spec, ok := cache.GetToolSpec(toolName)
	if !ok {
		return nil, fmt.Errorf("tool %s not found", toolName)
	}

	admin, err := auth.GetAdminCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin credentials: %v", err)
	}
	parts := strings.Split(spec.ID, "::")
	parameters := make(map[string]interface{}, 3)
	datasets := make(map[string]string, 2)
	for key, param := range params.Args {
		if strings.HasPrefix(key, "dataset:") {
			datasets[key] = param.(string)
		} else {
			parameters[key] = param
		}
	}

	run, err := tool.CreateToolRun(ctx, "_mcp", tool.CreateRunOptions{
		Name:       parts[1],
		Image:      parts[0],
		Parameters: parameters,
		Datasets:   datasets,
	}, admin.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create tool run: %v", err)
	}
	toolJob, err := tool.FromDBRun(run)
	if err != nil {
		return nil, fmt.Errorf("failed to convert tool run to tool job: %v", err)
	}

	DB := viper.Get("db").(*db.Queries)
	err = tool.RunTool(ctx, tool.RunToolOptions{
		DB:     DB,
		Tool:   toolJob,
		Env:    []string{},
		UserId: admin.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to run tool: %v", err)
	}

	results, err := toolJob.ListResults()
	if err != nil {
		return nil, fmt.Errorf("failed to list results: %v", err)
	}

	return results, nil
}
