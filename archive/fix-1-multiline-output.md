# Fix 1: Multi-line String Output Bug

**Priority:** Medium
**Effort:** Low
**Status:** ✅ COMPLETED

## Problem

The `outputText()` functions in `cmd/search.go` and `cmd/listempty.go` don't handle multi-line strings correctly. They output msgid/msgstr on a single line with literal `\n` in the string, which breaks the .po file format.

**Locations:**
- `cmd/search.go:125` (approximate)
- `cmd/listempty.go:outputText()`

**Current Code:**
```go
fmt.Printf("msgid \"%s\"\n", entry.MsgID)  // Breaks on newlines
```

**Issue:** If `entry.MsgID` contains `\n`, this outputs invalid .po format.

## Solution

Copy the multi-line handling logic from `cmd/translate.go:183-210`, which correctly outputs:

```go
if strings.Contains(entry.MsgID, "\n") {
    fmt.Println("msgid \"\"")
    for _, line := range strings.Split(entry.MsgID, "\n") {
        fmt.Printf("\"%s\\n\"\n", line)
    }
} else {
    fmt.Printf("msgid \"%s\"\n", entry.MsgID)
}
```

Same logic applies to `msgstr`.

## Implementation Steps

1. ✅ Read `cmd/translate.go` to understand the correct multi-line output logic
2. ✅ Read `cmd/search.go` to see current implementation
3. ✅ Fix `outputText()` in `cmd/search.go` for both msgid and msgstr
4. ✅ Read `cmd/listempty.go` to see current implementation
5. ✅ Fix `outputText()` in `cmd/listempty.go` for both msgid and msgstr
6. ✅ Test with a .po file containing multi-line strings

## Testing

Create a test .po file with multi-line msgid:

```po
msgid ""
"First line\n"
"Second line\n"
"Third line"
msgstr ""
"Première ligne\n"
"Deuxième ligne\n"
"Troisième ligne"
```

Verify output from `listempty` and `search` commands matches proper .po format.

## Impact

**Before Fix:** Multi-line strings output incorrectly, breaking .po file format
**After Fix:** Multi-line strings output correctly in proper .po format

## Notes

- The `translate.go` implementation already handles this correctly
- JSON output is unaffected (already works correctly)
- Only affects text (non-JSON) output mode
