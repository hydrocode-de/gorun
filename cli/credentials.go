package cli

import (
	"fmt"
	"os"

	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	refreshToken bool
	accessToken  bool
)

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

func init() {
	credentialsCmd.Flags().BoolVar(&refreshToken, "refresh-token", false, "Show only the refresh token")
	credentialsCmd.Flags().BoolVar(&accessToken, "access-token", false, "Show only the access token")

	rootCmd.AddCommand(credentialsCmd)
}
