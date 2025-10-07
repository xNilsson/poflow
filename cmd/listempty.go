package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nille/poflow/internal/parser"
	"github.com/spf13/cobra"
)

var (
	listEmptyLimit int
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
}

func runListEmpty(cmd *cobra.Command, args []string) error {
	// Determine input source
	var reader *os.File
	var err error

	if len(args) > 0 {
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
		if jsonOutput {
			data, err := json.Marshal(entry)
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(data))
		} else {
			fmt.Printf("msgid \"%s\"\n", entry.MsgID)
			fmt.Printf("msgstr \"%s\"\n", entry.MsgStr)
			if len(entry.References) > 0 {
				for _, ref := range entry.References {
					fmt.Printf("#: %s\n", ref)
				}
			}
			fmt.Println()
		}

		count++
	}

	// Check for errors
	if err := p.Err(); err != nil {
		return fmt.Errorf("parsing error: %w", err)
	}

	return nil
}
