package parser

import (
	"strings"
	"testing"
)

// TestSearchScenarios tests realistic search scenarios
func TestSearchScenarios(t *testing.T) {
	input := `# Test file
msgid "Welcome"
msgstr "Välkommen"

msgid "Sign In"
msgstr ""

msgid "Sign Out"
msgstr "Logga ut"

msgid "Profile"
msgstr ""

msgid "Settings"
msgstr "Inställningar"
`

	tests := []struct {
		name          string
		searchMsgID   string
		expectCount   int
		expectFirst   string
		caseVariation string
	}{
		{
			name:        "substring match",
			searchMsgID: "sign",
			expectCount: 2,
			expectFirst: "Sign In",
		},
		{
			name:        "exact match",
			searchMsgID: "welcome",
			expectCount: 1,
			expectFirst: "Welcome",
		},
		{
			name:          "case insensitive",
			searchMsgID:   "SETTINGS",
			expectCount:   1,
			expectFirst:   "Settings",
			caseVariation: "uppercase",
		},
		{
			name:        "no match",
			searchMsgID: "nonexistent",
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(input)
			p := NewParser(reader)

			searchPattern := strings.ToLower(tt.searchMsgID)
			count := 0
			var firstMatch string

			for {
				entry := p.Next()
				if entry == nil {
					break
				}

				if strings.Contains(strings.ToLower(entry.MsgID), searchPattern) {
					count++
					if firstMatch == "" {
						firstMatch = entry.MsgID
					}
				}
			}

			if err := p.Err(); err != nil {
				t.Fatalf("parser error: %v", err)
			}

			if count != tt.expectCount {
				t.Errorf("expected %d matches, got %d", tt.expectCount, count)
			}

			if tt.expectCount > 0 && firstMatch != tt.expectFirst {
				t.Errorf("expected first match %q, got %q", tt.expectFirst, firstMatch)
			}
		})
	}
}

// TestSearchValueScenarios tests searching by msgstr
func TestSearchValueScenarios(t *testing.T) {
	input := `msgid "Welcome"
msgstr "Välkommen"

msgid "Sign In"
msgstr ""

msgid "Sign Out"
msgstr "Logga ut"

msgid "Settings"
msgstr "Inställningar"
`

	tests := []struct {
		name        string
		searchMsgStr string
		expectCount int
		expectMsgID string
	}{
		{
			name:        "search translated value",
			searchMsgStr: "logga",
			expectCount: 1,
			expectMsgID: "Sign Out",
		},
		{
			name:        "search empty translation",
			searchMsgStr: "",
			expectCount: 4, // All entries have msgstr (even if empty)
			expectMsgID: "Welcome",
		},
		{
			name:        "no match in translations",
			searchMsgStr: "xyz",
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(input)
			p := NewParser(reader)

			searchPattern := strings.ToLower(tt.searchMsgStr)
			count := 0
			var firstMsgID string

			for {
				entry := p.Next()
				if entry == nil {
					break
				}

				if strings.Contains(strings.ToLower(entry.MsgStr), searchPattern) {
					count++
					if firstMsgID == "" {
						firstMsgID = entry.MsgID
					}
				}
			}

			if err := p.Err(); err != nil {
				t.Fatalf("parser error: %v", err)
			}

			if count != tt.expectCount {
				t.Errorf("expected %d matches, got %d", tt.expectCount, count)
			}

			if tt.expectCount > 0 && tt.expectMsgID != "" && firstMsgID != tt.expectMsgID {
				t.Errorf("expected first msgid %q, got %q", tt.expectMsgID, firstMsgID)
			}
		})
	}
}

// TestSearchLimit tests limiting search results
func TestSearchLimit(t *testing.T) {
	input := `msgid "Sign In"
msgstr ""

msgid "Sign Out"
msgstr "Logga ut"

msgid "Sign Up"
msgstr ""
`

	reader := strings.NewReader(input)
	p := NewParser(reader)

	searchPattern := "sign"
	limit := 2
	count := 0

	for {
		entry := p.Next()
		if entry == nil {
			break
		}

		if strings.Contains(strings.ToLower(entry.MsgID), searchPattern) {
			count++
			if count >= limit {
				break
			}
		}
	}

	if count != limit {
		t.Errorf("expected exactly %d results with limit, got %d", limit, count)
	}
}
