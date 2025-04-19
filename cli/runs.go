package cli

import (
	"fmt"

	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	listRuns bool
	filter   string
)

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Manage job runs",
	Run: func(cmd *cobra.Command, args []string) {
		DB := viper.Get("db").(*db.Queries)
		credentials, err := auth.GetAdminCredentials(cmd.Context())
		cobra.CheckErr(err)

		if listRuns {
			var runs []db.Run

			switch filter {
			case "":
				runs, err = DB.GetAllRuns(cmd.Context(), db.GetAllRunsParams{
					UserID: credentials.UserID,
				})
			}
			cobra.CheckErr(err)

			t := table.NewWriter()
			t.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
			t.AppendHeader(table.Row{"ID", "Name", "Title", "Created", "Status"})
			for _, run := range runs {
				t.AppendRow(table.Row{run.ID, run.Name, run.Title, run.CreatedAt, run.Status})
			}
			fmt.Println(t.Render())
			return
		}

		fmt.Printf(banner)
		cmd.Help()
	},
}

func init() {
	runsCmd.Flags().BoolVarP(&listRuns, "list", "l", false, "List all available runs")

	rootCmd.AddCommand(runsCmd)
}
