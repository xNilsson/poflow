#!/bin/bash
# Integration tests for poflow
# Tests all commands from TUTORIAL.md with real test data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0

# Change to the testdata directory
cd "$(dirname "$0")"

# Build the tool first
echo -e "${YELLOW}Building poflow...${NC}"
go build -o poflow .
echo -e "${GREEN}✓ Build successful${NC}\n"

# Helper function to run a test
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_exit_code="${3:-0}"

    echo -e "${YELLOW}Testing: ${test_name}${NC}"
    echo "  Command: $command"

    if eval "$command" > /dev/null 2>&1; then
        actual_exit_code=0
    else
        actual_exit_code=$?
    fi

    if [ "$actual_exit_code" -eq "$expected_exit_code" ]; then
        echo -e "  ${GREEN}✓ PASS${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "  ${RED}✗ FAIL (exit code: $actual_exit_code, expected: $expected_exit_code)${NC}"
        ((TESTS_FAILED++))
    fi
    echo
}

# Helper function to test command output
run_test_with_output() {
    local test_name="$1"
    local command="$2"
    local expected_pattern="$3"

    echo -e "${YELLOW}Testing: ${test_name}${NC}"
    echo "  Command: $command"

    output=$(eval "$command" 2>&1)

    if echo "$output" | grep -q "$expected_pattern"; then
        echo -e "  ${GREEN}✓ PASS (found: $expected_pattern)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "  ${RED}✗ FAIL (expected pattern: $expected_pattern)${NC}"
        echo "  Output: $output"
        ((TESTS_FAILED++))
    fi
    echo
}

echo "========================================"
echo "  poflow Integration Test Suite"
echo "========================================"
echo

# Test 1: Version command
run_test "poflow version" "./poflow version"

# Test 2: listempty with file
run_test "listempty with file" "./poflow listempty testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 3: listempty with --language flag
run_test "listempty with --language" "./poflow --config testdata/poflow.yml listempty --language sv"

# Test 4: listempty with --json
run_test "listempty with --json" "./poflow listempty --json testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 5: listempty with --limit
run_test "listempty with --limit" "./poflow listempty --limit 1 testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 6: search with plain text
run_test "search with plain text" "./poflow search 'Sign' testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 7: search with --language flag
run_test "search with --language" "./poflow --config testdata/poflow.yml search 'Ask a question' --language sv"

# Test 8: search with --re (regex)
run_test "search with regex" "./poflow search --re '^Sign' testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 9: search with --json
run_test "search with --json" "./poflow search --json 'Sign' testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 10: search with --limit
run_test "search with --limit" "./poflow search --limit 1 'Sign' testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 11: searchvalue with plain text
run_test "searchvalue with plain text" "./poflow searchvalue 'Logga' testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 12: searchvalue with --language flag
run_test "searchvalue with --language" "./poflow --config testdata/poflow.yml searchvalue 'ut' --language sv"

# Test 13: searchvalue with --re (regex)
run_test "searchvalue with regex" "./poflow searchvalue --re '^Välk' testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 14: searchvalue with --json
run_test "searchvalue with --json" "./poflow searchvalue --json 'ut' testdata/gettext/sv/LC_MESSAGES/default.po"

# Test 15: Help commands
run_test "poflow --help" "./poflow --help"
run_test "poflow search --help" "./poflow search --help"
run_test "poflow searchvalue --help" "./poflow searchvalue --help"
run_test "poflow listempty --help" "./poflow listempty --help"

# Test 16: Output validation tests
run_test_with_output "listempty output contains 'Sign In'" \
    "./poflow listempty testdata/gettext/sv/LC_MESSAGES/default.po" \
    "Sign In"

run_test_with_output "search output contains 'Sign Out'" \
    "./poflow search 'Sign' testdata/gettext/sv/LC_MESSAGES/default.po" \
    "Sign Out"

run_test_with_output "searchvalue output contains 'Logga ut'" \
    "./poflow searchvalue 'Logga' testdata/gettext/sv/LC_MESSAGES/default.po" \
    "Logga ut"

# Test 17: JSON output format validation
run_test_with_output "listempty JSON format valid" \
    "./poflow listempty --json testdata/gettext/sv/LC_MESSAGES/default.po" \
    '"msgid"'

run_test_with_output "search JSON format valid" \
    "./poflow search --json 'Sign' testdata/gettext/sv/LC_MESSAGES/default.po" \
    '"msgstr"'

# Test 18: stdin input
run_test_with_output "listempty from stdin" \
    "cat testdata/gettext/sv/LC_MESSAGES/default.po | ./poflow listempty" \
    "Sign In"

run_test_with_output "search from stdin" \
    "cat testdata/gettext/sv/LC_MESSAGES/default.po | ./poflow search 'Welcome'" \
    "Välkommen"

# Print summary
echo "========================================"
echo "  Test Summary"
echo "========================================"
echo -e "Total tests: $((TESTS_PASSED + TESTS_FAILED))"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi
