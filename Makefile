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
TARGET_DIR=${CURDIR} # Use built-in CURDIR

# Tools
TOOLS_DIR=$(TARGET_DIR)/tools
TOOLS_BIN_DIR=$(TOOLS_DIR)/bin
LINT_VERSION?=v1.59.1 # Updated to a recent version (adjust if needed)

# Correct GOBIN path using forward slashes absolute path
_ := $(shell mkdir -p $(subst \,/,$(TOOLS_BIN_DIR))) # Ensure dir exists for pwd
ABS_TOOLS_BIN_DIR_UNIX=$(shell cd $(TOOLS_BIN_DIR); pwd)
GOBIN?=$(ABS_TOOLS_BIN_DIR_UNIX)
GOLINT_PATH=$(TOOLS_BIN_DIR)/$(GOLINT)

# Detect OS (Simplified check for Windows)
ifeq ($(OS),Windows_NT)
	GOLINT_EXEC?=$(shell cygpath -w $(GOLINT_PATH))
	MKDIR_P = mkdir -p '$$(cygpath -u "$(1)")'
	RM_F=rm -f
	RM_RF=rm -rf
else
	GOLINT_EXEC?=$(GOLINT_PATH)
	MKDIR_P = mkdir -p "$(1)"
	RM_F=rm -f
	RM_RF=rm -rf
endif

# Ensure GOPATH/bin is in PATH for tools like gofumpt if not installed locally
export PATH := $(shell $(GOCMD) env GOPATH)/bin:$(PATH):$(abspath $(TOOLS_BIN_DIR))

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
fmt: ## Format code using gofumpt
	@echo "==> Formatting code..."
	$(GOFMT) -w . # Use '.' to format all in current dir and subdirs

# --- Lint Targets ---
lint: install-lint ## Run linters (installs if necessary)
	@echo "==> Linting code..."
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
	$(GOTEST) -v -race ./... # MODIFIED: Added -race flag

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
	# $(RM_RF) "$(TOOLS_BIN_DIR)" # Keep tools by default

# --- Tool Installation ---
install-tools: install-lint ## Install all required tools locally
	@echo "==> All tools installed."

# Installs golangci-lint locally if missing
# Note: Renamed 'check-lint' target above to avoid potential confusion with this rule
$(GOLINT_PATH):
	@echo "==> Installing $(GOLINT) $(LINT_VERSION)..."
	@$(call MKDIR_P,$(TOOLS_BIN_DIR))
	@echo "Using GOBIN=$(GOBIN) for go install"
	GOBIN=$(GOBIN) $(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINT_VERSION)
	@echo "Install command attempt finished."
	@if [ ! -f "$(GOLINT_PATH)" ] && [ ! -f "$(GOLINT_PATH).exe" ]; then \
		echo "Error: Failed to install $(GOLINT). Check Go installation, GOBIN path ($(GOBIN)), TOOLS_BIN_DIR ($(TOOLS_BIN_DIR)) and network."; \
		echo "Listing contents of $(TOOLS_BIN_DIR):"; \
		ls -l "$(TOOLS_BIN_DIR)" || echo "Could not list $(TOOLS_BIN_DIR)"; \
		exit 1; \
	else \
		echo "$(GOLINT) installation successful check PASSED."; \
	fi

# Explicit alias for install-lint if preferred over just depending on $(GOLINT_PATH)
install-lint: $(GOLINT_PATH) ## Install golangci-lint locally to ./tools/bin
	@echo "$(GOLINT) is installed at $(GOLINT_PATH) (or .exe equivalent)."