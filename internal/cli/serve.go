package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port int
var host string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the GoRun API server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("DBUG: serve")
	},
}

func init() {
	serveCmd.Flags().IntVar(&port, "port", 8080, "The port to listen on")
	serveCmd.Flags().StringVar(&host, "host", "127.0.0.1", "The host to listen on")

	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))

	rootCmd.AddCommand(serveCmd)
}
