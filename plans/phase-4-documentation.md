# Phase 4: Documentation + LLM Integration

**Status:** Not Started
**Goal:** Create comprehensive documentation and optimize for LLM usage.

## Tasks

- [ ] Write comprehensive README.md
- [ ] Add usage examples for all commands
- [ ] Create LLM prompt guide section
- [ ] Implement --json-help flag
- [ ] Ensure deterministic help output
- [ ] Add installation instructions
- [ ] Document configuration file format
- [ ] Create example config files
- [ ] Add troubleshooting section
- [ ] Write CONTRIBUTING.md

## README Structure

1. **Introduction** - What is poflow?
2. **Installation** - Binary releases, Homebrew, build from source
3. **Quick Start** - Basic usage examples
4. **Configuration** - Config file setup and options
5. **Commands** - Detailed command documentation
6. **LLM Integration** - How to use poflow with LLMs
7. **Examples** - Real-world workflows
8. **Development** - How to contribute

## LLM Integration Guide

Include in README:

```markdown
### Using poflow with LLMs

poflow is designed to be LLM-friendly with structured JSON output and clear commands.

#### Recommended workflow:

1. List untranslated strings:
   ```bash
   poflow listempty --json --limit 20 --language sv
   ```

2. Ask LLM to translate the JSON entries

3. Format translations as:
   ```
   msgid 1 = translation 1
   msgid 2 = translation 2
   ```

4. Apply translations:
   ```bash
   poflow translate --language sv translations.txt
   ```

#### Example prompt for LLMs:

"Please translate these English strings to Swedish. Output format: 'English = Swedish' one per line."
```

## Help Text Improvements

- Make help text educational, not just terse
- Include examples in help output
- Explain when to use JSON mode
- Show common workflows

## Deliverables

- Complete README.md
- LLM integration guide
- Example configuration files
- Clear, educational help text
- --json-help support for programmatic access
