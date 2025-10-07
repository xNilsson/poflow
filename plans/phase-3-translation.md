# Phase 3: Translation

**Status:** Not Started
**Goal:** Implement the translate command to merge translations into .po files.

## Tasks

- [ ] Parse translation input format (`msgid = msgstr`)
- [ ] Implement translate command with stdin support
- [ ] Implement translate command with file input
- [ ] Add --language flag for config-based path resolution
- [ ] Merge translations into existing .po file
- [ ] Preserve comments and references
- [ ] Handle missing msgids gracefully
- [ ] Support both direct file path and config-based resolution
- [ ] Write tests for translation merging

## Usage Patterns

### Config-based (preferred)

```bash
poflow translate --language sv translations.txt
```

Reads config from `poflow.yml`:
```yaml
gettext_path: "priv/gettext"
```

Resolves to: `priv/gettext/sv/LC_MESSAGES/default.po`

### Direct file path

```bash
cat translations.txt | poflow translate file.po > file_new.po
```

### Translation input format

```
Sign In = Logga in
Sign Out = Logga ut
Welcome = Välkommen
```

## Merge Behavior

1. Read existing .po file
2. Parse translation input
3. Match msgids
4. Update msgstr values
5. Preserve all other data (comments, references, etc.)
6. Write updated .po file

## Error Handling

- Warn if msgid not found in .po file
- Skip invalid lines in translation input
- Preserve original file if errors occur
- Optionally use --force to ignore warnings

## Test Cases

1. Simple translation merge
2. Multiple translations
3. Missing msgids (should warn)
4. Config-based path resolution
5. Direct file path
6. Stdin input
7. Preserve comments/references

## Deliverables

- `poflow translate` command
- Config-based path resolution
- Both stdin and file input support
- Proper error handling and warnings
