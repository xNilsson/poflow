package editor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/xnilsson/poflow/internal/model"
	"github.com/xnilsson/poflow/internal/output"
	"github.com/xnilsson/poflow/internal/parser"
)

// UpdateResult tracks what was updated in a file
type UpdateResult struct {
	FilePath     string
	EntriesFound int
	Updated      bool
	Error        error
}

// UpdateMsgIDInFile updates msgid in a single .po file
func UpdateMsgIDInFile(filePath, oldMsgID, newMsgID string, dryRun bool) (*UpdateResult, error) {
	result := &UpdateResult{FilePath: filePath}

	// Open and parse file
	file, err := os.Open(filePath)
	if err != nil {
		result.Error = err
		return result, err
	}
	defer file.Close()

	p := parser.NewParser(file)

	// Track if we found any matches
	foundMatch := false
	var updatedEntries []*model.MsgEntry

	// Parse all entries
	for {
		entry := p.Next()
		if entry == nil {
			break
		}

		// Check if msgid matches
		if entry.MsgID == oldMsgID {
			foundMatch = true
			result.EntriesFound++

			if !dryRun {
				// Update msgid (preserve msgstr!)
				entry.MsgID = newMsgID
				// Update msgid in RawLines to preserve comment order
				updateMsgIDInRawLines(entry, newMsgID)
			}
		}

		updatedEntries = append(updatedEntries, entry)
	}

	if err := p.Err(); err != nil {
		result.Error = err
		return result, err
	}

	// If no matches, nothing to do
	if !foundMatch {
		return result, nil
	}

	// If dry run, just report what would be changed
	if dryRun {
		return result, nil
	}

	// Write to temp file
	tempFile, err := os.CreateTemp("", "poflow-edit-*.po")
	if err != nil {
		result.Error = err
		return result, err
	}
	tempFileName := tempFile.Name()
	defer os.Remove(tempFileName)

	writer := bufio.NewWriter(tempFile)

	// Write header
	for _, line := range p.Header() {
		writer.WriteString(line + "\n")
	}

	// Write all entries (updated ones have new msgid)
	for _, entry := range updatedEntries {
		writer.WriteString(output.FormatEntry(entry))
	}

	writer.Flush()
	tempFile.Close()

	// Replace original file
	if err := os.Rename(tempFileName, filePath); err != nil {
		result.Error = err
		return result, err
	}

	result.Updated = true
	return result, nil
}

// updateMsgIDInRawLines updates the msgid in the entry's RawLines while preserving all comments and formatting
func updateMsgIDInRawLines(entry *model.MsgEntry, newMsgID string) {
	if len(entry.RawLines) == 0 {
		return
	}

	var updatedLines []string
	inMsgID := false

	for _, line := range entry.RawLines {
		trimmed := strings.TrimSpace(line)

		// Detect start of msgid
		if strings.HasPrefix(trimmed, "msgid ") {
			inMsgID = true

			// Replace msgid line with new value
			if strings.Contains(newMsgID, "\n") {
				// Multi-line msgid
				updatedLines = append(updatedLines, "msgid \"\"")
				for _, msgLine := range strings.Split(newMsgID, "\n") {
					updatedLines = append(updatedLines, fmt.Sprintf("\"%s\\n\"", escapeString(msgLine)))
				}
			} else {
				// Single-line msgid
				updatedLines = append(updatedLines, fmt.Sprintf("msgid \"%s\"", escapeString(newMsgID)))
			}
			continue
		}

		// Detect start of msgstr (end of msgid section)
		if strings.HasPrefix(trimmed, "msgstr ") {
			inMsgID = false
		}

		// Skip continuation lines of msgid (they start with ")
		if inMsgID && strings.HasPrefix(trimmed, "\"") {
			continue
		}

		// Keep all other lines (comments, references, msgstr, etc.)
		updatedLines = append(updatedLines, line)
	}

	// Update the entry's RawLines
	entry.RawLines = updatedLines
}

// escapeString escapes special characters for .po file format
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// UpdateMsgIDInFileWithSources updates msgid in .po file AND in source code files
func UpdateMsgIDInFileWithSources(filePath, oldMsgID, newMsgID string, dryRun bool, baseDir string) (*UpdateResult, error) {
	result := &UpdateResult{FilePath: filePath}

	// Open and parse file
	file, err := os.Open(filePath)
	if err != nil {
		result.Error = err
		return result, err
	}
	defer file.Close()

	p := parser.NewParser(file)

	// Track if we found any matches and collect source file references
	foundMatch := false
	var updatedEntries []*model.MsgEntry
	var sourceFiles []string // Collect all source file references

	// Parse all entries
	for {
		entry := p.Next()
		if entry == nil {
			break
		}

		// Check if msgid matches
		if entry.MsgID == oldMsgID {
			foundMatch = true
			result.EntriesFound++

			// Collect source file references from this entry
			for _, ref := range entry.References {
				// Parse "file.ex:123" format - can have multiple space-separated refs
				parts := strings.Fields(ref)
				for _, part := range parts {
					// Split by colon to separate file path from line number
					colonIdx := strings.LastIndex(part, ":")
					if colonIdx > 0 {
						filePath := part[:colonIdx]
						sourceFiles = append(sourceFiles, filePath)
					}
				}
			}

			if !dryRun {
				// Update msgid (preserve msgstr!)
				entry.MsgID = newMsgID
				// Update msgid in RawLines to preserve comment order
				updateMsgIDInRawLines(entry, newMsgID)
			}
		}

		updatedEntries = append(updatedEntries, entry)
	}

	if err := p.Err(); err != nil {
		result.Error = err
		return result, err
	}

	// If no matches, nothing to do
	if !foundMatch {
		return result, nil
	}

	// If dry run, just report what would be changed
	if dryRun {
		return result, nil
	}

	// Update source files
	for _, sourceFile := range sourceFiles {
		fullPath := sourceFile
		if baseDir != "" {
			fullPath = filepath.Join(baseDir, sourceFile)
		}

		if err := updateSourceFile(fullPath, oldMsgID, newMsgID); err != nil {
			// Don't fail the whole operation if a source file update fails
			// Just log and continue
			fmt.Printf("  Warning: failed to update source file %s: %v\n", sourceFile, err)
		}
	}

	// Write to temp file
	tempFile, err := os.CreateTemp("", "poflow-edit-*.po")
	if err != nil {
		result.Error = err
		return result, err
	}
	tempFileName := tempFile.Name()
	defer os.Remove(tempFileName)

	writer := bufio.NewWriter(tempFile)

	// Write header
	for _, line := range p.Header() {
		writer.WriteString(line + "\n")
	}

	// Write all entries (updated ones have new msgid)
	for _, entry := range updatedEntries {
		writer.WriteString(output.FormatEntry(entry))
	}

	writer.Flush()
	tempFile.Close()

	// Replace original file
	if err := os.Rename(tempFileName, filePath); err != nil {
		result.Error = err
		return result, err
	}

	result.Updated = true
	return result, nil
}

// updateSourceFile updates a single source file, replacing oldMsgID with newMsgID in gettext calls
func updateSourceFile(filePath, oldMsgID, newMsgID string) error {
	// Read the source file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	contentStr := string(content)

	// Replace in various gettext function calls
	// We use a simple string replacement approach:
	// Find: "oldMsgID"
	// Replace: "newMsgID"
	// This works because gettext strings are always quoted

	// Escape the old string for regex
	oldQuoted := `"` + regexp.QuoteMeta(oldMsgID) + `"`
	newQuoted := `"` + newMsgID + `"`

	// Use regex to find quoted strings
	re := regexp.MustCompile(oldQuoted)
	updatedContent := re.ReplaceAllString(contentStr, newQuoted)

	// If content changed, write it back
	if updatedContent != contentStr {
		if err := os.WriteFile(filePath, []byte(updatedContent), 0644); err != nil {
			return fmt.Errorf("failed to write source file: %w", err)
		}
	}

	return nil
}
