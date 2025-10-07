package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIssue1_CommentOrderPreserved verifies that comment order is preserved during edit
// Bug: https://github.com/xnilsson/poflow/issues/001
func TestIssue1_CommentOrderPreserved(t *testing.T) {
	// Create a temp directory for test files
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.po")

	// Original .po content with standard comment order:
	// 1. #: reference first
	// 2. #, flags second
	originalContent := `# HEADER COMMENT
msgid ""
msgstr ""

#: lib/monster_construction_web/components/pattern_components.ex:380
#, elixir-autogen, elixir-format
msgid "B: Fuller bust | C: Standard | D: Fuller hip"
msgstr "B: Fylligare byst | C: Standard | D: Fylligare höft"

`

	// Write test file
	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run the edit command
	oldMsgID := "B: Fuller bust | C: Standard | D: Fuller hip"
	newMsgID := "B: Athletic | C: Standard | D: Curvy"

	result, err := UpdateMsgIDInFile(testFile, oldMsgID, newMsgID, false)
	if err != nil {
		t.Fatalf("UpdateMsgIDInFile failed: %v", err)
	}

	if result.EntriesFound != 1 {
		t.Errorf("Expected 1 entry found, got %d", result.EntriesFound)
	}

	// Read the updated file
	updatedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	updatedStr := string(updatedContent)

	// Check 1: Reference comment should come BEFORE flags comment
	refIdx := strings.Index(updatedStr, "#: lib/monster_construction_web")
	flagsIdx := strings.Index(updatedStr, "#, elixir-autogen")

	if refIdx < 0 {
		t.Error("Reference comment (#:) not found in output")
	}
	if flagsIdx < 0 {
		t.Error("Flags comment (#,) not found in output")
	}
	if refIdx > flagsIdx {
		t.Errorf("Comment order changed! Reference (#:) at %d should come before flags (#,) at %d", refIdx, flagsIdx)
		t.Logf("Updated content:\n%s", updatedStr)
	}

	// Check 2: Flags comment should NOT have extra space (should be "#," not "# ,")
	if strings.Contains(updatedStr, "# , elixir-autogen") {
		t.Errorf("Flags comment has extra space: '# ,' instead of '#,'")
		t.Logf("Updated content:\n%s", updatedStr)
	}

	// Check 3: msgid should be updated
	if !strings.Contains(updatedStr, newMsgID) {
		t.Errorf("Updated msgid not found in output")
	}

	// Check 4: msgstr should be preserved
	if !strings.Contains(updatedStr, "B: Fylligare byst") {
		t.Errorf("Original msgstr was not preserved")
	}
}

// TestIssue1_MultipleCommentTypes tests various comment combinations
func TestIssue1_MultipleCommentTypes(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.po")

	originalContent := `# HEADER
msgid ""
msgstr ""

# Translator comment
#: lib/file.ex:100
#, elixir-autogen, elixir-format
# Another comment
msgid "Test"
msgstr "Testa"

`

	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := UpdateMsgIDInFile(testFile, "Test", "Updated", false)
	if err != nil {
		t.Fatalf("UpdateMsgIDInFile failed: %v", err)
	}

	if result.EntriesFound != 1 {
		t.Errorf("Expected 1 entry found, got %d", result.EntriesFound)
	}

	updatedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	updatedStr := string(updatedContent)

	// Verify exact order is preserved
	lines := strings.Split(updatedStr, "\n")

	// Find the entry (skip header)
	foundTranslator := false
	foundRef := false
	foundFlags := false
	foundAnother := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "# Translator comment" {
			foundTranslator = true
			t.Logf("Line %d: Translator comment found", i)
		}
		if strings.HasPrefix(trimmed, "#: lib/file.ex") {
			if !foundTranslator {
				t.Error("Reference comment appeared before translator comment")
			}
			foundRef = true
			t.Logf("Line %d: Reference found", i)
		}
		if strings.HasPrefix(trimmed, "#, elixir-autogen") {
			if !foundRef {
				t.Error("Flags comment appeared before reference comment")
			}
			foundFlags = true
			t.Logf("Line %d: Flags found", i)
		}
		if trimmed == "# Another comment" {
			if !foundFlags {
				t.Error("'Another comment' appeared before flags comment")
			}
			foundAnother = true
			t.Logf("Line %d: Another comment found", i)
		}
	}

	if !foundTranslator || !foundRef || !foundFlags || !foundAnother {
		t.Errorf("Not all comments found. Trans=%v Ref=%v Flags=%v Another=%v",
			foundTranslator, foundRef, foundFlags, foundAnother)
		t.Logf("Updated content:\n%s", updatedStr)
	}
}

// TestIssue2_SourceCodeUpdated verifies that source code files are updated when msgid changes
// Bug: https://github.com/xnilsson/poflow/issues/001
func TestIssue2_SourceCodeUpdated(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a source file
	sourceFile := filepath.Join(tmpDir, "lib", "pattern_components.ex")
	if err := os.MkdirAll(filepath.Dir(sourceFile), 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	sourceContent := `defmodule PatternComponents do
  def render(assigns) do
    ~H"""
    <div>{gettext("B: Fuller bust | C: Standard | D: Fuller hip")}</div>
    """
  end
end
`

	if err := os.WriteFile(sourceFile, []byte(sourceContent), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create a .po file with reference to the source file
	poFile := filepath.Join(tmpDir, "default.po")
	poContent := `# HEADER
msgid ""
msgstr ""

#: lib/pattern_components.ex:4
#, elixir-autogen, elixir-format
msgid "B: Fuller bust | C: Standard | D: Fuller hip"
msgstr "B: Fylligare byst | C: Standard | D: Fylligare höft"

`

	if err := os.WriteFile(poFile, []byte(poContent), 0644); err != nil {
		t.Fatalf("Failed to create .po file: %v", err)
	}

	// Run the edit command with source code updating enabled
	oldMsgID := "B: Fuller bust | C: Standard | D: Fuller hip"
	newMsgID := "B: Athletic | C: Standard | D: Curvy"

	result, err := UpdateMsgIDInFileWithSources(poFile, oldMsgID, newMsgID, false, tmpDir)
	if err != nil {
		t.Fatalf("UpdateMsgIDInFileWithSources failed: %v", err)
	}

	if result.EntriesFound != 1 {
		t.Errorf("Expected 1 entry found, got %d", result.EntriesFound)
	}

	// Read the updated source file
	updatedSource, err := os.ReadFile(sourceFile)
	if err != nil {
		t.Fatalf("Failed to read source file: %v", err)
	}

	updatedSourceStr := string(updatedSource)

	// Check that source code was updated
	if !strings.Contains(updatedSourceStr, newMsgID) {
		t.Errorf("Source file was not updated with new msgid")
		t.Logf("Expected to find: %s", newMsgID)
		t.Logf("Source content:\n%s", updatedSourceStr)
	}

	// Check that old msgid is gone
	if strings.Contains(updatedSourceStr, oldMsgID) {
		t.Errorf("Source file still contains old msgid")
		t.Logf("Old msgid should be removed: %s", oldMsgID)
		t.Logf("Source content:\n%s", updatedSourceStr)
	}

	// Verify the exact replacement
	expectedLine := `    <div>{gettext("B: Athletic | C: Standard | D: Curvy")}</div>`
	if !strings.Contains(updatedSourceStr, expectedLine) {
		t.Errorf("Source file does not contain expected line")
		t.Logf("Expected line: %s", expectedLine)
		t.Logf("Source content:\n%s", updatedSourceStr)
	}
}

// TestIssue2_MultipleSourceFiles tests updating multiple source files
func TestIssue2_MultipleSourceFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple source files
	sourceFile1 := filepath.Join(tmpDir, "lib", "file1.ex")
	sourceFile2 := filepath.Join(tmpDir, "lib", "file2.ex")
	if err := os.MkdirAll(filepath.Dir(sourceFile1), 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	source1Content := `defmodule File1 do
  def test do
    gettext("Welcome")
  end
end
`

	source2Content := `defmodule File2 do
  def test do
    dgettext("default", "Welcome")
  end
end
`

	if err := os.WriteFile(sourceFile1, []byte(source1Content), 0644); err != nil {
		t.Fatalf("Failed to create source file 1: %v", err)
	}
	if err := os.WriteFile(sourceFile2, []byte(source2Content), 0644); err != nil {
		t.Fatalf("Failed to create source file 2: %v", err)
	}

	// Create a .po file with references to both source files
	poFile := filepath.Join(tmpDir, "default.po")
	poContent := `# HEADER
msgid ""
msgstr ""

#: lib/file1.ex:3 lib/file2.ex:3
#, elixir-autogen, elixir-format
msgid "Welcome"
msgstr "Välkommen"

`

	if err := os.WriteFile(poFile, []byte(poContent), 0644); err != nil {
		t.Fatalf("Failed to create .po file: %v", err)
	}

	// Run the edit command
	result, err := UpdateMsgIDInFileWithSources(poFile, "Welcome", "Hello", false, tmpDir)
	if err != nil {
		t.Fatalf("UpdateMsgIDInFileWithSources failed: %v", err)
	}

	if result.EntriesFound != 1 {
		t.Errorf("Expected 1 entry found, got %d", result.EntriesFound)
	}

	// Check both source files were updated
	updatedSource1, _ := os.ReadFile(sourceFile1)
	updatedSource2, _ := os.ReadFile(sourceFile2)

	if !strings.Contains(string(updatedSource1), `gettext("Hello")`) {
		t.Errorf("Source file 1 was not updated correctly")
		t.Logf("Content:\n%s", string(updatedSource1))
	}

	if !strings.Contains(string(updatedSource2), `dgettext("default", "Hello")`) {
		t.Errorf("Source file 2 was not updated correctly")
		t.Logf("Content:\n%s", string(updatedSource2))
	}
}
