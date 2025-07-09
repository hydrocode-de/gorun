package cli

import (
	"fmt"
	"os"

	"github.com/hydrocode-de/gorun/api/mcp"
	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP (Model Control Protocol) commands",
	Long:  `Commands for interacting with the Model Control Protocol (MCP) server.`,
}

var mcpStdioCmd = &cobra.Command{
	Use:   "stdio",
	Short: "Start MCP server with stdio transport",
	Long: `Start the MCP (Model Control Protocol) server using standard input/output for communication.
This mode is typically used when integrating with IDEs or other tools that need to
communicate with the MCP server via stdio.`,
	Run: func(cmd *cobra.Command, args []string) {
		cache := viper.Get("cache").(*cache.Cache)
		transport := mcp.NewStdioTransport(cache)
		if err := transport.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running MCP stdio server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	mcpCmd.AddCommand(mcpStdioCmd)
	rootCmd.AddCommand(mcpCmd)
}
