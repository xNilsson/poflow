# Plan: Edit Text Command

**Status:** Draft
**Created:** 2025-10-07
**Priority:** Medium

## Problem Statement

When working with gettext `.po` files, the `msgid` field typically contains the English source text. When you want to **change** the source text (not just translate it), you currently need to:

1. Update the `msgid` in the Swedish (sv) `.po` file
2. Update the `msgid` in the English (en) `.po` file
3. Update the `msgid` in the `.pot` template file
4. Update all source code references that use that text

This is error-prone and tedious. A single text change requires updating multiple files and code locations.

### Example Scenario

You have:
```po
# sv/LC_MESSAGES/default.po
msgid "Sign In"
msgstr "Logga in"

# en/LC_MESSAGES/default.po
msgid "Sign In"
msgstr "Sign In"

# default.pot
msgid "Sign In"
msgstr ""
```

And in your code:
```elixir
gettext("Sign In")
```

You want to change "Sign In" to "Log In" everywhere. Currently you need to:
1. Change `msgid` in sv file to "Log In"
2. Change `msgid` in en file to "Log In"
3. Change `msgid` in .pot template to "Log In"
4. Update code: `gettext("Log In")`
5. Optionally update `msgstr` in en file to "Log In"

## Proposed Solution

Add a `poflow edit` command that automates updating `msgid` across all language files, the template file, and provides a report of code locations that need updating.

### Proposed Usage

```bash
# Edit a msgid across all files
poflow edit "Sign In" "Log In"

# Output:
# Updated msgid in 3 file(s):
#   ✓ priv/gettext/sv/LC_MESSAGES/default.po
#   ✓ priv/gettext/en/LC_MESSAGES/default.po
#   ✓ priv/gettext/default.pot
#
# Code references that need updating (3):
#   lib/auth/login.ex:42
#   lib/components/header.ex:18
#   test/auth_test.exs:67
#
# To update code references, run:
#   grep -r "Sign In" lib/ test/
```

### Behavior

1. **Find all `.po` and `.pot` files** in the configured gettext directory
2. **Search for the old msgid** in all files
3. **Update msgid** in all matching entries (preserving msgstr values)
4. **Update code references** (optional, or just report them)
5. **Show summary** of what was changed

### Flags

```bash
--dry-run        # Show what would be changed without modifying files
--all-languages  # Update all language files (default behavior)
--code-path      # Path to search for code references (default: lib/)
--update-code    # Actually update code references (default: false, just report)
```

## Technical Design

### Command Flow

1. Load config to find `gettext_path`
2. Find all `.po` and `.pot` files under `gettext_path`
3. For each file:
   - Parse with existing `parser.Parser`
   - Find entries where `entry.MsgID == oldText`
   - Update `entry.MsgID = newText`
   - Preserve `entry.MsgStr` (don't change translations)
   - Write back to file using `output.FormatEntry()`
4. Search code files for old text references
5. Report what was changed and what needs manual updates

### File Update Strategy

**Option 1: In-place update (like current `translate` command)**
- Read .po file
- Parse all entries
- Update matching msgid entries
- Write to temp file
- Replace original file

**Option 2: Smart line-by-line replacement**
- Read .po file line by line
- Detect `msgid "old text"` lines
- Replace with `msgid "new text"`
- Preserve all other lines exactly
- More efficient, preserves file structure better

**Recommendation:** Use Option 1 (in-place update with parser) because:
- Reuses existing tested parser code
- Handles multi-line msgid correctly
- Handles escape sequences properly
- Consistent with `translate` command behavior

### Code Structure

```
cmd/
  edit.go              # New command implementation

internal/
  parser/
    po.go              # Existing parser (no changes needed)
  output/
    output.go          # Existing output formatter (no changes needed)
  editor/
    edit.go            # New package for edit operations
    code_search.go     # Search code files for references
```

### Implementation Steps

1. **Phase 1: Basic msgid replacement**
   - Create `cmd/edit.go` with basic structure
   - Find all .po/.pot files in gettext_path
   - Parse each file and update matching msgid entries
   - Write back to files
   - Add `--dry-run` flag

2. **Phase 2: Code reference detection**
   - Search code files for old text
   - Report file paths and line numbers
   - Show helpful grep command

3. **Phase 3: Code reference updating** (optional)
   - Add `--update-code` flag
   - Parse code files (basic text replacement)
   - Update references automatically
   - Handle escape sequences in different languages

## Edge Cases to Handle

### Multi-line msgid
```po
msgid ""
"This is a very long "
"text that spans multiple lines"
msgstr "Translation here"
```

**Solution:** Parser already handles this. When we update `entry.MsgID`, the `output.FormatEntry()` will re-serialize it correctly.

### Escaped characters
```po
msgid "Click \"OK\" to continue"
msgstr "Klicka \"OK\" för att fortsätta"
```

**Solution:** Parser already handles unescaping. When we write back, `output.escapeString()` handles re-escaping.

### Partial matches
```po
msgid "Sign In"
msgstr "Logga in"

msgid "Sign In Here"
msgstr "Logga in här"
```

**Problem:** If we search for "Sign In", should we match both?

**Solution:**
- Default behavior: **Exact match only** (`entry.MsgID == oldText`)
- Add `--partial` flag for substring matching
- Add `--regex` flag for advanced matching

### Multiple msgid matches in same file
```po
# Entry 1
msgid "Submit"
msgstr "Skicka"

# Entry 2
msgid "Submit"
msgstr "Skicka in"
```

**Problem:** Duplicate msgid (shouldn't happen in valid .po files, but might)

**Solution:** Update all matches, show count in summary

### File permissions
```bash
# What if .po file is read-only?
poflow edit "old" "new"
# Error: failed to write priv/gettext/sv/LC_MESSAGES/default.po: permission denied
```

**Solution:** Show clear error message, suggest checking file permissions

## Alternative Approaches

### Alternative 1: Shell script
```bash
#!/bin/bash
OLD="$1"
NEW="$2"
find priv/gettext -name "*.po" -o -name "*.pot" | \
  xargs sed -i "s/msgid \"$OLD\"/msgid \"$NEW\"/g"
```

**Pros:** Simple, no code changes needed
**Cons:**
- Doesn't handle multi-line msgid
- Doesn't handle escape sequences correctly
- Doesn't validate .po file format
- No dry-run or safety checks
- Doesn't help with code references

### Alternative 2: Use existing gettext tools
```bash
# Extract new .pot
mix gettext.extract

# Merge into .po files
mix gettext.merge
```

**Pros:** Uses official tools, well-tested
**Cons:**
- Requires changing code first
- Doesn't help find all code references
- Separate step for each language
- Doesn't update old msgid across files

### Alternative 3: Manual find/replace in editor
**Pros:** Full control, visual feedback
**Cons:**
- Error-prone (easy to miss files)
- No validation
- Tedious for multiple files
- Hard to track what was changed

**Recommendation:** Implement `poflow edit` command because:
- Safer than shell scripts (validates .po format)
- Faster than manual editing
- Helps find code references
- Consistent with poflow's philosophy
- Reuses existing parser/output code

## Success Criteria

A successful implementation should:

1. ✅ Find and update msgid in all .po and .pot files
2. ✅ Preserve msgstr (translations) exactly
3. ✅ Handle multi-line msgid correctly
4. ✅ Handle escape sequences correctly
5. ✅ Show clear summary of changes
6. ✅ Support `--dry-run` for safety
7. ✅ Detect code references that need updating
8. ✅ Work with existing config system
9. ✅ Fail gracefully with helpful errors
10. ✅ Be well-documented with examples

## Testing Strategy

### Unit Tests
```go
// internal/editor/edit_test.go
TestFindFiles         // Find all .po/.pot files
TestUpdateMsgID       // Update msgid in entry
TestPreserveTranslation // Ensure msgstr unchanged
TestMultilineMsgID    // Handle multi-line msgid
TestEscapedCharacters // Handle escape sequences
```

### Integration Tests
```bash
# test_edit.sh
# Create test .po files
# Run poflow edit
# Verify changes
# Test --dry-run
# Test code reference detection
```

### Manual Testing
```bash
# Basic edit
./poflow edit "Sign In" "Log In"

# Dry run
./poflow edit --dry-run "Sign In" "Log In"

# With config
./poflow --config test.yml edit "old" "new"
```

## Documentation Updates

### README.md
Add new section under "Commands":

```markdown
### `edit` - Update msgid Across Files

Update a msgid (source text) across all language files and templates.

Usage:
  poflow edit "old text" "new text"
  poflow edit --dry-run "old text" "new text"

This updates:
- All .po language files
- .pot template file
- Reports code references that need updating
```

### TUTORIAL.md
Add practical example:

```markdown
## Changing Source Text

When you need to change the English source text, use `poflow edit`:

$ poflow edit "Sign In" "Log In"
Updated msgid in 3 file(s):
  ✓ priv/gettext/sv/LC_MESSAGES/default.po
  ✓ priv/gettext/en/LC_MESSAGES/default.po
  ✓ priv/gettext/default.pot

Code references (2):
  lib/auth.ex:42
  lib/header.ex:18
```

## Dependencies

### Existing Code (Reuse)
- `internal/parser/po.go` - Parse .po files ✅
- `internal/output/output.go` - Format .po entries ✅
- `internal/config/config.go` - Load config ✅

### New Code (Implement)
- `cmd/edit.go` - Command implementation
- `internal/editor/edit.go` - Core edit logic
- `internal/editor/files.go` - File discovery
- `internal/editor/code_search.go` - Code reference search

### External Dependencies
None needed (use standard library).

## Risks and Mitigations

### Risk 1: Data loss if update fails mid-operation
**Mitigation:**
- Write to temp files first (like `translate` command)
- Only replace original files if all updates succeed
- Provide `--dry-run` to preview changes

### Risk 2: Breaking code references
**Mitigation:**
- Detect code references and report them clearly
- Don't update code by default (require `--update-code` flag)
- Provide example grep command for manual verification

### Risk 3: Corrupting .po file format
**Mitigation:**
- Use existing tested parser/output code
- Validate files can be parsed after update
- Preserve exact file structure (comments, references, etc.)

### Risk 4: Partial matches updating wrong entries
**Mitigation:**
- Default to exact match only
- Show count of entries that will be updated
- Require `--partial` flag for substring matching

## Timeline Estimate

- **Phase 1** (Basic msgid replacement): 2-3 hours
- **Phase 2** (Code reference detection): 1-2 hours
- **Phase 3** (Code updating, optional): 2-3 hours
- **Testing**: 1-2 hours
- **Documentation**: 1 hour

**Total:** 7-11 hours for full implementation

**Minimum Viable:** Phase 1 + basic testing = 3-4 hours

## Future Enhancements

1. **Interactive mode:** Prompt for confirmation before each change
2. **Regex support:** `poflow edit --regex "Sign.*" "Log.*"`
3. **Batch edits:** `poflow edit --file changes.txt` with multiple old→new pairs
4. **Undo support:** `poflow edit --undo` to revert last edit
5. **Git integration:** Auto-commit changes with descriptive message
6. **Multi-language code search:** Detect references in Elixir, JS, etc.

## Related Issues

This addresses the common workflow problem where:
- Developers want to change UI text
- Text is duplicated across multiple .po files
- Code references need to be found and updated
- Current tools don't automate this process

## References

- [GNU gettext manual](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html)
- [Existing poflow translate command](cmd/translate.go) - Similar file update pattern
- [Existing poflow search command](cmd/search.go) - Similar search pattern

## Notes

This command fills a gap in the current poflow feature set:
- `listempty` - Find untranslated entries ✅
- `search` - Find entries by msgid ✅
- `translate` - Update msgstr (translations) ✅
- `edit` - Update msgid (source text) ⭐ NEW

The edit command completes the workflow by handling source text changes, which is a common operation when iterating on UI text.
