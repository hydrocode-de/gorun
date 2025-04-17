package main

import (
	"os"

	"github.com/hydrocode-de/gorun/cli"
)

func main() {
	cli.Execute()
	os.Exit(0)

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

}
