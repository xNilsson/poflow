package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is set during build
	Version = "0.1.0-dev"
	// Commit is set during build
	Commit = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of poflow",
	Long:  `All software has versions. This is poflow's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("poflow version %s (commit: %s)\n", Version, Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
