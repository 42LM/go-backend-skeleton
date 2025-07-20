package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionCmd does what one would expect.
// It prints the version of the go backend skeleton.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. This is go backend skeletons.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gbs version " + version)
	},
}
