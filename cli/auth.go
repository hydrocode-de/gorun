package cli

import (
	"fmt"
	"os"

	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	refreshToken bool
	accessToken  bool
	isAdmin      bool
	password     string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication and users for GoRun",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var credentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Show Admin credentials for GoRun",
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := auth.GetAdminCredentials(cmd.Context())
		cobra.CheckErr(err)

		if refreshToken {
			fmt.Println(credentials.RefreshToken)
			return
		}

		if accessToken {
			fmt.Println(credentials.AccessToken)
			return
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleColoredBlackOnBlueWhite)
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"Key", "Value"})
		t.AppendRow(table.Row{"Email", credentials.Email})
		t.AppendRow(table.Row{"User ID", credentials.UserID})
		t.AppendRow(table.Row{"Refresh token", credentials.RefreshToken})
		t.AppendRow(table.Row{"Expires at", credentials.ExpiresAt})
		t.AppendRow(table.Row{"Access token", credentials.AccessToken})

		t.Render()
	},
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Run: func(cmd *cobra.Command, args []string) {
		DB := viper.Get("db").(*db.Queries)
		secret := viper.GetString("secret")

		_, err := auth.CreateUser(cmd.Context(), DB, args[0], password, isAdmin, secret)
		cobra.CheckErr(err)
	},
}

func init() {
	credentialsCmd.Flags().BoolVar(&refreshToken, "refresh-token", false, "Show only the refresh token")
	credentialsCmd.Flags().BoolVar(&accessToken, "access-token", false, "Show only the access token")

	createUserCmd.Flags().BoolVar(&isAdmin, "admin", false, "Create an admin user")
	createUserCmd.Flags().StringVar(&password, "password", "", "The password for the new user. If you don't provide one, the user can't login.")

	authCmd.AddCommand(credentialsCmd)
	authCmd.AddCommand(createUserCmd)
	rootCmd.AddCommand(authCmd)
}
