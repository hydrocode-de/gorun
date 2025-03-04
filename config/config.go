package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/helper"
	"github.com/hydrocode-de/gorun/sql"
)

type APIConfig struct {
	Port          int
	MaxUploadSize int64
	dbPath        string
	db            *db.Queries
	docker        *client.Client
	baseMountPath string
	BaseTempDir   string
	MaxTempAge    time.Duration
	Cache         cache.Cache
}

func (c *APIConfig) Validate() error {
	if c.Port == 0 {
		return fmt.Errorf("port is required")
	}

	// check if it is a valid sqlite3 connection string
	if _, err := os.Stat(c.dbPath); os.IsNotExist(err) {
		return fmt.Errorf("the database file %s does not exist", c.dbPath)
	}

	return nil
}

func (c *APIConfig) Load(docker *client.Client) error {
	// Get port from environment
	portStr, isSet := os.LookupEnv("GORUN_PORT")
	if !isSet {
		c.Port = 8080
	} else {
		// Convert string to int
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid port number: %v", err)
		}
		c.Port = port
	}

	// Get DB connection string - place a new one in the home directory if not set
	dbURL := os.Getenv("GORUN_DB")
	if dbURL == "" {
		p := path.Join(os.Getenv("HOME"), "gorun")
		os.MkdirAll(p, 0755)
		c.dbPath = path.Join(p, "gorun.db")
	} else {
		c.dbPath = dbURL
	}

	c.docker = docker

	// Initialize the database driver
	drv, err := sql.CreateDB(c.dbPath)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := db.New(drv)
	c.db = dbQueries

	// have a cache
	c.Cache = cache.Cache{}
	c.Cache.Reset()

	// setup the base mount path
	c.baseMountPath = os.Getenv("GORUN_MOUNT_PATH")
	if c.baseMountPath == "" {
		c.baseMountPath = path.Join(os.Getenv("HOME"), "gorun", "mounts")
	}

	// some defaults
	if c.MaxUploadSize == 0 {
		c.MaxUploadSize = 1024 * 1024 * 1024 * 2 // 2GB
	}

	if c.BaseTempDir == "" {
		c.BaseTempDir = path.Join(os.TempDir(), "gorun")
	}
	if c.MaxTempAge == 0 {
		c.MaxTempAge = 12 * time.Hour
	}

	return nil
}

func (c *APIConfig) GetDB() *db.Queries {
	return c.db
}

func (c *APIConfig) GetMountConfig() map[string]string {
	return make(map[string]string)
}

func (c *APIConfig) GetDockerClient() *client.Client {
	return c.docker
}

func (c *APIConfig) CreateNewMountPaths(strategy string) map[string]string {
	mounts := make(map[string]string)

	if strategy == "_random" {
		level := helper.GetRandomString(12)
		mounts["/in"] = path.Join(c.baseMountPath, level, "in")
		mounts["/out"] = path.Join(c.baseMountPath, level, "out")
	} else {
		mounts["/in"] = path.Join(c.baseMountPath, strategy, "in")
		mounts["/out"] = path.Join(c.baseMountPath, strategy, "out")
	}

	for _, hostPath := range mounts {
		os.MkdirAll(hostPath, 0755)
	}

	return mounts
}

func (c *APIConfig) GetMountPath() string {
	return c.baseMountPath
}
