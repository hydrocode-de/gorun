package cli

import (
	"fmt"
	"strings"

	"github.com/hydrocode-de/gorun/internal/toolImage"
	"github.com/spf13/cobra"
)

var image string
var verbose bool

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a docker image to be tool-spec compliant",
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := toolImage.ReadToolSpec(cmd.Context(), image)
		if err != nil {
			fmt.Println("The tool image is not tool-spec compliant")
		} else {
			fmt.Println("The tool image is tool-spec compliant")
		}
		if verbose {
			cobra.CheckErr(err)
			fmt.Printf("\nNumber of Tools: %d\n", len(spec.Tools))
			for _, tool := range spec.Tools {
				fmt.Printf("\n- %s\n", tool.Title)
				desc := tool.Description
				if len(desc) > 70 {
					desc = desc[:70] + "..."
				}
				fmt.Printf("  Description: %s\n", desc)
				fmt.Printf("  Parameters: %d (", len(tool.Parameters))
				names := make([]string, 0, len(tool.Parameters))
				for name := range tool.Parameters {
					names = append(names, name)
				}
				fmt.Printf("%s", strings.Join(names, ", "))
				fmt.Printf(")\n")
				fmt.Println()
			}
		}
	},
}

func init() {
	inspectCmd.Flags().StringVar(&image, "image", "", "The image to inspect")
	inspectCmd.MarkFlagRequired("image")
	inspectCmd.Flags().BoolVar(&verbose, "verbose", false, "Verbose output")

	rootCmd.AddCommand(inspectCmd)
}
