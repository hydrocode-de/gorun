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
	"github.com/hydrocode-de/gorun/internal/auth"
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
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for range ticker.C {
			log.Println("Running cleanup")
			err := files.Cleanup(&config)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	go toolImage.ReadAllTools(context.Background(), docker, &config.Cache)
	go func() {
		for range ticker.C {
			log.Println("Checking for new images...")
			_, err := toolImage.ReadAllTools(context.Background(), docker, &config.Cache)
			if err != nil {
				log.Fatalf("Errored while updating tool images: %s", err)
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
			server := api.EnableCORS(mux, "*")
			log.Printf("GoRun server listening on port %d\n", config.Port)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), server))
		}

		if os.Args[1] == "list" {
			tools, err := toolImage.ReadAllTools(context.Background(), docker, &config.Cache)
			if err != nil {
				log.Fatal(err)
			}
			for _, tool := range tools {
				log.Println(tool)
			}
		}

		if os.Args[1] == "cleanup" {
			err = files.Cleanup(&config)
			if err != nil {
				log.Fatal(err)
			}
		}

		if os.Args[1] == "find" {
			if len(os.Args) < 3 {
				log.Fatal("missing pattern - use like gorun find \"<pattern>\"")
			}
			matches, err := files.Find(os.Args[2], config.GetMountPath(), "all")
			if err != nil {
				log.Fatal(err)
			}

			for _, match := range matches {
				fmt.Printf("- %v\n", match)
			}
		}

		if os.Args[1] == "create-key" {
			key, err := auth.CreateNewApiKey(context.Background(), config.GetDB(), 0)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Created a new API key for GoRun. Store this key, as you can't retrieve it later:")
			fmt.Println(key)
			fmt.Println()
		}

	} else {
		log.Println("No command line arguments provided")
	}
}
