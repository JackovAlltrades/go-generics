#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e
# Treat unset variables as an error when substituting.
set -u
# Prevent errors in pipelines from being masked.
set -o pipefail

# --- Configuration ---
# Specify the packages/directories to check (./... means all in current dir and subdirs)
PACKAGES_TO_CHECK="./..."
# Default commit message if none is provided
DEFAULT_COMMIT_MESSAGE="chore: Run auto checks and updates"
# Branch to push to
GIT_BRANCH="main" # Or your current development branch

# --- Check for Prerequisites ---
if ! command -v gofumpt &> /dev/null; then
    echo "ERROR: gofumpt could not be found. Please install it."
    echo "       go install mvdan.cc/gofumpt@latest"
    exit 1
fi

if ! command -v golangci-lint &> /dev/null; then
    echo "ERROR: golangci-lint could not be found. Please install it."
    echo "       go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    exit 1
fi

# --- Check Script Arguments ---
# Allow skipping git operations with --no-git
do_git=true
if [[ "${1:-}" == "--no-git" ]]; then
    do_git=false
    echo "INFO: Running checks only (--no-git specified)."
fi

# --- Steps ---

echo "STEP 1: Formatting code with gofumpt..."
gofumpt -w ${PACKAGES_TO_CHECK}
echo "Formatting complete."
echo

echo "STEP 2: Linting code with golangci-lint..."
golangci-lint run ${PACKAGES_TO_CHECK}
echo "Linting complete."
echo

echo "STEP 3: Running tests (with race detector)..."
# The -race flag adds significant overhead but finds concurrency bugs. Remove if too slow for frequent runs.
go test -v -race ${PACKAGES_TO_CHECK}
echo "Tests passed."
echo

# --- Optional Git Operations ---
if [[ "$do_git" == true ]]; then
    echo "STEP 4: Performing Git operations..."

    # Check for uncommitted changes introduced by the script (e.g., formatting)
    if ! git diff --quiet HEAD --; then
        echo "INFO: Code changes detected (likely formatting)."

        echo "Staging changes..."
        git add .

        echo "Committing changes..."
        # Prompt for commit message, using default if empty
        read -p "Enter commit message [Default: ${DEFAULT_COMMIT_MESSAGE}]: " commit_msg
        commit_msg="${commit_msg:-${DEFAULT_COMMIT_MESSAGE}}"
        git commit -m "$commit_msg"

        echo "Pushing changes to branch '${GIT_BRANCH}'..."
        git push origin "${GIT_BRANCH}"

        echo "Git operations complete."
    else
        echo "INFO: No code changes detected by script. Nothing to commit or push."
    fi
else
    echo "INFO: Skipping Git operations."
fi

echo
echo "***************************"
echo "*** All checks passed! ***"
echo "***************************"

exit 0