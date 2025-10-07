package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Translation represents a single translation pair
type Translation struct {
	MsgID  string
	MsgStr string
}

// ParseTranslations parses translation input in the format: msgid = msgstr
// Returns a map of msgid -> msgstr for fast lookups
func ParseTranslations(r io.Reader) (map[string]string, error) {
	translations := make(map[string]string)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse "msgid = msgstr" format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format, expected 'msgid = msgstr', got: %s", lineNum, line)
		}

		msgid := strings.TrimSpace(parts[0])
		msgstr := strings.TrimSpace(parts[1])

		if msgid == "" {
			return nil, fmt.Errorf("line %d: msgid cannot be empty", lineNum)
		}

		translations[msgid] = msgstr
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading translations: %w", err)
	}

	return translations, nil
}
