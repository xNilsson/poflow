# poflow

A lightweight CLI tool for working with GNU gettext `.po` translation files.

**poflow** helps developers, translators, and LLMs navigate large `.po` files — searching, listing, and updating translation entries in a structured, automatable way.

## Features

- 🚀 **Fast**: Streaming parser processes files line by line (no in-memory AST)
- 🔍 **Powerful search**: Search by msgid or msgstr with regex support
- 📝 **List untranslated**: Quickly find entries that need translation
- 🔄 **Merge translations**: Apply translations from text files back to .po files
- ✏️ **Edit source text**: Update msgid across all language files and templates
- 🤖 **LLM-friendly**: JSON output and clear structure for automation
- ⚙️ **Configurable**: Project-specific config files for easy workflow
- 📦 **Portable**: Single Go binary with no dependencies

## Documentation

- **[Installation Guide](INSTALLATION.md)** - Complete installation and setup instructions
- **[Tutorial](TUTORIAL.md)** - Step-by-step guide with real examples
- **[README](README.md)** - This file, full reference documentation

## Quick Start

```bash
# 1. Install poflow
go install github.com/xnilsson/poflow@latest

# 2. Initialize in your project
cd /path/to/your/project
poflow init

# 3. Start using poflow
poflow listempty --language sv
```

For detailed installation instructions, see [INSTALLATION.md](INSTALLATION.md).

For a guided tutorial, see [TUTORIAL.md](TUTORIAL.md).

## Installation

### Method 1: Go Install (Recommended)

```bash
go install github.com/xnilsson/poflow@latest
```

### Method 2: Download Binary

Download from the [releases page](https://github.com/xnilsson/poflow/releases):

```bash
# macOS/Linux
curl -L https://github.com/xnilsson/poflow/releases/latest/download/poflow-$(uname -s)-$(uname -m) -o poflow
chmod +x poflow
sudo mv poflow /usr/local/bin/
```

### Method 3: Build from Source

```bash
git clone https://github.com/xnilsson/poflow.git
cd poflow
go build -o poflow .
sudo mv poflow /usr/local/bin/
```

See [INSTALLATION.md](INSTALLATION.md) for complete setup instructions.

## Basic Usage

```bash
# Initialize config in your project
poflow init

# List all untranslated entries
poflow listempty --language sv

# List first 10 untranslated entries in JSON format
poflow listempty --json --limit 10 translations.po

# Search for entries by msgid
poflow search "Welcome" translations.po
poflow search --re "^Sign" translations.po

# Search for entries by msgstr (translation)
poflow searchvalue "Välkommen" translations.po

# Apply translations from a file
poflow translate --language sv translations.txt

# Update source text (msgid) across all language files
poflow edit "Sign In" "Log In"
```

## Configuration

### Quick Setup with `poflow init`

The easiest way to configure poflow:

```bash
cd your_project
poflow init
```

This auto-detects your project structure and creates `poflow.yml` with appropriate settings.

**Custom path:**
```bash
poflow init --path translations
```

### Manual Configuration

Create a `poflow.yml` file in your project root:

```yaml
# poflow.yml
gettext_path: "priv/gettext"
```

With this config, poflow resolves `.po` files as:
```
{gettext_path}/{lang}/LC_MESSAGES/default.po
```

### Config File Locations

poflow searches for config files in this order:

1. `./poflow.yml` (current directory)
2. `./poflow.json` (current directory)
3. `~/.config/poflow/config.yml` (user home)

You can also specify a config file explicitly:

```bash
poflow --config /path/to/config.yml listempty
```

### Config Examples

**Phoenix/Elixir projects:**
```yaml
gettext_path: "priv/gettext"
```

**Rails projects:**
```yaml
gettext_path: "config/locales"
```

**Custom setup:**
```yaml
gettext_path: "translations"
```

## Commands

### `init` - Initialize Configuration

Create a `poflow.yml` config file with auto-detected settings.

```bash
# Auto-detect project structure
poflow init

# Specify custom path
poflow init --path translations

# Force overwrite existing config
poflow init --path priv/gettext
```

**Output:**
```
Detected Phoenix/Elixir project
✓ Created poflow.yml
✓ Configured gettext_path: priv/gettext

Next steps:
  1. Edit poflow.yml if needed
  2. Run: poflow listempty --language <lang>
  3. See: poflow --help for all commands
```

**What it detects:**
- Phoenix/Elixir projects (`mix.exs` → `priv/gettext`)
- Rails projects (`config/application.rb` → `config/locales`)
- Node.js projects (`package.json` → `translations`)
- Existing gettext directory structures

See [INSTALLATION.md](INSTALLATION.md) for more setup details.

### `listempty` - List Untranslated Entries

List all entries with empty translations (`msgstr`).

```bash
# List all empty entries
poflow listempty file.po

# JSON output (one entry per line)
poflow listempty --json file.po

# Limit to first 10 entries
poflow listempty --limit 10 file.po

# Use config file to resolve path
poflow listempty --language sv --json

# From stdin
cat file.po | poflow listempty
```

**Output format:**

Plain text:
```
msgid "Sign In"
msgstr ""

msgid "Sign Out"
msgstr ""
```

JSON (`--json`):
```json
{"msgid":"Sign In","msgstr":""}
{"msgid":"Sign Out","msgstr":""}
```

### `search` - Search by msgid

Search for translation entries where the msgid matches a pattern.

```bash
# Plain substring search (case-insensitive)
poflow search "Welcome" file.po

# Regex search
poflow search --re "^Login" file.po

# JSON output
poflow search --json "error" file.po

# Limit results
poflow search --limit 5 "button" file.po

# From stdin
cat file.po | poflow search "Welcome"
```

### `searchvalue` - Search by msgstr

Search for translation entries where the msgstr (translation) matches a pattern.

```bash
# Plain substring search
poflow searchvalue "Välkommen" file.po

# Regex search
poflow searchvalue --re "^Tack" file.po

# JSON output
poflow searchvalue --json "fel" file.po

# From stdin
cat file.po | poflow searchvalue "error"
```

### `edit` - Update Source Text

Update a msgid (source text) across all language files and templates.

**Usage:**

```bash
# Update msgid across all files
poflow edit "Sign In" "Log In"

# Preview changes without modifying files
poflow edit --dry-run "Sign In" "Log In"
```

**What it does:**

1. Finds all `.po` files in your gettext directory
2. Finds the `.pot` template file (if it exists)
3. Updates the msgid in all matching entries
4. Preserves translations (msgstr) exactly
5. Shows a summary of changes

**Example:**

```bash
$ poflow edit "Sign In" "Log In"

  ✓ priv/gettext/sv/LC_MESSAGES/default.po (1 entries)
  ✓ priv/gettext/en/LC_MESSAGES/default.po (1 entries)
  ✓ priv/gettext/default.pot (1 entries)

Updated 3 file(s) with 3 total entries
```

**When to use:**

- You want to change the English source text
- The change should apply to all language files
- You need to keep translations intact

**Flags:**

- `--dry-run` - Preview changes without modifying files

### `translate` - Merge Translations

Apply translations from a text file into a `.po` file.

**Translation file format (`translations.txt`):**
```
Sign In = Logga in
Sign Out = Logga ut
Welcome = Välkommen
```

**Usage:**

```bash
# Using config file with translation file
poflow translate --language sv translations.txt

# Using --file / -F flag (explicit)
poflow translate --language sv -F translations.txt
poflow translate --language sv --file translations.txt

# Using heredoc (no file needed)
poflow translate --language sv <<EOF
Sign In = Logga in
Sign Out = Logga ut
EOF

# Using pipe
echo "Sign In = Logga in" | poflow translate --language sv

# Direct file path
poflow translate file.po < translations.txt > file_new.po

# Force mode (continue even if msgid not found)
poflow translate --force --language sv translations.txt
```

**How it works:**

1. Reads the original `.po` file
2. Parses translation pairs from input
3. Updates matching msgid entries with new msgstr values
4. Outputs updated `.po` file to stdout

## Global Flags

All commands support these flags:

- `--json` - Output in JSON format (one entry per line)
- `--config <file>` - Specify config file path
- `--quiet` - Suppress progress output

## Using poflow with LLMs

poflow is designed to be LLM-friendly with structured JSON output and clear commands.

### Recommended Workflow

1. **List untranslated strings:**

```bash
poflow listempty --json --limit 20 --language sv
```

2. **Ask your LLM to translate the entries:**

Send the JSON output to an LLM with a prompt like:

> "Please translate these English strings to Swedish. The input is JSON with msgid and msgstr fields. Output translations in the format: 'English text = Swedish translation' (one per line)."

3. **Save translations to a file:**

Save the LLM's output to `translations.txt`:
```
Sign In = Logga in
Sign Out = Logga ut
Welcome = Välkommen
```

4. **Apply translations:**

```bash
poflow translate --language sv translations.txt
```

### Example LLM Prompt

```
Please translate these English strings to Swedish:

[paste JSON output from poflow listempty --json]

Output format (one per line):
English = Swedish

Keep technical terms, placeholders like %{name}, and formatting intact.
```

### Why JSON Output?

- **Structured**: Easy to parse programmatically
- **Line-delimited**: One entry per line for streaming
- **Complete**: Includes msgid, msgstr, comments, and references
- **Deterministic**: Consistent output for reliable automation

### Programmatic Usage

```bash
# Count untranslated entries
poflow listempty --json file.po | wc -l

# Extract just the msgids
poflow listempty --json file.po | jq -r .msgid

# Filter entries by reference
poflow listempty --json file.po | jq 'select(.references[]? | contains("login"))'

# Pipe through LLM API
poflow listempty --json --limit 10 file.po | \
  llm "Translate to Swedish, output as 'EN = SV'" | \
  poflow translate --language sv
```

## Real-World Examples

### Translate Next 10 Untranslated Strings

```bash
# Get untranslated strings
poflow listempty --json --limit 10 --language sv > to_translate.json

# Have LLM translate them (manually or via API)
cat to_translate.json | llm "Translate to Swedish" > translations.txt

# Apply translations
poflow translate --language sv translations.txt
```

### Find All Login-Related Strings

```bash
# Search by msgid
poflow search --re "(?i)login|sign.in" file.po

# Get JSON for further processing
poflow search --json --re "(?i)login|sign.in" file.po | jq .
```

### Check Translation Coverage

```bash
# Count total entries
total=$(poflow search --json "." file.po | wc -l)

# Count untranslated
untranslated=$(poflow listempty --json file.po | wc -l)

# Calculate percentage
echo "Translation coverage: $(( (total - untranslated) * 100 / total ))%"
```

### Batch Process Multiple Languages

```bash
for lang in sv no da fi; do
  echo "Processing $lang..."
  poflow listempty --json --language $lang | \
    llm "Translate to $lang" | \
    poflow translate --language $lang
done
```

## Troubleshooting

### Config File Not Found

**Problem:** `poflow` can't find your config file.

**Solution:**
- Check config file location: `./poflow.yml` or `~/.config/poflow/config.yml`
- Verify YAML syntax (indentation matters!)
- Use `--config /path/to/config.yml` to specify explicitly

### File Not Found with --language Flag

**Problem:** `poflow translate --language sv` says file not found.

**Solution:**
- Ensure `poflow.yml` exists with `gettext_path` configured
- Check that the resolved path exists: `{gettext_path}/{lang}/LC_MESSAGES/default.po`
- Verify directory structure matches gettext convention

### Regex Pattern Not Matching

**Problem:** `poflow search --re "pattern"` returns no results.

**Solution:**
- Test your regex pattern separately (use a tool like regex101.com)
- Remember that patterns are case-sensitive with `--re` (unless using `(?i)` flag)
- Without `--re`, search is case-insensitive substring matching

### Translation Not Applied

**Problem:** `poflow translate` runs but translation doesn't appear in output.

**Solution:**
- Verify the msgid in your translation file exactly matches the msgid in the .po file
- Check for extra whitespace or special characters
- Use `--force` flag to see warnings about unmatched msgids
- Ensure output is redirected correctly (stdout goes to file or pipe)

### Multi-line Strings Not Parsing

**Problem:** Entries with multi-line msgid or msgstr are truncated.

**Solution:**
- This is expected behavior for text output format
- Use `--json` flag for complete multi-line string support
- JSON output preserves exact string content with proper escaping

### Performance Issues with Large Files

**Problem:** Processing large .po files is slow.

**Solution:**
- Use `--limit` flag to process fewer entries
- Redirect stdin/stdout efficiently (avoid displaying large output in terminal)
- Use `--quiet` flag to suppress progress output
- Consider splitting very large .po files if possible

## .po File Format

poflow works with GNU gettext `.po` files, which have this structure:

```
# Comment
#: reference/to/source.ex:42
msgid "English text"
msgstr "Translated text"

msgid "Another string"
msgstr ""
```

### Supported Features

- ✅ Single-line and multi-line strings
- ✅ Escaped quotes and special characters
- ✅ Comments (translator, extracted, reference)
- ✅ Empty translations
- ✅ msgid and msgstr parsing

### Limitations

- ❌ msgid_plural / msgstr[n] (plural forms)
- ❌ msgctxt (context)
- ❌ Fuzzy flags (#, fuzzy)
- ❌ .pot template files (treated as .po)

Plural forms and context support may be added in future versions.

## Development

### Building

```bash
go build -o poflow .
```

### Running Tests

**Unit Tests:**
```bash
go test ./...
```

**Integration Tests:**
```bash
./test_integration.sh
```

The integration test suite validates all commands from the TUTORIAL.md with real test data. It includes:
- All command flags (--json, --limit, --language, --re)
- File and stdin input modes
- Output format validation
- Config file resolution

### Code Structure

```
poflow/
├── cmd/                   # Cobra commands
│   ├── root.go           # Root command + config
│   ├── listempty.go      # List untranslated
│   ├── search.go         # Search by msgid
│   ├── searchvalue.go    # Search by msgstr
│   ├── translate.go      # Apply translations
│   └── version.go        # Version info
├── internal/
│   ├── config/           # Config file handling
│   ├── parser/           # .po file parser
│   ├── model/            # Data structures
│   └── util/             # Helper functions
└── main.go               # Entry point
```

### Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Created by Nille ([@xnilsson](https://github.com/xnilsson))

## Documentation

- **[INSTALLATION.md](INSTALLATION.md)** - Complete installation and setup guide
- **[TUTORIAL.md](TUTORIAL.md)** - Step-by-step tutorial with real examples
- **[README.md](README.md)** - This file, full reference documentation

## Changelog

### v0.2.0 (October 7, 2025)

**Enhancements:**
- Add `poflow edit` for changing msgid in po files and source code.

### v0.1.1 (October 7, 2025)

**Enhancements:**
- Added `--file` / `-F` flag to `translate` command for explicit file input specification
- Updated help text with heredoc and pipe examples for stdin input
- Documentation improvements across README, INSTALLATION, and TUTORIAL

**Bug Fixes:**
- Fixed multi-line string output in text format for all commands
- Fixed comment output to include proper `# ` prefix
- Fixed JSON output to use `encoding/json` for proper escaping
- Deduplicated output functions into shared `internal/output` package

**Internal:**
- Created `internal/output/output.go` with shared output functions
- Removed ~100 lines of duplicated code across commands
- Improved code maintainability and consistency

### v0.1.0 (October 7, 2025)

Initial release with core features:
- `listempty` - List untranslated entries
- `search` - Search by msgid
- `searchvalue` - Search by msgstr
- `translate` - Merge translations into .po files
- JSON output support for all commands
- Config file support (YAML/JSON)
- Streaming parser for efficient processing
- Full documentation and tutorial

## Links

- [GitHub Repository](https://github.com/xnilsson/poflow)
- [Issue Tracker](https://github.com/xnilsson/poflow/issues)
- [GNU gettext Documentation](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html)

---

**poflow** - workflow utility for gettext .po files 🌊
