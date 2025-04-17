package cli

import (
	"fmt"

	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/files"
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
		tools, err := toolImage.ReadAllTools(cmd.Context(), cache, viper.GetBool("verbose"))
		cobra.CheckErr(err)

		fmt.Printf("Found %d tools:\n", len(tools))
		for _, name := range tools {
			fmt.Printf("|- %s\n", name)
		}
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up temporary files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running cleanup...")
		err := files.Cleanup()
		cobra.CheckErr(err)
		fmt.Println("Cleanup completed successfully")
	},
}

func init() {
	listCmd.Flags().BoolVar(&verbose, "verbose", false, "Verbose output")
	viper.BindPFlag("verbose", listCmd.Flags().Lookup("verbose"))

	toolsCmd.AddCommand(listCmd)
	toolsCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(toolsCmd)
}
