package toolImage

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/toolSpec"
)

func ReadAllTools(ctx context.Context, c *client.Client, cache *cache.Cache) ([]string, error) {
	summary, err := c.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}
	var tools []string
	for _, img := range summary {
		if len(img.RepoTags) == 0 {
			continue
		}
		imgTag := img.RepoTags[0]
		//log.Printf("checking %s...\n", imgTag)
		image, ok := cache.GetImageSpec(imgTag)
		if !ok {
			spec, err := ReadToolSpec(ctx, c, imgTag)
			if err != nil {
				cache.SetImageSpec(imgTag, toolSpec.SpecFile{})
			} else {
				cache.SetImageSpec(imgTag, spec)
				for name, tool := range spec.Tools {
					cache.SetToolSpec(name, &tool)
					tools = append(tools, fmt.Sprintf("%s::%s", imgTag, name))
				}
			}
		} else {
			for name := range image.Tools {
				tools = append(tools, fmt.Sprintf("%s::%s", imgTag, name))
			}
		}
	}

	return tools, nil
}

func LoadToolSpec(ctx context.Context, c *client.Client, toolSlug string, cache *cache.Cache) (toolSpec.ToolSpec, error) {
	chunks := strings.Split(toolSlug, "::")
	if len(chunks) == 1 {
		spec, ok := cache.GetToolSpec(toolSlug)
		if !ok {
			return toolSpec.ToolSpec{}, fmt.Errorf("the tool %s was not found in the cache. Try to call like <image-name>::<tool-name>", toolSlug)
		}
		return *spec, nil
	}

	if len(chunks) == 2 {
		imageName := chunks[0]
		toolName := chunks[1]
		spec, ok := cache.GetImageSpec(imageName)
		if !ok {
			specFile, err := ReadToolSpec(ctx, c, imageName)
			if err != nil {
				return toolSpec.ToolSpec{}, err
			}
			cache.SetImageSpec(imageName, specFile)
			for name, tool := range specFile.Tools {
				cache.SetToolSpec(name, &tool)
			}
			tool, ok := specFile.Tools[toolName]
			if !ok {
				return toolSpec.ToolSpec{}, fmt.Errorf("the tool %s was not found in the image %s", toolName, imageName)
			}
			return tool, nil
		} else {
			tool, ok := spec.Tools[toolName]
			if !ok {
				return toolSpec.ToolSpec{}, fmt.Errorf("the tool %s was not found in the image %s", toolName, imageName)
			}
			return tool, nil
		}
	}
	return toolSpec.ToolSpec{}, fmt.Errorf("invalid tool slug: %s", toolSlug)
}

func ReadToolSpec(ctx context.Context, c *client.Client, imageName string) (toolSpec.SpecFile, error) {
	cont, err := c.ContainerCreate(ctx, &container.Config{
		Image:      imageName,
		Entrypoint: []string{"cat"},
		Cmd:        []string{"/src/tool.yml"},
	}, &container.HostConfig{}, nil, nil, "")
	if err != nil {
		return toolSpec.SpecFile{}, err
	}
	defer c.ContainerRemove(ctx, cont.ID, container.RemoveOptions{})

	if err = c.ContainerStart(ctx, cont.ID, container.StartOptions{}); err != nil {
		return toolSpec.SpecFile{}, err
	}

	statusCh, errCh := c.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return toolSpec.SpecFile{}, err
	case <-statusCh:
	}

	logReader, err := c.ContainerLogs(ctx, cont.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return toolSpec.SpecFile{}, err
	}
	defer logReader.Close()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	stdcopy.StdCopy(stdout, stderr, logReader)

	// check if the sdterr is empty
	if stderr.Len() != 0 {
		return toolSpec.SpecFile{}, fmt.Errorf("the container errored while identifying the tool spec: %v", stderr.String())
	}
	if stdout.Len() == 0 {
		return toolSpec.SpecFile{}, fmt.Errorf("the container did not respond")
	}
	out := stdout.Bytes()

	spec, err := toolSpec.LoadToolSpec(out)
	if err != nil {
		return toolSpec.SpecFile{}, fmt.Errorf("the container %s did not contain a valid tool-spec at /src/tool.yml: %v", imageName, err)
	}

	return spec, nil
}
