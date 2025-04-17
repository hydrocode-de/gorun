package main

import (
	"log"
	"os"

	"github.com/hydrocode-de/gorun/cli"
)

func main() {
	cli.Execute()
	os.Exit(0)

	// DEV section
	// read in a argv for a docker image name
	if len(os.Args) == 1 {
		log.Println("No command line arguments provided")
	}

	// if os.Args[1] == "find" {
	// 	if len(os.Args) < 3 {
	// 		log.Fatal("missing pattern - use like gorun find \"<pattern>\"")
	// 	}
	// 	matches, err := files.Find(os.Args[2], config.GetMountPath(), "all")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	for _, match := range matches {
	// 		fmt.Printf("- %v\n", match)
	// 	}
	// }

	// if os.Args[1] == "auth" {
	// 	if len(os.Args) < 3 {
	// 		log.Fatal("missing command - use like gorun auth create")
	// 	}

	// 	if os.Args[2] == "credentials" {
	// 		credentials, err := config.GetAdminCredentials()
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		if len(os.Args) > 2 && os.Args[2] == "--token" {
	// 			fmt.Printf("%s\n", credentials.RefreshToken)
	// 			return
	// 		}
	// 		fmt.Println("Admin credentials:")
	// 		fmt.Printf("Email:      %s\n", credentials.Email)
	// 		fmt.Printf("User ID:       %s\n", credentials.UserID)
	// 		fmt.Printf("Refresh token: %s\n", credentials.RefreshToken)
	// 		fmt.Printf("Expires at:    %s\n", credentials.ExpiresAt)
	// 		fmt.Printf("Access token  \n\n%s\n", credentials.AccessToken)
	// 	}

	// 	if os.Args[2] == "create-user" {
	// 		if len(os.Args) < 6 {
	// 			log.Fatal("missing arguments - use like gorun auth create-user <email> <password> <is-admin>")
	// 		}

	// 		_, err := auth.CreateUser(context.Background(), config.GetDB(), os.Args[3], os.Args[4], os.Args[5] == "true", config.Secret)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		log.Println("User created")
	// 	}
	// }
}
