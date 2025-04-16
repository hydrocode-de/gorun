package toolImage

import "github.com/docker/docker/client"

func WithDocker(fn func(c *client.Client) error, strategy string) error {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	return fn(c)
}
