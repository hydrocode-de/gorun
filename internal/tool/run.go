package tool

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/hydrocode-de/gorun/internal/db"
)

type RunToolOptions struct {
	DB   *db.Queries
	Tool Tool
	Env  []string
	Cmd  []string
}

func RunTool(ctx context.Context, c *client.Client, opt RunToolOptions, user_id string) error {
	// create a function to update the database
	updateDB := func(status string, origError error) {
		switch status {
		case "started":
			_, err := opt.DB.StartRun(ctx, db.StartRunParams{
				ID:     opt.Tool.ID,
				UserID: user_id,
			})
			if err != nil {
				log.Fatal(err)
			}
		case "finished":
			_, err := opt.DB.FinishRun(ctx, opt.Tool.ID)
			if err != nil {
				log.Fatal(err)
			}
		case "errored":
			_, err := opt.DB.RunErrored(ctx, db.RunErroredParams{
				ID: opt.Tool.ID,
				ErrorMessage: sql.NullString{
					String: fmt.Sprintf("the execution of the tool (%v) container (%v) errored unexpectedly: %v", opt.Tool.Name, opt.Tool.Image, origError),
					Valid:  true,
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	tool := &opt.Tool
	mounts := make([]mount.Mount, 0, len(tool.Mounts))
	for containerPath, hostPath := range tool.Mounts {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: hostPath,
			Target: containerPath,
		})
	}
	//fmt.Println(mounts)

	config := container.Config{
		Image:        tool.Image,
		Tty:          false,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
	}

	if len(opt.Cmd) != 0 {
		fmt.Printf("Custom CMD: %v\n", opt.Cmd)
		config.Cmd = opt.Cmd
	}
	fmt.Printf("running tool %v with image: %v\n", tool.Name, tool.Image)
	cont, err := c.ContainerCreate(ctx, &config, &container.HostConfig{
		Mounts: mounts,
	}, nil, nil, "")
	if err != nil {
		updateDB("errored", err)
		return err
	}
	defer c.ContainerRemove(ctx, cont.ID, container.RemoveOptions{})
	fmt.Printf("container created: %v\n", cont)

	if err = c.ContainerStart(ctx, cont.ID, container.StartOptions{}); err != nil {
		updateDB("errored", err)
		fmt.Println("starting container failed")
		return err
	}
	updateDB("started", nil)

	statusCh, errCh := c.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		fmt.Println(err)
		updateDB("errored", err)
		return err
	case <-statusCh:
		fmt.Println("container finished")
	}

	logReader, err := c.ContainerLogs(ctx, cont.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return err
	}
	defer logReader.Close()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	stdcopy.StdCopy(stdout, stderr, logReader)

	// create log files in the mounted out volume
	for _, mount := range mounts {
		if mount.Target == "/out" {
			os.WriteFile(path.Join(mount.Source, "STDOUT.log"), stdout.Bytes(), 0644)
			os.WriteFile(path.Join(mount.Source, "STDERR.log"), stderr.Bytes(), 0644)
			break
		}
	}
	updateDB("finished", nil)
	return nil
}
