# Bug Report: poflow translate outputs entire .po file instead of updating in place

## Date
2025-10-07

## Issue
When running `poflow translate --language sv translations.txt`, the tool outputs the entire .po file to stdout instead of updating the file in place or providing a clear indication of what changed.

## Expected Behavior
One of:
1. Update the .po file directly and show a summary of changes (e.g., "Updated 2 translations")
2. Output instructions on what to do next (e.g., "Redirect output to file: poflow translate --language sv translations.txt > temp.po")
3. Have a `--in-place` or `-i` flag to update directly

## Actual Behavior
The entire .po file (thousands of lines) is dumped to stdout, making it unclear:
- Whether the translations were applied successfully
- Which entries were updated
- What the user should do with the output

## Command Used
```bash
poflow translate --language sv translations.txt
```

## translations.txt content
```
Ask a question = Ställ en fråga
Feedback = Feedback
```

## Suggested Improvements
1. Add a `--dry-run` flag to preview changes without output
2. Add `--in-place` or `-i` flag to update file directly
3. Show a summary of what was changed (e.g., "Updated 2 translations: 'Ask a question', 'Feedback'")
4. If outputting to stdout is intended, document this clearly in help text

## Workaround
The user needs to manually redirect output:
```bash
poflow translate --language sv translations.txt > temp.po
mv temp.po priv/gettext/sv/LC_MESSAGES/default.po
```

But this isn't obvious from the tool's behavior.
