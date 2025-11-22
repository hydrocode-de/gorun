package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/files"
	"github.com/hydrocode-de/gorun/internal/helper"
	"github.com/hydrocode-de/gorun/internal/toolImage"
	"github.com/hydrocode-de/gorun/pkg/toolspec"
	"github.com/spf13/viper"
)

type CreateRunOptions struct {
	Name       string
	Image      string
	Parameters map[string]interface{}
	Datasets   map[string]string
}

func CreateToolRun(ctx context.Context, mountStrategy string, opts CreateRunOptions, user_id string) (db.Run, error) {
	DB := viper.Get("db").(*db.Queries)
	mountPath := viper.GetString("mount_path")

	spec, err := toolImage.ReadToolSpec(ctx, opts.Image)
	if err != nil {
		return db.Run{}, err
	}
	toolSpec, err := spec.GetTool(opts.Name)
	if err != nil {
		return db.Run{}, err
	}

	mounts := files.CreateNewMountPaths(mountPath, mountStrategy)
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
	inputJSON, err := json.MarshalIndent(toolspec.InputFile{
		opts.Name: toolspec.ToolInput{
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
	runData, err := DB.CreateRun(ctx, db.CreateRunParams{
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
