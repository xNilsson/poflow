# Phase 0: Project Setup

**Status:** ✅ COMPLETED
**Goal:** Initialize the poflow project with proper structure and configuration management.

## Tasks

- [x] Initialize Go module
- [x] Install Cobra CLI framework
- [x] Install Viper for configuration management
- [x] Create basic project structure (cmd/, internal/)
- [x] Implement config file reading (poflow.yml/poflow.json)
- [x] Set up version command
- [x] Create basic root command with help text

## Configuration File Structure

The config file should support:

```yaml
# poflow.yml
gettext_path: "priv/gettext"  # Base path for gettext files
# Will resolve to: {gettext_path}/{lang}/LC_MESSAGES/default.po
```

Or JSON format:

```json
{
  "gettext_path": "priv/gettext"
}
```

## Deliverables

- Working Go module
- Basic CLI structure with `poflow --version` and `poflow --help`
- Config file loading from:
  - `./poflow.yml` or `./poflow.json` (current directory)
  - `~/.config/poflow/config.yml` (user config)
