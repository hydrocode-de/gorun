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

	// create a ticker to renew the admin credentials every 50 minutes
	tickerAdmin := time.NewTicker(time.Minute * 50)
	go func() {
		for range tickerAdmin.C {
			log.Println("Renewing admin credentials")
			if _, err := config.GetAdminCredentials(); err != nil {
				log.Printf("Failed to renew admin credentials: %s...\n", err)
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
	if len(os.Args) == 1 {
		log.Println("No command line arguments provided")
	}

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

	if os.Args[1] == "auth" {
		if len(os.Args) < 3 {
			log.Fatal("missing command - use like gorun auth create")
		}

		if os.Args[2] == "credentials" {
			credentials, err := config.GetAdminCredentials()
			if err != nil {
				log.Fatal(err)
			}

			if len(os.Args) > 2 && os.Args[2] == "--token" {
				fmt.Printf("%s\n", credentials.RefreshToken)
				return
			}
			fmt.Println("Admin credentials:")
			fmt.Printf("Email:      %s\n", credentials.Email)
			fmt.Printf("User ID:       %s\n", credentials.UserID)
			fmt.Printf("Refresh token: %s\n", credentials.RefreshToken)
			fmt.Printf("Expires at:    %s\n", credentials.ExpiresAt)
			fmt.Printf("Access token  \n\n%s\n", credentials.AccessToken)
		}

		if os.Args[2] == "create-user" {
			if len(os.Args) < 6 {
				log.Fatal("missing arguments - use like gorun auth create-user <email> <password> <is-admin>")
			}

			_, err := auth.CreateUser(context.Background(), config.GetDB(), os.Args[3], os.Args[4], os.Args[5] == "true", config.Secret)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("User created")
		}
	}
}
