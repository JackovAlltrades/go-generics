# Makefile for go-generics project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT?=gofumpt # Allows overriding via environment if needed
GOLINT=golangci-lint
GOCOVER=$(GOTEST) -coverprofile=coverage.out
GOCOVER_HTML=$(GOCMD) tool cover -html=coverage.out
GOINSTALL=$(GOCMD) install

# Project structure
PKG_LIST?=$(shell $(GOCMD) list ./... | grep -v /vendor/)
# Use '.' as a reliable default for the project root if CURDIR is problematic
TARGET_DIR?=.

# Tools - Define paths relative to TARGET_DIR
TOOLS_DIR := $(TARGET_DIR)/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
LINT_VERSION?=v1.64.8 # Use the version compatible with Go 1.22+ selected GOBIN path definition ---
# Get the absolute path using Make's built-in 'abspath' for better cross-platform reliability
# This avoids complex shell expansions during global variable setup.
ABS_TOOLS_BIN_DIR_MAKE := $(abspath $(TOOLS_BIN_DIR))
# Convert the absolute path to use forward slashes (required by Go tools)
ABS_TOOLS_BIN_DIR_UNIX := $(subst \,/,$(ABS_TOOLS_BIN_DIR_MAKE))
# Set GOBIN using this derived absolute Unix-style path
GOBIN?=$(ABS_TOOLS_BIN_DIR_UNIX)
# Define the expected linter path using the relative TOOLS_BIN_DIR
GOLINT_PATH := $(TOOLS_BIN_DIR)/$(GOLINT)

# --- OS Detection and Helpers ---
# Detect OS (Simplified check for Windows)
ifeq ($(OS),Windows_NT)
	# On Windows, GOLINT_EXEC needs the Windows path for execution
	# It's safer to check if the .exe exists explicitly later if needed
	GOLINT_EXEC?=$(shell cygpath -w $(GOLINT_PATH))
	# Define MKDIR_P to handle potential drive letters and use forward slashes via cygpath
	# Always ensure the path passed to cygpath is relative or absolute as needed
	# CORRECTED LINE: Removed single quotes around $$(cygpath...) to allow substitution
	MKDIR_P = mkdir -p $$(cygpath -u "$(1)") # Create using Unix path for consistency
	RM_F=rm -f
	RM_RF=rm -rf
else
	# Assume Unix-like
	GOLINT_EXEC?=$(GOLINT_PATH)
	MKDIR_P = mkdir -p "$(1)"
	RM_F=rm -f
	RM_RF=rm -rf
endif

# Ensure GOPATH/bin and local tools bin are in PATH
export PATH := $(shell $(GOCMD) env GOPATH)/bin:$(PATH):$(ABS_TOOLS_BIN_DIR_UNIX)

# --- Phony Targets ---
# Declare targets that don't represent files
.PHONY: all build test cover cover-html clean fmt lint lint-check \
        install-tools install-lint check bench bench-functional help

# --- Default Target ---
all: check ## Run all checks (fmt, lint, test) by default

# --- Help Target ---
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Build Target ---
build: ## Build the main application (if any - currently none)
	@echo "==> Building..."
	# Placeholder: $(GOBUILD) -o $(TARGET_DIR)/go-generics ./cmd/app # Adjust if needed

# --- Format Target ---
fmt: ## Format code using gofumpt with extra rules
	@echo "==> Formatting code (with extra rules)..."
	$(GOFMT) -extra -w . # Add -extra flag here
# --- Lint Targets ---
lint: install-lint ## Run linters (installs if necessary)
	@echo "==> Linting code..."
	# Use the potentially OS-specific execution path variable
	"$(GOLINT_EXEC)" run ./...

# Target to check if linter is installed and accessible without running it
lint-check: ## Check if linter is installed and show version
	@_LINT_EXEC_PATH=""; \
	if [ -f "$(GOLINT_PATH).exe" ]; then \
		_LINT_EXEC_PATH="$(GOLINT_PATH).exe"; \
	elif [ -f "$(GOLINT_PATH)" ]; then \
		_LINT_EXEC_PATH="$(GOLINT_PATH)"; \
	fi; \
	if [ -n "$$_LINT_EXEC_PATH" ]; then \
		_LINT_EXEC_WIN=$$(cygpath -w $$_LINT_EXEC_PATH 2>/dev/null || echo $$_LINT_EXEC_PATH); \
		echo "$(GOLINT) found at: $$_LINT_EXEC_PATH"; \
		"$$_LINT_EXEC_WIN" --version; \
	else \
		echo "$(GOLINT) not found at $(GOLINT_PATH) or $(GOLINT_PATH).exe. Run 'make install-lint'."; \
		exit 1; \
	fi

# --- Test & Coverage Targets ---
test: ## Run tests with race detector
	@echo "==> Running tests (with race detector)..."
	go test -v -race ./...    # Use the go command directly

cover: ## Generate coverage report
	@echo "==> Generating coverage report..."
	$(GOCOVER) ./...

cover-html: cover ## Open coverage report in browser
	@echo "==> Opening coverage report in browser..."
	$(GOCOVER_HTML)

# --- Combined Check Target ---
check: fmt lint test ## Run format, lint, and test sequentially
	@echo "==> All checks passed."

# --- Benchmark Targets ---
bench: ## Run benchmarks for all packages
	@echo "==> Running all benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

bench-functional: ## Run benchmarks for the functional package only
	@echo "==> Running functional package benchmarks..."
	$(GOTEST) -bench=. -benchmem ./functional/...

# --- Clean Target ---
clean: ## Clean build artifacts and coverage file
	@echo "==> Cleaning..."
	$(GOCLEAN)
	$(RM_F) "$(TARGET_DIR)/go-generics" # Remove binary if exists
	$(RM_F) "$(TARGET_DIR)/coverage.out"
	# $(RM_RF) "$(TOOLS_BIN_DIR)" # Keep tools by default - uncomment to also remove tools/bin on clean

# --- Tool Installation ---
install-tools: install-lint ## Install all required tools locally
	@echo "==> All tools installed."

# Installs golangci-lint locally if missing - This is the file trigger rule
# It depends on the existence of the *file* $(GOLINT_PATH)
$(GOLINT_PATH):
	@echo "==> Installing $(GOLINT) $(LINT_VERSION)..."
	# Explicitly ensure the TOOLS_BIN_DIR exists *before* installing into it
	# Use the OS-aware MKDIR_P helper via 'call'
	@$(call MKDIR_P,$(TOOLS_BIN_DIR))
	@echo "Using GOBIN=$(GOBIN) for go install"
	# Run the installation using the reliably set GOBIN
	GOBIN=$(GOBIN) $(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINT_VERSION)
	@echo "Install command attempt finished."
	# Verify installation success
	@if [ ! -f "$(GOLINT_PATH)" ] && [ ! -f "$(GOLINT_PATH).exe" ]; then \
		echo "Error: Failed to install $(GOLINT). Check Go installation, GOBIN path ($(GOBIN)), TOOLS_BIN_DIR ($(TOOLS_BIN_DIR)) and network."; \
		echo "Listing contents of $(TOOLS_BIN_DIR):"; \
		ls -l "$(TOOLS_BIN_DIR)" || echo "Could not list $(TOOLS_BIN_DIR)"; \
		exit 1; \
	else \
		echo "$(GOLINT) installation successful check PASSED."; \
	fi

# Explicit alias target for install-lint if preferred over just depending implicitly
install-lint: $(GOLINT_PATH) ## Install golangci-lint locally to ./tools/bin
	@echo "$(GOLINT) is installed at $(GOLINT_PATH) (or .exe equivalent)."