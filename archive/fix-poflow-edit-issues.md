# Bug Report: poflow edit Issues

**Date**: 2025-10-07
**Command**: `poflow edit "B: Fuller bust | C: Standard | D: Fuller hip" "B: Athletic | C: Standard | D: Curvy"`

## Issue 1: Comment Order Changed

The `poflow edit` command changed the order of comments in the `.po` and `.pot` files.

**Before**:
```
#: lib/monster_construction_web/components/pattern_components.ex:380
#, elixir-autogen, elixir-format
msgid "B: Fuller bust | C: Standard | D: Fuller hip"
```

**After**:
```
# , elixir-autogen, elixir-format
#: lib/monster_construction_web/components/pattern_components.ex:380
msgid "B: Athletic | C: Standard | D: Curvy"
```

**Problems**:
1. The `#, elixir-autogen, elixir-format` line moved from second position to first position
2. A space was added after the `#` (became `# ,` instead of `#,`)

**Impact**: This changes the file structure unnecessarily and may cause issues with gettext tooling that expects standard comment ordering.

## Issue 2: Source Code Not Updated

The `poflow edit` command updated the msgid in translation files but **did not update the source code** that references this string.

**File not updated**: `lib/monster_construction_web/components/pattern_components.ex:380`

**Current state**:
```elixir
{gettext("B: Fuller bust | C: Standard | D: Fuller hip")}
```

**Expected state**:
```elixir
{gettext("B: Athletic | C: Standard | D: Curvy")}
```

**Impact**: The source code still contains the old msgid, which means:
- The translation will not work until the source code is manually updated
- Running `mix gettext.extract` will likely recreate the old msgid
- This defeats the purpose of the `edit` command

## Expected Behavior

`poflow edit` should:
1. Preserve the original comment order in `.po` and `.pot` files
2. Update the msgid in all translation files (✓ works)
3. **Update the source code files** that reference the msgid (✗ missing feature)

## Files Affected

- `web/priv/gettext/default.pot`
- `web/priv/gettext/en/LC_MESSAGES/default.po`
- `web/priv/gettext/sv/LC_MESSAGES/default.po`
- `lib/monster_construction_web/components/pattern_components.ex` (should be affected but wasn't)

## Workaround

Manually update the source code file after running `poflow edit`.
