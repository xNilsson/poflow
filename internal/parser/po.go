package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/nille/poflow/internal/model"
)

// Parser streams .po file entries one by one without loading entire file into memory
type Parser struct {
	scanner *bufio.Scanner
	err     error
}

// NewParser creates a new streaming parser for .po files
func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
	}
}

// Next returns the next entry from the .po file, or nil when done
func (p *Parser) Next() *model.MsgEntry {
	if p.err != nil {
		return nil
	}

	var entry model.MsgEntry
	var msgidLines []string
	var msgstrLines []string
	inMsgID := false
	inMsgStr := false

	for p.scanner.Scan() {
		line := p.scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip empty lines between entries
		if trimmed == "" {
			if entry.MsgID != "" {
				// We have a complete entry
				entry.MsgStr = strings.Join(msgstrLines, "")
				return &entry
			}
			// Reset state if we hit empty line without msgid
			inMsgID = false
			inMsgStr = false
			continue
		}

		// Handle comments
		if strings.HasPrefix(trimmed, "#:") {
			// Reference comment
			ref := strings.TrimSpace(trimmed[2:])
			entry.References = append(entry.References, ref)
			continue
		} else if strings.HasPrefix(trimmed, "#") {
			// Other comments
			comment := strings.TrimSpace(trimmed[1:])
			entry.Comments = append(entry.Comments, comment)
			continue
		}

		// Handle msgid
		if strings.HasPrefix(trimmed, "msgid ") {
			inMsgID = true
			inMsgStr = false
			msgidLines = []string{unquote(trimmed[6:])}
			continue
		}

		// Handle msgstr
		if strings.HasPrefix(trimmed, "msgstr ") {
			// Save accumulated msgid
			entry.MsgID = strings.Join(msgidLines, "")
			inMsgID = false
			inMsgStr = true
			msgstrLines = []string{unquote(trimmed[7:])}
			continue
		}

		// Handle continuation lines (quoted strings on their own lines)
		if strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"") {
			if inMsgID {
				msgidLines = append(msgidLines, unquote(trimmed))
			} else if inMsgStr {
				msgstrLines = append(msgstrLines, unquote(trimmed))
			}
			continue
		}
	}

	// Handle last entry in file (no trailing empty line)
	if len(msgidLines) > 0 {
		entry.MsgID = strings.Join(msgidLines, "")
		entry.MsgStr = strings.Join(msgstrLines, "")
		return &entry
	}

	// Check for scanner errors
	if err := p.scanner.Err(); err != nil {
		p.err = err
	}

	return nil
}

// Err returns any error encountered during parsing
func (p *Parser) Err() error {
	return p.err
}

// unquote removes surrounding quotes and handles escape sequences
func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// Handle common escape sequences
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\t", "\t")
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "\\\\", "\\")

	return s
}

// ParseAll reads all entries from a .po file (convenience method for testing)
func ParseAll(r io.Reader) ([]*model.MsgEntry, error) {
	parser := NewParser(r)
	var entries []*model.MsgEntry

	for {
		entry := parser.Next()
		if entry == nil {
			break
		}
		entries = append(entries, entry)
	}

	if err := parser.Err(); err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return entries, nil
}
