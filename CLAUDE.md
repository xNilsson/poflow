# CLAUDE.md

This file provides guidance to Claude Code when working on the poflow project.

## Starting the Day

When beginning work on poflow, follow these steps:

1. **Read the current day's log** (if it exists): `logs/YYYY-MM-DD.md`
2. **Check the main plan**: `plan.md` for overall project vision
3. **Review the active phase plan**: Look in `plans/` for the current phase being worked on
4. **Check Phase 0 status**: `plans/phase-0-setup.md` to see what's completed
5. **Run the tool**: `./poflow --help` to verify it builds and works

### Creating Daily Logs

At the start of each session:
- Create or update `logs/YYYY-MM-DD.md`
- Document what you're working on
- Track key decisions and progress
- Note any blockers or questions

## Project Overview

**poflow** is a lightweight CLI tool for working with GNU gettext `.po` translation files.

### Purpose
Help developers, translators, and LLMs navigate large `.po` files — searching, listing, and updating translation entries in a structured, automatable way.

### Key Design Principles
- **Fast**: Stream files line by line (no in-memory AST)
- **LLM-friendly**: Clear help text, JSON output, predictable behavior
- **Configurable**: Use config files to set project-specific paths
- **Portable**: Single Go binary, no dependencies

### Target Usage Pattern

```bash
# List untranslated entries
poflow listempty --language sv --json

# Search for specific strings
poflow search "Welcome" --json

# Apply translations
poflow translate --language sv translations.txt
```

## Project Structure

```
poflow/
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command + config loading
│   ├── version.go         # Version command
│   ├── search.go          # (TODO) Search by msgid
│   ├── searchvalue.go     # (TODO) Search by msgstr
│   ├── listempty.go       # (TODO) List untranslated
│   └── translate.go       # (TODO) Merge translations
├── internal/
│   ├── config/            # Config file handling
│   │   └── config.go
│   ├── parser/            # .po file streaming parser
│   │   └── po.go          # (TODO)
│   ├── model/             # Data structures
│   │   └── entry.go       # MsgEntry struct
│   └── util/              # Helper functions
│       └── io.go          # (TODO)
├── plans/                 # Phase-based planning docs
│   ├── phase-0-setup.md
│   ├── phase-1-parsing.md
│   ├── phase-2-searching.md
│   ├── phase-3-translation.md
│   └── phase-4-documentation.md
├── archive/               # Completed plans
├── logs/                  # Daily development logs
├── main.go               # Entry point
├── plan.md               # Master plan document
└── example.poflow.yml    # Example configuration
```

## Development Phases

### Phase 0: Project Setup ✅ COMPLETED
- Go module initialized
- Cobra + Viper configured
- Basic CLI structure with version command
- Config file support implemented

### Phase 1: Core Parsing + listempty (CURRENT)
- Implement .po file streaming parser
- Create MsgEntry iterator
- Build `listempty` command
- Add JSON output support

### Phase 2: Searching
- `search` command (search msgid)
- `searchvalue` command (search msgstr)
- Regex and substring matching

### Phase 3: Translation
- `translate` command to merge translations
- Config-based path resolution
- Support `--language` flag

### Phase 4: Documentation + LLM Integration
- Write comprehensive README
- Add LLM usage guide
- Implement `--json-help`

## Configuration

The tool uses Viper to read config files in this order:
1. `./poflow.yml` (current directory)
2. `./poflow.json` (current directory)
3. `~/.config/poflow/config.yml` (user home)

### Config Format

```yaml
gettext_path: "priv/gettext"
```

This resolves to: `priv/gettext/{lang}/LC_MESSAGES/default.po` when using `--language` flag.

## Common Development Commands

```bash
# Build the tool
go build -o poflow .

# Run tests
go test ./...

# Run the tool
./poflow --help
./poflow version

# Install locally (optional)
go install

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```

## Code Style Guidelines

### Error Handling
- Return errors, don't panic (except in main/init)
- Use `fmt.Errorf` with `%w` for error wrapping
- Provide helpful error messages

### CLI Commands
- Use Cobra's `RunE` for commands that can error
- Support both file and stdin input where applicable
- Always support `--json` flag for structured output
- Keep help text educational, not just terse

### Parser Design
- Stream files line by line (don't load entire file)
- Use `bufio.Scanner` for reading
- Yield entries one at a time (iterator pattern)
- Handle multi-line strings and escaped quotes

### Testing
- Write tests for parser (critical component)
- Test edge cases: empty files, malformed entries, multi-line strings
- Use table-driven tests for multiple cases

## When Working on New Commands

1. Create `cmd/{command}.go` with Cobra command
2. Add command to root in `init()` function
3. Implement business logic in `internal/` packages
4. Support `--json` flag for output
5. Write helpful help text with examples
6. Update relevant phase plan in `plans/`
7. Update daily log with progress

## LLM-Specific Notes

This tool is designed to be LLM-friendly:
- All commands should support `--json` output
- JSON output should be line-delimited (one entry per line)
- Help text should be clear and include examples
- Commands should be composable (Unix philosophy)
- Prefer deterministic output

## Best Practices

### When Starting a New Feature
1. Read the relevant phase plan in `plans/`
2. Check if there are any notes in today's log
3. Update the phase plan checklist as you work
4. Document decisions in the daily log

### When Completing a Phase
1. Mark all tasks complete in phase plan
2. Update phase status to "COMPLETED"
3. Move phase plan to `archive/`
4. Update `plan.md` if needed
5. Update daily log with summary

### Daily Log Format
Use this structure for daily logs:

```markdown
# Development Log - Month Day, Year

## Session N: Brief Description

### Goals
- List what you plan to work on

### Completed
- ✅ What was completed

### In Progress
- 🔄 What is ongoing

### Blocked
- ❌ Any blockers

### Key Decisions
- Important architectural or design decisions

### Next Steps
- What to do next

### Notes
- Any other relevant information
```

## Architecture Decisions

### Why Streaming Parser?
.po files can be large (thousands of entries). Streaming avoids memory issues and enables fast processing.

### Why Cobra + Viper?
- Cobra: Industry standard for Go CLIs, excellent help generation
- Viper: Flexible config management with multiple format support

### Why Config File?
While the original plan was stateless (stdin/stdout), a config file makes the tool much easier to use in real projects. Users can run `poflow translate --language sv translations.txt` without specifying full paths each time.

## Troubleshooting

### Build Errors
```bash
go mod tidy
go build -o poflow .
```

### Config Not Loading
- Check config file location: `./poflow.yml` or `~/.config/poflow/config.yml`
- Verify YAML syntax
- Run with `--config /path/to/config.yml` to specify explicitly

### Import Errors
- Run `go mod tidy`
- Check that all imports use `github.com/nille/poflow/...`

## Resources

- [Cobra Documentation](https://github.com/spf13/cobra)
- [Viper Documentation](https://github.com/spf13/viper)
- [GNU gettext .po file format](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html)
- [Go bufio.Scanner](https://pkg.go.dev/bufio#Scanner)

## Contact & Contributing

This is Nille's personal project. When making changes:
- Keep the code simple and readable
- Follow Go idioms and conventions
- Write tests for critical functionality
- Update documentation as you go
- Log your work in daily logs
