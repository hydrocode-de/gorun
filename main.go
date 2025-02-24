package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/hydrocode-de/gorun/api"
	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/files"
	"github.com/hydrocode-de/gorun/internal/toolImage"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer docker.Close()

	config := config.APIConfig{}
	err = config.Load(docker)
	if err != nil {
		log.Fatal(err)
	}
	err = config.Validate()
	if err != nil {
		log.Fatal(err)
	}

	// create a ticker to cleanup uploads
	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		for range ticker.C {
			fmt.Println("Running cleanup")
			err := files.Cleanup(&config)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// DEV section
	// read in a argv for a docker image name
	if len(os.Args) > 1 {
		if os.Args[1] == "serve" {
			mux, err := api.CreateServer(&config)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("GoRun server listening on port %d\n", config.Port)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), mux))
		}

		if os.Args[1] == "list" {
			tools, err := toolImage.ReadAllTools(context.Background(), docker, &config.Cache)
			if err != nil {
				log.Fatal(err)
			}
			for _, tool := range tools {
				fmt.Println(tool)
			}
		}

	} else {
		fmt.Println("No command line arguments provided")
	}
}
