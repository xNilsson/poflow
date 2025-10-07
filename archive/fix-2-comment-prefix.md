# Fix 2: Comment Output Missing Prefix

**Priority:** Low
**Effort:** Trivial
**Status:** ✅ COMPLETED

## Problem

Comments in .po files are output without the proper `# ` prefix in text output mode.

**Locations:**
- `cmd/translate.go:184` (approximate)
- `cmd/search.go:127` (approximate)
- Possibly `cmd/searchvalue.go` and `cmd/listempty.go`

**Current Code:**
```go
for _, comment := range entry.Comments {
    fmt.Println(comment)  // Missing "# " prefix
}
```

**Issue:** Comments should be prefixed with `# ` to maintain valid .po file format.

## Solution

Add `# ` prefix when outputting comments:

```go
for _, comment := range entry.Comments {
    fmt.Printf("# %s\n", comment)
}
```

## Implementation Steps

1. ✅ Identify all files with comment output
2. ✅ Fix `cmd/search.go`
3. ✅ Fix `cmd/searchvalue.go`
4. ✅ Fix `cmd/listempty.go`
5. ✅ Fix `cmd/translate.go`
6. ✅ Test with .po file containing comments

## Testing

Create test .po file with comments:

```po
# This is a translator comment
#: reference.ex:42
msgid "Hello"
msgstr "Bonjour"
```

Verify comments output correctly with `# ` prefix.

## Impact

**Before Fix:** Comments output without prefix, breaking .po format
**After Fix:** Comments output with proper `# ` prefix

## Notes

- Very simple fix
- Only affects text output (JSON already includes comments in structured format)
- Improves .po format compliance
