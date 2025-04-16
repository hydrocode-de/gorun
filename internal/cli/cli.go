package cli

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hydrocode-de/gorun/internal/cache"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/helper"
	"github.com/hydrocode-de/gorun/sql"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const banner = `
   _____       _____             
  / ____|     |  __ \            
 | |  __  ___ | |__) |     _ __   
 | | |_ |/ _ \|  _  / | | |  _ \  
 | |__| | (_) | | \ \ |_| | | | | 
  \_____|\___/|_|  \_\__,_|_| |_| 
`

var rootCmd = &cobra.Command{
	Use:   "gorun",
	Short: "GoRun operates tool-spec compliant research tools",
	Long: banner + `
GoRun is a CLI tool that operates tool-spec compliant research tools.

The tool specification is available at https://voforwater.github.io/tool-spec/
You ran gorun without a command. Please refer to the section below to learn
about all available commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println(viper.AllKeys())
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initApplicationConfig)

}

func initApplicationConfig() {
	// Load .env file first
	godotenv.Load()

	viper.SetEnvPrefix("gorun")
	viper.AutomaticEnv()

	viper.SetDefault("port", 8080)
	viper.SetDefault("path", path.Join(os.Getenv("HOME"), "gorun"))
	viper.SetDefault("db_path", path.Join(viper.GetString("path"), "gorun.db"))
	viper.SetDefault("mount_path", path.Join(viper.GetString("path"), "mounts"))
	viper.SetDefault("temp_dir", path.Join(os.TempDir(), "gorun"))
	viper.SetDefault("max_upload_size", 1024*1024*1024*2) // 2GB
	viper.SetDefault("max_temp_age", 12*time.Hour)
	viper.SetDefault("secret", helper.GetRandomString(32))

	c := &cache.Cache{}
	c.Reset()
	viper.Set("cache", c)

	err := os.MkdirAll(viper.GetString("path"), 0755)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to create GorunBasePath directory: %w", err))
	}

	// Ensure the database directory exists
	dbDir := path.Dir(viper.GetString("db_path"))
	err = os.MkdirAll(dbDir, 0755)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to create database directory: %w", err))
	}

	// Ensure the mount directory exists
	err = os.MkdirAll(viper.GetString("mount_path"), 0755)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to create mount directory: %w", err))
	}

	// Ensure the temp directory exists
	err = os.MkdirAll(viper.GetString("temp_dir"), 0755)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to create temp directory: %w", err))
	}

	// Initialize the database driver
	drv, err := sql.CreateDB(viper.GetString("db_path"))
	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to create database driver: %w", err))
	}
	dbQueries := db.New(drv)
	viper.Set("db", dbQueries)
}
