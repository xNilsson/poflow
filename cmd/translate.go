package cmd

import (
	"fmt"
	"os"

	"github.com/nille/poflow/internal/config"
	"github.com/nille/poflow/internal/output"
	"github.com/nille/poflow/internal/parser"
	"github.com/spf13/cobra"
)

var translateFlags struct {
	language string
	force    bool
}

var translateCmd = &cobra.Command{
	Use:   "translate [po-file]",
	Short: "Merge translations into a .po file",
	Long: `Merge translations into a .po file from translation input.

Translation input format (one per line):
  msgid = msgstr

Examples:
  Sign In = Logga in
  Sign Out = Logga ut
  Welcome = Välkommen

Usage patterns:

  # Config-based (uses poflow.yml to resolve path)
  poflow translate --language sv translations.txt

  # Direct file path
  cat translations.txt | poflow translate file.po > file_new.po

  # From stdin translations
  echo "Sign In = Logga in" | poflow translate --language sv

Config file format (poflow.yml):
  gettext_path: "priv/gettext"

Resolves to: {gettext_path}/{lang}/LC_MESSAGES/default.po`,
	RunE: runTranslate,
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringVarP(&translateFlags.language, "language", "l", "", "language code (e.g., sv, en)")
	translateCmd.Flags().BoolVarP(&translateFlags.force, "force", "f", false, "continue even if msgids not found")
}

func runTranslate(cmd *cobra.Command, args []string) error {
	jsonOutput, _ := cmd.Flags().GetBool("json")
	quiet, _ := cmd.Flags().GetBool("quiet")

	// Determine the .po file to translate
	var poFilePath string
	if translateFlags.language != "" {
		// Config-based path resolution
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w (hint: create poflow.yml with gettext_path)", err)
		}
		poFilePath, err = cfg.ResolvePOPath(translateFlags.language)
		if err != nil {
			return err
		}
		if !quiet {
			fmt.Fprintf(os.Stderr, "Resolved path: %s\n", poFilePath)
		}
	} else if len(args) > 0 {
		// Direct file path
		poFilePath = args[0]
	} else {
		return fmt.Errorf("either --language or po-file argument is required")
	}

	// Determine translation input source
	var translationInput *os.File
	if len(args) > 1 {
		// Translation file provided as second argument
		f, err := os.Open(args[1])
		if err != nil {
			return fmt.Errorf("failed to open translation file: %w", err)
		}
		defer f.Close()
		translationInput = f
	} else {
		// Read translations from stdin
		translationInput = os.Stdin
	}

	// Parse translations
	translations, err := parser.ParseTranslations(translationInput)
	if err != nil {
		return fmt.Errorf("failed to parse translations: %w", err)
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "Loaded %d translations\n", len(translations))
	}

	// Open and read the .po file
	poFile, err := os.Open(poFilePath)
	if err != nil {
		return fmt.Errorf("failed to open .po file: %w", err)
	}
	defer poFile.Close()

	// Parse and merge
	p := parser.NewParser(poFile)
	notFound := []string{}
	updated := 0

	for {
		entry := p.Next()
		if entry == nil {
			break
		}

		// Check if we have a translation for this msgid
		if newMsgStr, ok := translations[entry.MsgID]; ok {
			entry.MsgStr = newMsgStr
			updated++
			delete(translations, entry.MsgID) // Mark as found
		}

		// Output the entry (possibly updated)
		if err := output.OutputEntry(entry, jsonOutput); err != nil {
			return err
		}
	}

	if err := p.Err(); err != nil {
		return fmt.Errorf("error parsing .po file: %w", err)
	}

	// Check for unfound translations
	if len(translations) > 0 {
		for msgid := range translations {
			notFound = append(notFound, msgid)
		}

		if !quiet {
			fmt.Fprintf(os.Stderr, "\nWarning: %d msgid(s) not found in .po file:\n", len(notFound))
			for _, msgid := range notFound {
				fmt.Fprintf(os.Stderr, "  - %s\n", msgid)
			}
		}

		if !translateFlags.force {
			return fmt.Errorf("some translations not applied (use --force to ignore)")
		}
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "\nUpdated %d entries\n", updated)
	}

	return nil
}

