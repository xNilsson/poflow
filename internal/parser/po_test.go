package parser

import (
	"strings"
	"testing"
)

func TestParser_SimplePair(t *testing.T) {
	input := `msgid "Hello"
msgstr "Hej"

`
	parser := NewParser(strings.NewReader(input))
	entry := parser.Next()

	if entry == nil {
		t.Fatal("expected entry, got nil")
	}

	if entry.MsgID != "Hello" {
		t.Errorf("expected msgid 'Hello', got '%s'", entry.MsgID)
	}

	if entry.MsgStr != "Hej" {
		t.Errorf("expected msgstr 'Hej', got '%s'", entry.MsgStr)
	}

	// Should be no more entries
	if parser.Next() != nil {
		t.Error("expected no more entries")
	}
}

func TestParser_EmptyTranslation(t *testing.T) {
	input := `msgid "Untranslated"
msgstr ""

`
	parser := NewParser(strings.NewReader(input))
	entry := parser.Next()

	if entry == nil {
		t.Fatal("expected entry, got nil")
	}

	if entry.MsgID != "Untranslated" {
		t.Errorf("expected msgid 'Untranslated', got '%s'", entry.MsgID)
	}

	if entry.MsgStr != "" {
		t.Errorf("expected empty msgstr, got '%s'", entry.MsgStr)
	}

	if !entry.IsEmpty() {
		t.Error("expected entry.IsEmpty() to be true")
	}
}

func TestParser_MultilineString(t *testing.T) {
	input := `msgid ""
"This is a long "
"multiline string"
msgstr ""
"Detta är en lång "
"flerradig sträng"

`
	parser := NewParser(strings.NewReader(input))
	entry := parser.Next()

	if entry == nil {
		t.Fatal("expected entry, got nil")
	}

	expectedMsgID := "This is a long multiline string"
	if entry.MsgID != expectedMsgID {
		t.Errorf("expected msgid '%s', got '%s'", expectedMsgID, entry.MsgID)
	}

	expectedMsgStr := "Detta är en lång flerradig sträng"
	if entry.MsgStr != expectedMsgStr {
		t.Errorf("expected msgstr '%s', got '%s'", expectedMsgStr, entry.MsgStr)
	}
}

func TestParser_WithComments(t *testing.T) {
	input := `# This is a comment
#: lib/web/live/page.ex:24
msgid "Welcome"
msgstr "Välkommen"

`
	parser := NewParser(strings.NewReader(input))
	entry := parser.Next()

	if entry == nil {
		t.Fatal("expected entry, got nil")
	}

	if len(entry.Comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(entry.Comments))
	}

	if entry.Comments[0] != "This is a comment" {
		t.Errorf("expected comment 'This is a comment', got '%s'", entry.Comments[0])
	}

	if len(entry.References) != 1 {
		t.Errorf("expected 1 reference, got %d", len(entry.References))
	}

	if entry.References[0] != "lib/web/live/page.ex:24" {
		t.Errorf("expected reference 'lib/web/live/page.ex:24', got '%s'", entry.References[0])
	}
}

func TestParser_EscapedQuotes(t *testing.T) {
	input := `msgid "Say \"Hello\""
msgstr "Säg \"Hej\""

`
	parser := NewParser(strings.NewReader(input))
	entry := parser.Next()

	if entry == nil {
		t.Fatal("expected entry, got nil")
	}

	expectedMsgID := `Say "Hello"`
	if entry.MsgID != expectedMsgID {
		t.Errorf("expected msgid '%s', got '%s'", expectedMsgID, entry.MsgID)
	}

	expectedMsgStr := `Säg "Hej"`
	if entry.MsgStr != expectedMsgStr {
		t.Errorf("expected msgstr '%s', got '%s'", expectedMsgStr, entry.MsgStr)
	}
}

func TestParser_MultipleEntries(t *testing.T) {
	input := `msgid "First"
msgstr "Första"

msgid "Second"
msgstr ""

msgid "Third"
msgstr "Tredje"
`
	parser := NewParser(strings.NewReader(input))

	// First entry
	entry := parser.Next()
	if entry == nil || entry.MsgID != "First" {
		t.Errorf("expected first entry with msgid 'First'")
	}

	// Second entry
	entry = parser.Next()
	if entry == nil || entry.MsgID != "Second" || entry.MsgStr != "" {
		t.Errorf("expected second entry with msgid 'Second' and empty msgstr")
	}

	// Third entry
	entry = parser.Next()
	if entry == nil || entry.MsgID != "Third" {
		t.Errorf("expected third entry with msgid 'Third'")
	}

	// No more entries
	if parser.Next() != nil {
		t.Error("expected no more entries")
	}
}

func TestParser_NewlineEscapes(t *testing.T) {
	input := `msgid "Line one\nLine two"
msgstr "Rad ett\nRad två"

`
	parser := NewParser(strings.NewReader(input))
	entry := parser.Next()

	if entry == nil {
		t.Fatal("expected entry, got nil")
	}

	expectedMsgID := "Line one\nLine two"
	if entry.MsgID != expectedMsgID {
		t.Errorf("expected msgid with newline, got '%s'", entry.MsgID)
	}

	expectedMsgStr := "Rad ett\nRad två"
	if entry.MsgStr != expectedMsgStr {
		t.Errorf("expected msgstr with newline, got '%s'", entry.MsgStr)
	}
}

func TestParseAll(t *testing.T) {
	input := `msgid "One"
msgstr "Ett"

msgid "Two"
msgstr ""

msgid "Three"
msgstr "Tre"
`
	entries, err := ParseAll(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}

	if entries[0].MsgID != "One" {
		t.Errorf("expected first msgid 'One', got '%s'", entries[0].MsgID)
	}

	if entries[1].MsgStr != "" {
		t.Errorf("expected second msgstr to be empty, got '%s'", entries[1].MsgStr)
	}

	if entries[2].MsgID != "Three" {
		t.Errorf("expected third msgid 'Three', got '%s'", entries[2].MsgID)
	}
}
