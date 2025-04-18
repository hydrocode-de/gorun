package cli

import (
	"fmt"
	"strings"

	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var (
	listUsers bool
	password  string
	isAdmin   bool
	delete    bool
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage GoRun users",
	Run: func(cmd *cobra.Command, args []string) {
		DB := viper.Get("db").(*db.Queries)

		if listUsers {
			users, err := DB.GetAllUsers(cmd.Context())
			cobra.CheckErr(err)

			t := table.NewWriter()
			t.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
			t.AppendHeader(table.Row{"ID", "Email", "Is Admin", "Job count"})
			for _, user := range users {
				t.AppendRow(table.Row{user.ID, user.Email, user.IsAdmin, user.RunCount})
			}
			fmt.Println(t.Render())
			return
		}

		if len(args) == 0 {
			fmt.Println("No user provided. You need to provide a user id or email.\n")
			cmd.Help()
			return
		}

		var user db.User
		var err error
		if strings.Contains(args[0], "@") {
			user, err = DB.GetUserByEmail(cmd.Context(), args[0])
		} else {
			user, err = DB.GetUserByID(cmd.Context(), args[0])
		}
		cobra.CheckErr(err)
		if user == (db.User{}) {
			cobra.CheckErr(fmt.Errorf("user %s not found", args[0]))
		}

		if delete {
			err = DB.DeleteUser(cmd.Context(), user.ID)
			cobra.CheckErr(err)
			fmt.Println("User deleted successfully!")
			return
		}

		if password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			cobra.CheckErr(err)
			user, err = DB.UpdateUserPassword(cmd.Context(), db.UpdateUserPasswordParams{
				ID:           user.ID,
				PasswordHash: string(hashedPassword),
			})
			cobra.CheckErr(err)
			fmt.Println("Password updated!")
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
		t.AppendHeader(table.Row{"Key", "Value"})
		t.AppendRows([]table.Row{
			{"ID", user.ID},
			{"Email", user.Email},
			{"Is Admin", user.IsAdmin},
		})
		fmt.Println(t.Render())
	},
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cobra.CheckErr(fmt.Errorf("email is required"))
		}

		DB := viper.Get("db").(*db.Queries)
		secret := viper.GetString("secret")

		_, err := auth.CreateUser(cmd.Context(), DB, args[0], password, isAdmin, secret)
		cobra.CheckErr(err)
		fmt.Println("User created successfully!")
	},
}

func init() {
	userCmd.Flags().BoolVarP(&listUsers, "list", "l", false, "List all users")
	userCmd.Flags().StringVar(&password, "password", "", "Change the password for the selected user")
	userCmd.Flags().BoolVarP(&delete, "delete", "d", false, "Delete the selected user")

	createUserCmd.Flags().BoolVar(&isAdmin, "admin", false, "Create an admin user")
	createUserCmd.Flags().StringVar(&password, "password", "", "The password for the new user")

	userCmd.AddCommand(createUserCmd)
	rootCmd.AddCommand(userCmd)
}
