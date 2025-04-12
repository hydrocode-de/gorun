package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/helper"
	"github.com/hydrocode-de/gorun/internal/toolImage"
)

type CreateRunOptions struct {
	Name       string
	Image      string
	Parameters map[string]interface{}
	Datasets   map[string]string
}

type ToolInput struct {
	Parameters map[string]interface{} `json:"parameters"`
	Datasets   map[string]string      `json:"data"`
}
type inputFile map[string]ToolInput

func CreateToolRun(ctx context.Context, mountStrategy string, opts CreateRunOptions, user_id string, c *config.APIConfig) (db.Run, error) {
	spec, err := toolImage.ReadToolSpec(ctx, c.GetDockerClient(), opts.Image)
	if err != nil {
		return db.Run{}, err
	}
	toolSpec, err := spec.GetTool(opts.Name)
	if err != nil {
		return db.Run{}, err
	}

	mounts := c.CreateNewMountPaths(mountStrategy)
	datasets := make(map[string]string)

	for dataName, dataPath := range opts.Datasets {
		containerPath := path.Join(mounts["/in"], filepath.Base(dataPath))
		err := helper.CopyFile(dataPath, containerPath)
		if err != nil {
			return db.Run{}, err
		}
		datasets[dataName] = path.Join("/in", filepath.Base(dataPath))
	}

	// create the input file
	inputJSON, err := json.MarshalIndent(inputFile{
		opts.Name: {
			Parameters: opts.Parameters,
			Datasets:   datasets,
		},
	}, "", "\t")
	if err != nil {
		return db.Run{}, err
	}
	err = os.WriteFile(path.Join(mounts["/in"], "inputs.json"), inputJSON, 0644)
	if err != nil {
		return db.Run{}, err
	}

	// marshal the other stuff
	parJSON, parErr := json.Marshal(opts.Parameters)
	dataJSON, dataErr := json.Marshal(datasets)
	mountJSON, mountErr := json.Marshal(mounts)
	if dataErr != nil || mountErr != nil || parErr != nil {
		return db.Run{}, fmt.Errorf("failed to marshal parameters and mount points")
	}

	// create the database entry
	runData, err := c.GetDB().CreateRun(ctx, db.CreateRunParams{
		Name:        opts.Name,
		Title:       toolSpec.Title,
		Description: toolSpec.Description,
		DockerImage: opts.Image,
		Parameters:  string(parJSON),
		Data:        string(dataJSON),
		Mounts:      string(mountJSON),
		UserID:      user_id,
	})
	if err != nil {
		return db.Run{}, err
	}

	return runData, nil
}
