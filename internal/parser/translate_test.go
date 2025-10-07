package parser

import (
	"strings"
	"testing"
)

func TestParseTranslations_Simple(t *testing.T) {
	input := `Sign In = Logga in
Sign Out = Logga ut
Welcome = Välkommen`

	translations, err := ParseTranslations(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{
		"Sign In":  "Logga in",
		"Sign Out": "Logga ut",
		"Welcome":  "Välkommen",
	}

	if len(translations) != len(expected) {
		t.Fatalf("expected %d translations, got %d", len(expected), len(translations))
	}

	for msgid, expectedMsgstr := range expected {
		if msgstr, ok := translations[msgid]; !ok {
			t.Errorf("missing translation for %q", msgid)
		} else if msgstr != expectedMsgstr {
			t.Errorf("for %q: expected %q, got %q", msgid, expectedMsgstr, msgstr)
		}
	}
}

func TestParseTranslations_WithComments(t *testing.T) {
	input := `# This is a comment
Sign In = Logga in
# Another comment
Sign Out = Logga ut`

	translations, err := ParseTranslations(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(translations) != 2 {
		t.Fatalf("expected 2 translations, got %d", len(translations))
	}

	if translations["Sign In"] != "Logga in" {
		t.Errorf("unexpected translation for Sign In: %q", translations["Sign In"])
	}
}

func TestParseTranslations_WithEmptyLines(t *testing.T) {
	input := `Sign In = Logga in

Sign Out = Logga ut

Welcome = Välkommen`

	translations, err := ParseTranslations(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(translations) != 3 {
		t.Fatalf("expected 3 translations, got %d", len(translations))
	}
}

func TestParseTranslations_WithWhitespace(t *testing.T) {
	input := `  Sign In  =  Logga in
Sign Out=Logga ut
  Welcome  =Välkommen`

	translations, err := ParseTranslations(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{
		"Sign In":  "Logga in",
		"Sign Out": "Logga ut",
		"Welcome":  "Välkommen",
	}

	for msgid, expectedMsgstr := range expected {
		if msgstr, ok := translations[msgid]; !ok {
			t.Errorf("missing translation for %q", msgid)
		} else if msgstr != expectedMsgstr {
			t.Errorf("for %q: expected %q, got %q", msgid, expectedMsgstr, msgstr)
		}
	}
}

func TestParseTranslations_InvalidFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "no equals sign",
			input: "Sign In",
		},
		{
			name:  "empty msgid",
			input: " = Logga in",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTranslations(strings.NewReader(tt.input))
			if err == nil {
				t.Errorf("expected error for input %q, got nil", tt.input)
			}
		})
	}
}

func TestParseTranslations_MsgstrWithEquals(t *testing.T) {
	// msgstr can contain equals signs
	input := `Equation = x = y + 1`

	translations, err := ParseTranslations(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(translations) != 1 {
		t.Fatalf("expected 1 translation, got %d", len(translations))
	}

	if translations["Equation"] != "x = y + 1" {
		t.Errorf("expected %q, got %q", "x = y + 1", translations["Equation"])
	}
}

func TestParseTranslations_EmptyMsgstr(t *testing.T) {
	input := `Sign In =
Sign Out = Logga ut`

	translations, err := ParseTranslations(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(translations) != 2 {
		t.Fatalf("expected 2 translations, got %d", len(translations))
	}

	if translations["Sign In"] != "" {
		t.Errorf("expected empty msgstr for Sign In, got %q", translations["Sign In"])
	}

	if translations["Sign Out"] != "Logga ut" {
		t.Errorf("unexpected translation for Sign Out: %q", translations["Sign Out"])
	}
}
