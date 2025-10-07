## Project Plan: **poflow**

**Purpose:**
`poflow` is a lightweight CLI for working with GNU gettext `.po` translation files.
It helps developers, translators, and LLMs navigate large `.po` files — searching, listing, and updating translation entries in a structured, automatable way.

---

### 🧭 Design principles

* **Fast:** Stream files line by line (no in-memory AST needed).
* **Predictable:** Deterministic text output and JSON modes.
* **Few clear commands:** Fewer verbs, more power.
* **LLM-friendly:** Clear help text, descriptive JSON, and short examples.
* **Configurable:** Use a config file to set paths for .po/.pot files.
* **Stateless:** Input and output via stdin/stdout by default.
* **Portable:** Single binary (Go) without dependencies.

---

### 🚀 MVP Goals (v0.1)

Implement a CLI that can:

1. **Search** `.po` files by `msgid` (regex or substring).
2. **Searchvalue** by `msgstr`.
3. **Listempty** — list untranslated or empty entries.
4. **Translate** — take a list of translated lines and merge them into a `.po` file (pipe or file).
5. Provide structured, LLM-parseable help output (`--json-help` or just a predictable `--help`).

---

### 📁 Project Layout

```
poflow/
 ├── cmd/
 │   ├── root.go          # cobra/viper setup
 │   ├── search.go
 │   ├── searchvalue.go
 │   ├── listempty.go
 │   ├── translate.go
 │   └── version.go
 ├── internal/
 │   ├── config/
 │   │   └── config.go    # config file handling
 │   ├── parser/
 │   │   └── po.go        # .po streaming parser
 │   ├── model/
 │   │   └── entry.go     # MsgEntry struct
 │   └── util/
 │       └── io.go        # helpers for reading stdin/out
 ├── plans/               # phase-based planning docs
 ├── archive/             # completed plans
 ├── logs/                # daily development logs
 ├── README.md
 ├── go.mod
 ├── main.go
 └── LICENSE
```

Use **Cobra** for CLI (clean help text, composable commands).
Use **Viper** for configuration file management.

---

### 🧱 Core data structure

```go
type MsgEntry struct {
    MsgID     string   `json:"msgid"`
    MsgStr    string   `json:"msgstr"`
    Comments  []string `json:"comments,omitempty"`
    References []string `json:"references,omitempty"`
}
```

Parsing approach:

* Detect new entry at `msgid`.
* Accumulate until blank line.
* Handle multi-line strings and escaped quotes.

---

### 🧩 Commands

#### `poflow search`

Find messages by ID pattern.

```bash
poflow search --re "Pattern" file.po
poflow search --plain "Welcome" --json
```

Output:

```
#: lib/web/live/page.ex:24
msgid "Welcome"
msgstr "Välkommen"
```

or JSON mode:

```json
{"msgid":"Welcome","msgstr":"Välkommen"}
```

#### `poflow searchvalue`

Search in translations (`msgstr`).

```bash
poflow searchvalue --re "Tack" file.po
```

#### `poflow listempty`

List entries with empty translations.

```bash
poflow listempty --limit 10 file.po
poflow listempty --lang sv --json
```

Example output:

```
msgid "Sign In"
msgstr ""
```

JSON example:

```json
{"msgid":"Sign In","msgstr":""}
```

#### `poflow translate`

Read translations from stdin (or file) and apply to `.po`.

Example input (translations.txt):

```
Sign In = Logga in
Sign Out = Logga ut
```

Command:

```bash
poflow translate --language sv translations.txt
```

The command will automatically locate the correct `.po` file based on the configuration file settings. The config file (`poflow.yml` or `poflow.json`) should specify the base path for gettext files:

```yaml
# poflow.yml
gettext_path: "priv/gettext"  # Will look for priv/gettext/{lang}/LC_MESSAGES/default.po
```

Alternative usage with explicit file:

```bash
cat translations.txt | poflow translate file.po > file_new.po
```

---

### ⚙️ Global flags

```
--json           Output in JSON lines (one entry per line)
--lang <code>    Optional language tag (sv, fr, etc.)
--limit <n>      Limit number of entries
--re             Interpret pattern as regex
--plain          Plain substring match (default)
--quiet          Suppress progress output
```

---

### 📘 LLM Education Strategy

#### 1. Help Text (`poflow help`)

Make help text *didactic*, not just terse:

```
poflow: workflow utility for gettext .po files

Usage:
  poflow [command] [flags]

Commands:
  search        Search msgid by regex or text
  searchvalue   Search msgstr by regex or text
  listempty     List untranslated entries
  translate     Apply new translations from input

All commands can output JSON lines for programmatic or LLM usage.
Example:
  poflow listempty --json | jq .

Use `--json-help` for machine-readable help metadata.
```

#### 2. Recommended prompt for LLMs (in README)

Add this snippet:

```
### Recommended prompt for LLMs

You can use `poflow` to query and edit gettext `.po` files.  
Always use JSON output for structured access.

Examples:
- To list untranslated strings:
  `poflow listempty --json --limit 20 priv/gettext/sv/LC_MESSAGES/default.po`
- To search for a specific key:
  `poflow search --re "pattern" --json file.po`
```

---

### 🧩 Stretch goals (v0.2+)

* `poflow stats` → count total, translated, untranslated.
* `poflow merge` → combine `.po` fragments.
* `poflow json` → export full `.po` as structured JSON for LLM batch translation.
* `poflow apply` → apply translations from JSON back to `.po`.
* Multi-language directory mode: operate on all files in `priv/gettext/*/LC_MESSAGES/*.po`.

---

### 🧠 Development Plan

**Phase 0: Project Setup**

* Initialize Go module.
* Set up Cobra + Viper.
* Create project structure (plans/, archive/, logs/).
* Implement config file reading (YAML/JSON support).

**Phase 1: Core parsing + listempty**

* Implement `MsgEntry` parsing.
* Build iterator that yields entries one by one.
* Implement `listempty`.

**Phase 2: Searching**

* Add `search` and `searchvalue` commands.
* Add regex/substring support.
* Add `--json` output.

**Phase 3: Translation**

* Implement `translate` subcommand merging translations into existing `.po`.
* Support config-based path resolution with `--language` flag.

**Phase 4: Documentation + LLM integration**

* Write `README.md` with examples and LLM prompt guide.
* Add `--json-help` and ensure help messages are deterministic.

---

### 🔧 Example small snippet (Go)

```go
func ParsePO(r io.Reader) ([]MsgEntry, error) {
    scanner := bufio.NewScanner(r)
    var entries []MsgEntry
    var e MsgEntry
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if strings.HasPrefix(line, "msgid") {
            if e.MsgID != "" {
                entries = append(entries, e)
                e = MsgEntry{}
            }
            e.MsgID = unquote(line[5:])
        } else if strings.HasPrefix(line, "msgstr") {
            e.MsgStr = unquote(line[6:])
        } else if line == "" && e.MsgID != "" {
            entries = append(entries, e)
            e = MsgEntry{}
        }
    }
    if e.MsgID != "" {
        entries = append(entries, e)
    }
    return entries, scanner.Err()
}
```

---

### 🧾 License & Packaging

* License: MIT
* Language: Go (1.23+)
* CLI framework: `spf13/cobra`
* Distribution: Homebrew tap, standalone binary, GitHub Releases
