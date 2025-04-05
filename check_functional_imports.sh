#!/bin/bash

PROJECT_ROOT="." # Assumes run from C:\MyProjects-25\go-generics
FUNCTIONAL_DIR="${PROJECT_ROOT}/functional"
MODULE_PATH="github.com/JackovAlltrades/go-generics" # Your module path

echo "--- Checking non-test *.go files in ${FUNCTIONAL_DIR} ---"
echo

# Use find to get only .go files NOT ending in _test.go
find "${FUNCTIONAL_DIR}" -name '*.go' -not -name '*_test.go' | while IFS= read -r file; do
  echo "Checking file: $file"

  # 1. Check for direct import of the test package (very unlikely, but a direct cycle)
  if grep -q "import \".*${MODULE_PATH}/functional_test\"" "$file"; then
    echo "  ERROR: Found direct import of test package in $file"
  fi
  # Check for aliased import too
  if grep -q "\w\+ \".*${MODULE_PATH}/functional_test\"" "$file"; then
    echo "  ERROR: Found aliased import of test package in $file"
  fi

  # 2. Check for usage of potentially test-only identifiers (ptr, person)
  #    Use word boundaries (\b) to avoid matching parts of other words.
  if grep -q -E '\b(ptr|person)\b' "$file"; then
    echo "  WARNING: Found usage of 'ptr' or 'person' in $file."
    echo "           Verify these are NOT the helpers defined only in helpers_test.go."
    grep -n -E '\b(ptr|person)\b' "$file" # Show line numbers
  fi

done

echo
echo "--- Check Complete ---"