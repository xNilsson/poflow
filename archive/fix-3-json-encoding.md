# Fix 3: Use encoding/json for JSON Output

**Priority:** Low
**Effort:** Low
**Status:** ✅ COMPLETED

## Problem

The `outputEntryJSON()` function in `cmd/translate.go:169-181` manually constructs JSON strings instead of using Go's `encoding/json` package.

**Location:** `cmd/translate.go:169-181`

**Current Code:**
```go
fmt.Printf(`{"msgid":%q,"msgstr":%q%s%s}`+"\n", ...)  // Manual JSON construction
```

**Risk:** Potential escaping bugs with complex strings containing quotes, backslashes, or special characters in comments/references.

## Solution

Use `json.Marshal()` like in `cmd/search.go:121`:

```go
import "encoding/json"

type jsonEntry struct {
    MsgID      string   `json:"msgid"`
    MsgStr     string   `json:"msgstr"`
    Comments   []string `json:"comments,omitempty"`
    References []string `json:"references,omitempty"`
}

func outputEntryJSON(entry *model.MsgEntry) error {
    j := jsonEntry{
        MsgID:      entry.MsgID,
        MsgStr:     entry.MsgStr,
        Comments:   entry.Comments,
        References: entry.References,
    }
    data, err := json.Marshal(j)
    if err != nil {
        return err
    }
    fmt.Println(string(data))
    return nil
}
```

## Implementation Steps

1. ✅ Read `cmd/search.go` to see how it uses `json.Marshal()`
2. ✅ Read `cmd/translate.go` to see current manual JSON construction
3. ✅ Replace manual JSON with `json.Marshal()` in `cmd/translate.go`
4. ✅ Check if other commands also manually construct JSON
5. ✅ Test with complex strings (quotes, backslashes, unicode)

## Testing

Test with .po entries containing special characters:

```po
msgid "String with \"quotes\" and \\ backslashes"
msgstr "Chaîne with 'quotes' and unicode: 你好"

# Comment with "quotes" and special chars: $@!
msgid "Another"
msgstr "Autre"
```

Verify JSON output is valid and properly escaped.

## Impact

**Before Fix:** Potential escaping bugs with complex strings
**After Fix:** Robust JSON output using standard library

## Notes

- Most commands already use `json.Marshal()` correctly
- Only `translate.go` appears to use manual JSON construction
- This makes the code more maintainable and safer
