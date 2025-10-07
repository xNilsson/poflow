# Phase 1: Core Parsing + listempty

**Status:** ✅ COMPLETED
**Goal:** Implement the .po file parser and the first command (listempty).

## Tasks

- [x] Create MsgEntry struct in internal/model/entry.go
- [x] Implement streaming .po parser in internal/parser/po.go
- [x] Handle multi-line strings
- [x] Handle escaped quotes
- [x] Handle comments and references
- [x] Implement listempty command
- [x] Add --json output support
- [x] Add --limit flag
- [x] Write tests for parser

## Key Components

### MsgEntry Structure

```go
type MsgEntry struct {
    MsgID      string   `json:"msgid"`
    MsgStr     string   `json:"msgstr"`
    Comments   []string `json:"comments,omitempty"`
    References []string `json:"references,omitempty"`
}
```

### Parser Behavior

- Stream files line by line (no full in-memory AST)
- Detect new entry at `msgid`
- Accumulate until blank line
- Yield entries one by one

## Test Cases

1. Simple msgid/msgstr pair
2. Multi-line strings
3. Escaped quotes
4. Comments and references
5. Empty msgstr
6. File with mixed translated/untranslated entries

## Deliverables

- Working parser that can read .po files
- `poflow listempty file.po` command
- Both text and JSON output modes
