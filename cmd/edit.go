package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xnilsson/poflow/internal/config"
	"github.com/xnilsson/poflow/internal/editor"
)

var editFlags struct {
	dryRun bool
}

var editCmd = &cobra.Command{
	Use:   "edit <old-msgid> <new-msgid>",
	Short: "Update msgid across all .po and .pot files",
	Long: `Update a msgid (source text) across all language files and templates.

This command:
  1. Finds all .po files in your gettext directory
  2. Finds the .pot template file (if it exists)
  3. Updates the msgid in all matching entries
  4. Updates source code files that reference the msgid
  5. Preserves translations (msgstr) exactly
  6. Reports what was changed

Examples:
  # Update "Sign In" to "Log In" across all files and source code
  poflow edit "Sign In" "Log In"

  # Preview changes without modifying files
  poflow edit --dry-run "Sign In" "Log In"`,
	Args: cobra.ExactArgs(2),
	RunE: runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().BoolVar(&editFlags.dryRun, "dry-run", false, "show what would be changed without modifying files")
}

func runEdit(cmd *cobra.Command, args []string) error {
	oldMsgID := args[0]
	newMsgID := args[1]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Find all .po files
	poFiles, err := cfg.GetAllPOFiles()
	if err != nil {
		return fmt.Errorf("failed to find .po files: %w", err)
	}

	// Find .pot template (optional, might not exist)
	potFile, err := cfg.GetPOTFile()
	if err == nil {
		poFiles = append(poFiles, potFile)
	}

	if len(poFiles) == 0 {
		return fmt.Errorf("no .po or .pot files found in gettext directory")
	}

	if editFlags.dryRun {
		fmt.Println("DRY RUN - No files will be modified\n")
	}

	// Get current working directory as base for source file paths
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Update each file
	totalUpdated := 0
	totalEntries := 0

	for _, filePath := range poFiles {
		// Always update source files along with .po files
		result, err := editor.UpdateMsgIDInFileWithSources(filePath, oldMsgID, newMsgID, editFlags.dryRun, baseDir)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", filePath, err)
			continue
		}

		if result.EntriesFound > 0 {
			totalEntries += result.EntriesFound
			if result.Updated || editFlags.dryRun {
				totalUpdated++
				status := "✓"
				if editFlags.dryRun {
					status = "→"
				}
				fmt.Printf("  %s %s (%d entries)\n", status, filePath, result.EntriesFound)
			}
		}
	}

	// Summary
	fmt.Printf("\n")
	if totalEntries == 0 {
		fmt.Printf("No entries found matching \"%s\"\n", oldMsgID)
	} else if editFlags.dryRun {
		fmt.Printf("Would update %d file(s) with %d total entries\n", totalUpdated, totalEntries)
		fmt.Printf("\nRun without --dry-run to apply changes\n")
	} else {
		fmt.Printf("Updated %d file(s) with %d total entries\n", totalUpdated, totalEntries)
	}

	return nil
}
