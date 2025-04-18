package cli

import (
	"fmt"

	"github.com/hydrocode-de/gorun/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GoRun",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoRun version %s (commit: %s, built: %s)\n", version.Version, version.Commit, version.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
