package cmd

import (
	"bufio"
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
	stdout   bool
}

var translateCmd = &cobra.Command{
	Use:   "translate [po-file] [translation-file]",
	Short: "Merge translations into a .po file",
	Long: `Merge translations into a .po file from translation input.

Translation input format (one per line):
  msgid = msgstr

Examples:
  Sign In = Logga in
  Sign Out = Logga ut
  Welcome = Välkommen

BEHAVIOR:

  By default, the .po file is updated IN-PLACE and a summary is shown:

    $ poflow translate --language sv translations.txt
    Resolved path: priv/gettext/sv/LC_MESSAGES/default.po
    Loaded 2 translations

    Updated 2 translation(s) in priv/gettext/sv/LC_MESSAGES/default.po:
      ✓ Ask a question
      ✓ Feedback

  To output to stdout instead (for piping), use --stdout:

    $ poflow translate --language sv translations.txt --stdout > output.po

USAGE PATTERNS:

  # Config-based (uses poflow.yml to resolve path)
  poflow translate --language sv translations.txt

  # Direct file path
  poflow translate file.po translations.txt

  # From stdin translations
  echo "Sign In = Logga in" | poflow translate --language sv

  # Output to stdout for piping
  poflow translate --language sv translations.txt --stdout > new.po

Config file format (poflow.yml):
  gettext_path: "priv/gettext"

Resolves to: {gettext_path}/{lang}/LC_MESSAGES/default.po`,
	RunE: runTranslate,
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringVarP(&translateFlags.language, "language", "l", "", "language code (e.g., sv, en)")
	translateCmd.Flags().BoolVarP(&translateFlags.force, "force", "f", false, "continue even if msgids not found")
	translateCmd.Flags().BoolVar(&translateFlags.stdout, "stdout", false, "output to stdout instead of updating file in-place")
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
	var translationArg string

	if translateFlags.language != "" {
		// When using --language, the first arg (if present) is the translation file
		if len(args) > 0 {
			translationArg = args[0]
		}
	} else {
		// When using direct path, the second arg (if present) is the translation file
		if len(args) > 1 {
			translationArg = args[1]
		}
	}

	if translationArg != "" {
		f, err := os.Open(translationArg)
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
	updatedMsgIDs := []string{}

	// Determine output destination
	var outputWriter *bufio.Writer
	var tempFile *os.File

	if translateFlags.stdout {
		// Output to stdout
		outputWriter = bufio.NewWriter(os.Stdout)
	} else {
		// Write to temp file for in-place update
		var err error
		tempFile, err = os.CreateTemp("", "poflow-*.po")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tempFile.Name())
		outputWriter = bufio.NewWriter(tempFile)
	}

	// Write file header first
	headerWritten := false

	for {
		entry := p.Next()
		if entry == nil {
			break
		}

		// Write header before first entry
		if !headerWritten {
			header := p.Header()
			for _, line := range header {
				outputWriter.WriteString(line + "\n")
			}
			headerWritten = true
		}

		// Check if we have a translation for this msgid
		if newMsgStr, ok := translations[entry.MsgID]; ok {
			entry.MsgStr = newMsgStr
			updated++
			updatedMsgIDs = append(updatedMsgIDs, entry.MsgID)
			delete(translations, entry.MsgID) // Mark as found
		}

		// Output the entry (possibly updated)
		if translateFlags.stdout && jsonOutput {
			if err := output.OutputEntry(entry, jsonOutput); err != nil {
				return err
			}
		} else {
			// Write as .po format
			if _, err := outputWriter.WriteString(output.FormatEntry(entry)); err != nil {
				return fmt.Errorf("failed to write entry: %w", err)
			}
		}
	}

	if err := p.Err(); err != nil {
		return fmt.Errorf("error parsing .po file: %w", err)
	}

	// Flush output
	if err := outputWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	// If in-place mode, replace original file with temp file
	if !translateFlags.stdout {
		if err := tempFile.Close(); err != nil {
			return fmt.Errorf("failed to close temp file: %w", err)
		}
		if err := os.Rename(tempFile.Name(), poFilePath); err != nil {
			return fmt.Errorf("failed to replace original file: %w", err)
		}
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

	// Show summary (unless quiet or stdout mode with non-JSON output)
	if !quiet && !translateFlags.stdout {
		fmt.Fprintf(os.Stderr, "\nUpdated %d translation(s) in %s:\n", updated, poFilePath)
		for _, msgid := range updatedMsgIDs {
			fmt.Fprintf(os.Stderr, "  ✓ %s\n", msgid)
		}
	} else if !quiet && translateFlags.stdout {
		fmt.Fprintf(os.Stderr, "\nUpdated %d entries\n", updated)
	}

	return nil
}

