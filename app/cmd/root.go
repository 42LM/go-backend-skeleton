// Package cmd holds the CLI command definitions.
//
// This file was partially generated with spf13/Cobra 🐍,
// read more about the cobra generator:
//   - https://github.com/spf13/cobra-cli/blob/main/README.md
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// used to inject value at build time using linker flag.
//
// go build -ldflags "-X go-backend-skeleton/app/cmd.version=v0.0.0" -o gbs
var version string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     "gbs",
	Version: version,
	Short:   "gbs - go backend skeleton",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("exit:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
