package cli

import (
	"fmt"

	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/toolImage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Manage tools from the command line",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tools",
	Run: func(cmd *cobra.Command, args []string) {
		cache := viper.Get("cache").(*cache.Cache)
		tools, err := toolImage.ReadAllTools(cmd.Context(), cache)
		cobra.CheckErr(err)

		fmt.Printf("Found %d tools:\n", len(tools))
		for _, name := range tools {
			fmt.Printf("|- %s\n", name)
		}
	},
}

func init() {
	toolsCmd.AddCommand(listCmd)
	rootCmd.AddCommand(toolsCmd)
}
