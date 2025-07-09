package cli

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hydrocode-de/gorun/api"
	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/files"
	"github.com/hydrocode-de/gorun/internal/toolImage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port int
var host string
var noAuth bool
var enableMcp bool
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the GoRun API server",
	Run: func(cmd *cobra.Command, args []string) {
		if noAuth && host != "127.0.0.1" {
			logrus.Warn("You are running the server with no authentication and a non-localhost host. This is not recommended and might expose your server to the public internet.")
		}
		startBackgroundTasks(cmd.Context())

		mux, err := api.CreateServer(enableMcp)
		cobra.CheckErr(err)

		server := api.EnableCORS(mux, "*")
		logrus.Infof("GoRun server listening on  http://%s:%d", host, port)
		logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), server))
	},
}

func startBackgroundTasks(ctx context.Context) {
	cleanupTicker := time.NewTicker(time.Minute * 5)
	go func() {
		for range cleanupTicker.C {
			logrus.Info("Running cleanup")
			err := files.Cleanup()
			cobra.CheckErr(err)
		}
	}()

	toolsTicker := time.NewTicker(time.Minute * 5)
	go func() {
		for range toolsTicker.C {
			logrus.Info("Checking for new tools")
			cache := viper.Get("cache").(*cache.Cache)
			_, err := toolImage.ReadAllTools(ctx, cache, false)
			cobra.CheckErr(err)
		}
	}()

	adminTicker := time.NewTicker(time.Minute * 50)
	go func() {
		for range adminTicker.C {
			logrus.Debug("Renewing admin credentials")
			if _, err := auth.GetAdminCredentials(ctx); err != nil {
				logrus.Errorf("Failed to renew admin credentials: %v", err)
			}
		}
	}()
}

func init() {
	serveCmd.Flags().IntVar(&port, "port", 8080, "The port to listen on")
	serveCmd.Flags().StringVar(&host, "host", "127.0.0.1", "The host to listen on")
	serveCmd.Flags().BoolVar(&noAuth, "no-auth", false, "Disable authentication")
	serveCmd.Flags().BoolVar(&enableMcp, "enable-mcp", true, "Enable MCP")
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("host", serveCmd.Flags().Lookup("host"))
	viper.BindPFlag("no_auth", serveCmd.Flags().Lookup("no-auth"))
	viper.BindPFlag("enable_mcp", serveCmd.Flags().Lookup("enable-mcp"))

	// Configure logrus
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	rootCmd.AddCommand(serveCmd)
}
