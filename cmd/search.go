package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nille/poflow/internal/config"
	"github.com/nille/poflow/internal/output"
	"github.com/nille/poflow/internal/parser"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search PATTERN [file]",
	Short: "Search for entries by msgid pattern",
	Long: `Search for translation entries where the msgid matches the given pattern.

By default, uses plain substring matching (case-insensitive).
Use --re for regex pattern matching.

Examples:
  poflow search "Welcome" file.po
  poflow search --re "^Login" file.po
  poflow search --json "error" file.po
  cat file.po | poflow search "Welcome"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runSearch,
}

var searchFlags struct {
	useRegex bool
	limit    int
	language string
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&searchFlags.useRegex, "re", false, "use regex pattern matching")
	searchCmd.Flags().IntVar(&searchFlags.limit, "limit", 0, "maximum number of entries to output (0 = no limit)")
	searchCmd.Flags().StringVar(&searchFlags.language, "language", "", "language code (uses config to resolve path)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	// Compile regex if needed
	var re *regexp.Regexp
	var err error
	if searchFlags.useRegex {
		re, err = regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	} else {
		// For plain matching, convert to lowercase for case-insensitive search
		pattern = strings.ToLower(pattern)
	}

	// Determine input source
	var reader *os.File

	// Handle --language flag
	if searchFlags.language != "" {
		cfg, cfgErr := config.Load()
		if cfgErr != nil {
			return fmt.Errorf("failed to load config: %w", cfgErr)
		}
		path, pathErr := cfg.ResolvePOPath(searchFlags.language)
		if pathErr != nil {
			return fmt.Errorf("failed to resolve path: %w", pathErr)
		}
		var openErr error
		reader, openErr = os.Open(path)
		if openErr != nil {
			return fmt.Errorf("failed to open file: %w", openErr)
		}
		defer reader.Close()
	} else if len(args) == 2 {
		// File input
		reader, err = os.Open(args[1])
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer reader.Close()
	} else {
		// Stdin input
		reader = os.Stdin
	}

	// Create parser
	p := parser.NewParser(reader)

	// Get JSON output flag
	jsonOutput, _ := cmd.Flags().GetBool("json")

	// Process entries
	count := 0
	for {
		entry := p.Next()
		if entry == nil {
			break
		}

		// Check if msgid matches pattern
		matches := false
		if searchFlags.useRegex {
			matches = re.MatchString(entry.MsgID)
		} else {
			// Case-insensitive substring match
			matches = strings.Contains(strings.ToLower(entry.MsgID), pattern)
		}

		if !matches {
			continue
		}

		// Output the entry
		if err := output.OutputEntry(entry, jsonOutput); err != nil {
			return err
		}

		count++
		if searchFlags.limit > 0 && count >= searchFlags.limit {
			break
		}
	}

	if err := p.Err(); err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	return nil
}

