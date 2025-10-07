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

// OutputEntryText outputs an entry in .po text format
func OutputEntryText(entry *model.MsgEntry) error {
	// Output comments with # prefix
	for _, comment := range entry.Comments {
		fmt.Printf("# %s\n", comment)
	}

	// Output references with #: prefix
	for _, ref := range entry.References {
		fmt.Printf("#: %s\n", ref)
	}

	// Output msgid (handle multi-line)
	if strings.Contains(entry.MsgID, "\n") {
		fmt.Println("msgid \"\"")
		for _, line := range strings.Split(entry.MsgID, "\n") {
			fmt.Printf("\"%s\\n\"\n", line)
		}
	} else {
		fmt.Printf("msgid \"%s\"\n", entry.MsgID)
	}

	// Output msgstr (handle multi-line)
	if strings.Contains(entry.MsgStr, "\n") {
		fmt.Println("msgstr \"\"")
		for _, line := range strings.Split(entry.MsgStr, "\n") {
			fmt.Printf("\"%s\\n\"\n", line)
		}
	} else {
		fmt.Printf("msgstr \"%s\"\n", entry.MsgStr)
	}

	fmt.Println() // Blank line between entries
	return nil
}
