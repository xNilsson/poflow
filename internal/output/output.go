package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nille/poflow/internal/model"
)

// OutputEntry outputs a single entry in text or JSON format
func OutputEntry(entry *model.MsgEntry, jsonFormat bool) error {
	if jsonFormat {
		return OutputEntryJSON(entry)
	}
	return OutputEntryText(entry)
}

// OutputEntryJSON outputs an entry in JSON format (line-delimited)
func OutputEntryJSON(entry *model.MsgEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// FormatEntry returns an entry formatted as .po text
func FormatEntry(entry *model.MsgEntry) string {
	var sb strings.Builder

	// Output comments with # prefix
	for _, comment := range entry.Comments {
		sb.WriteString(fmt.Sprintf("# %s\n", comment))
	}

	// Output references with #: prefix
	for _, ref := range entry.References {
		sb.WriteString(fmt.Sprintf("#: %s\n", ref))
	}

	// Output msgid (handle multi-line)
	if strings.Contains(entry.MsgID, "\n") {
		sb.WriteString("msgid \"\"\n")
		for _, line := range strings.Split(entry.MsgID, "\n") {
			sb.WriteString(fmt.Sprintf("\"%s\\n\"\n", line))
		}
	} else {
		sb.WriteString(fmt.Sprintf("msgid \"%s\"\n", entry.MsgID))
	}

	// Output msgstr (handle multi-line)
	if strings.Contains(entry.MsgStr, "\n") {
		sb.WriteString("msgstr \"\"\n")
		for _, line := range strings.Split(entry.MsgStr, "\n") {
			sb.WriteString(fmt.Sprintf("\"%s\\n\"\n", line))
		}
	} else {
		sb.WriteString(fmt.Sprintf("msgstr \"%s\"\n", entry.MsgStr))
	}

	sb.WriteString("\n") // Blank line between entries
	return sb.String()
}

// OutputEntryText outputs an entry in .po text format
func OutputEntryText(entry *model.MsgEntry) error {
	fmt.Print(FormatEntry(entry))
	return nil
}
