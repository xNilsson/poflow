# Fix 4: Deduplicate Output Functions

**Priority:** Medium
**Effort:** Medium
**Status:** ✅ COMPLETED

## Problem

Each command has its own `outputText()` and `outputJSON()` functions with similar logic, leading to code duplication.

**Locations:**
- `cmd/listempty.go`
- `cmd/search.go`
- `cmd/searchvalue.go`
- `cmd/translate.go`

**Issue:**
- Code duplication across 4 files
- Inconsistencies in output format
- Harder to maintain (need to fix bugs in multiple places)
- Violates DRY principle

## Solution

Create shared output functions in `internal/output/output.go`:

```go
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

// OutputEntryJSON outputs an entry in JSON format
func OutputEntryJSON(entry *model.MsgEntry) error {
    type jsonEntry struct {
        MsgID      string   `json:"msgid"`
        MsgStr     string   `json:"msgstr"`
        Comments   []string `json:"comments,omitempty"`
        References []string `json:"references,omitempty"`
    }

    j := jsonEntry{
        MsgID:      entry.MsgID,
        MsgStr:     entry.MsgStr,
        Comments:   entry.Comments,
        References: entry.References,
    }

    data, err := json.Marshal(j)
    if err != nil {
        return fmt.Errorf("failed to marshal JSON: %w", err)
    }

    fmt.Println(string(data))
    return nil
}
```

## Implementation Steps

1. ✅ Create `internal/output/output.go` with shared functions
2. ✅ Implement `OutputEntry()`, `OutputEntryText()`, `OutputEntryJSON()`
3. ✅ Update `cmd/listempty.go` to use shared functions
4. ✅ Update `cmd/search.go` to use shared functions
5. ✅ Update `cmd/searchvalue.go` to use shared functions
6. ✅ Update `cmd/translate.go` to use shared functions
7. ✅ Remove old duplicate `outputText()` and `outputJSON()` functions
8. ✅ Test all commands with both text and JSON output
9. ✅ Verify multi-line strings and comments work correctly

## Testing

Test each command with:
- Text output (`poflow <command> ...`)
- JSON output (`poflow <command> --json ...`)
- Multi-line strings
- Comments and references
- Empty translations

## Impact

**Before Fix:**
- ~80-100 lines of duplicated code across 4 files
- Inconsistent output handling
- Need to fix bugs in multiple places

**After Fix:**
- Single source of truth for output formatting
- Easier maintenance
- Consistent behavior across all commands
- Fixes issues #1, #2, and #3 automatically

## Notes

- This fix should be done **after** fixes #1, #2, and #3 (or incorporate them)
- Alternatively, do this fix **first** and it will solve #1, #2, #3 automatically
- The shared functions will ensure consistent output across all commands
- May need to adjust for command-specific needs (but most output is identical)

## Decision

**Approach:** Do this fix LAST (after fixes #1, #2, #3) to ensure we understand all the edge cases before refactoring. This way we can incorporate all the fixes into the shared functions.
