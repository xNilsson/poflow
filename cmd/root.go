package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "poflow",
	Short: "A workflow utility for GNU gettext .po files",
	Long: `poflow: workflow utility for gettext .po files

poflow helps developers, translators, and LLMs navigate large .po files —
searching, listing, and updating translation entries in a structured,
automatable way.

Features:
  • Fast streaming parser (no in-memory AST needed)
  • Search by msgid or msgstr with regex support
  • List untranslated entries
  • Merge translations from text files
  • JSON output for programmatic/LLM usage
  • Config file support for project-specific paths`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./poflow.yml or ~/.config/poflow/config.yml)")
	rootCmd.PersistentFlags().Bool("json", false, "output in JSON format")
	rootCmd.PersistentFlags().Bool("quiet", false, "suppress progress output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory first
		viper.AddConfigPath(".")
		viper.SetConfigName("poflow")

		// Search in home directory config folder
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home + "/.config/poflow")
		}
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if !viper.GetBool("quiet") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
