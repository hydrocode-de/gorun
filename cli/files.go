package cli

import (
	"fmt"
	"log"

	"github.com/hydrocode-de/gorun/internal/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type targetFlag files.Target

func (t *targetFlag) Set(value string) error {
	*t = targetFlag(value)
	return nil
}

func (t *targetFlag) Type() string {
	return "target"
}

func (t *targetFlag) String() string {
	return string(*t)
}

var location targetFlag = targetFlag(files.TargetAll)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Manage input and output files",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var findFileCmd = &cobra.Command{
	Use:   "find",
	Short: "Find file(s) by name or filename pattern",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("missing arguments - use like gorun files find <pattern>")
		}

		mountPath := viper.GetString("mount_path")
		matches, err := files.Find(args[0], mountPath, files.Target(location))
		cobra.CheckErr(err)

		for _, match := range matches {
			fmt.Printf("%s\n", match.AbsPath)
		}

	},
}

func init() {
	findFileCmd.Flags().Var(&location, "location", "The location to search for files. Can be 'all', 'input', 'output' or 'both'")

	findFileCmd.Flags().String("mount-path", "", "The mount base path to search for files")
	viper.BindPFlag("mount_path", findFileCmd.Flags().Lookup("mount-path"))

	filesCmd.AddCommand(findFileCmd)
	rootCmd.AddCommand(filesCmd)
}
