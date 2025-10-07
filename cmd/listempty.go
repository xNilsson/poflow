package cmd

import (
	"fmt"
	"os"

	"github.com/xnilsson/poflow/internal/config"
	"github.com/xnilsson/poflow/internal/output"
	"github.com/xnilsson/poflow/internal/parser"
	"github.com/spf13/cobra"
)

var (
	listEmptyLimit    int
	listEmptyLanguage string
)

var listemptyCmd = &cobra.Command{
	Use:   "listempty [file]",
	Short: "List untranslated entries",
	Long: `List all entries with empty translations (msgstr).

Examples:
  poflow listempty file.po
  poflow listempty --json file.po
  poflow listempty --limit 10 file.po
  cat file.po | poflow listempty
  poflow listempty --language sv --json`,
	RunE: runListEmpty,
}

func init() {
	rootCmd.AddCommand(listemptyCmd)
	listemptyCmd.Flags().IntVar(&listEmptyLimit, "limit", 0, "limit number of entries (0 = no limit)")
	listemptyCmd.Flags().StringVar(&listEmptyLanguage, "language", "", "language code (uses config to resolve path)")
}

func runListEmpty(cmd *cobra.Command, args []string) error {
	// Determine input source
	var reader *os.File
	var err error

	// Handle --language flag
	if listEmptyLanguage != "" {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		path, err := cfg.ResolvePOPath(listEmptyLanguage)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}
		reader, err = os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer reader.Close()
	} else if len(args) > 0 {
		// Read from file
		reader, err = os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer reader.Close()
	} else {
		// Read from stdin
		reader = os.Stdin
	}

	// Create parser
	p := parser.NewParser(reader)

	// Get output format from global flag
	jsonOutput, _ := cmd.Flags().GetBool("json")
	count := 0

	// Stream entries
	for {
		entry := p.Next()
		if entry == nil {
			break
		}

		// Skip non-empty entries
		if !entry.IsEmpty() {
			continue
		}

		// Check limit
		if listEmptyLimit > 0 && count >= listEmptyLimit {
			break
		}

		// Output entry
		if err := output.OutputEntry(entry, jsonOutput); err != nil {
			return err
		}

		count++
	}

	// Check for errors
	if err := p.Err(); err != nil {
		return fmt.Errorf("parsing error: %w", err)
	}

	return nil
}
