package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nille/poflow/internal/model"
	"github.com/nille/poflow/internal/parser"
	"github.com/spf13/cobra"
)

var searchvalueCmd = &cobra.Command{
	Use:   "searchvalue PATTERN [file]",
	Short: "Search for entries by msgstr pattern",
	Long: `Search for translation entries where the msgstr (translation) matches the given pattern.

By default, uses plain substring matching (case-insensitive).
Use --re for regex pattern matching.

Examples:
  poflow searchvalue "Välkommen" file.po
  poflow searchvalue --re "^Tack" file.po
  poflow searchvalue --json "fel" file.po
  cat file.po | poflow searchvalue "error"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runSearchValue,
}

var searchvalueFlags struct {
	useRegex bool
	limit    int
}

func init() {
	rootCmd.AddCommand(searchvalueCmd)
	searchvalueCmd.Flags().BoolVar(&searchvalueFlags.useRegex, "re", false, "use regex pattern matching")
	searchvalueCmd.Flags().IntVar(&searchvalueFlags.limit, "limit", 0, "maximum number of entries to output (0 = no limit)")
}

func runSearchValue(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	// Compile regex if needed
	var re *regexp.Regexp
	var err error
	if searchvalueFlags.useRegex {
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
	if len(args) == 2 {
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

		// Check if msgstr matches pattern
		matches := false
		if searchvalueFlags.useRegex {
			matches = re.MatchString(entry.MsgStr)
		} else {
			// Case-insensitive substring match
			matches = strings.Contains(strings.ToLower(entry.MsgStr), pattern)
		}

		if !matches {
			continue
		}

		// Output the entry
		if jsonOutput {
			outputJSONValue(*entry)
		} else {
			outputTextValue(*entry)
		}

		count++
		if searchvalueFlags.limit > 0 && count >= searchvalueFlags.limit {
			break
		}
	}

	if err := p.Err(); err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	return nil
}

func outputJSONValue(entry model.MsgEntry) {
	data, _ := json.Marshal(entry)
	fmt.Println(string(data))
}

func outputTextValue(entry model.MsgEntry) {
	// Output comments and references
	for _, comment := range entry.Comments {
		fmt.Println(comment)
	}
	for _, ref := range entry.References {
		fmt.Printf("#: %s\n", ref)
	}

	// Output msgid and msgstr
	fmt.Printf("msgid \"%s\"\n", entry.MsgID)
	fmt.Printf("msgstr \"%s\"\n", entry.MsgStr)
	fmt.Println()
}
