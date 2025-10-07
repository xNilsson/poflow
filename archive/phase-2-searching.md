# Phase 2: Searching

**Status:** COMPLETED
**Goal:** Add search and searchvalue commands with regex/substring support.

## Tasks

- [x] Implement `search` command (search in msgid)
- [x] Implement `searchvalue` command (search in msgstr)
- [x] Add --re flag for regex matching
- [x] Add --plain flag for substring matching (default)
- [x] Add --json output support
- [x] Add --limit flag
- [x] Handle case-insensitive matching
- [x] Write tests for search functionality

## Commands

### search

Search for entries by msgid pattern:

```bash
poflow search --re "Pattern" file.po
poflow search --plain "Welcome" file.po
poflow search --json "Login" file.po
```

### searchvalue

Search for entries by msgstr pattern:

```bash
poflow searchvalue --re "Tack" file.po
poflow searchvalue --plain "Välkommen" file.po
```

## Output Examples

Text mode:
```
#: lib/web/live/page.ex:24
msgid "Welcome"
msgstr "Välkommen"
```

JSON mode:
```json
{"msgid":"Welcome","msgstr":"Välkommen","references":["lib/web/live/page.ex:24"]}
```

## Test Cases

1. Regex pattern matching
2. Plain substring matching
3. Case sensitivity
4. No matches found
5. Multiple matches
6. JSON output format

## Deliverables

- `poflow search` command
- `poflow searchvalue` command
- Both regex and plain text search modes
- JSON output support
