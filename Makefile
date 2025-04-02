# Makefile for go-generics project

# === Go Tooling ===
GO        := go
GOFMT     := gofmt -w -s
GOTEST    := go test
GOBUILD   := go build
GOTOOLDIR := $(shell $(GO) env GOPATH)/bin

# === Tool Versions (Pin for Reproducibility) ===
# Find latest versions on GitHub releases pages
GOLANGCI_LINT_VERSION := v1.57.2
GOSEC_VERSION         := v2.19.0 # Example gosec version
GOVULNCHECK_VERSION   := latest  # govulncheck often recommended as 'latest'
# GOCOGNIT_VERSION := v0.1.0 # If installing gocognit separately

# === Go Project Info ===
MODULE_PKGS := $(shell $(GO) list ./...)

# === Local Tool Paths (using go install to GOPATH/bin) ===
# Assumes $(GOTOOLDIR) is in PATH
GOLINT_CMD      := $(GOTOOLDIR)/golangci-lint
GOSEC_CMD       := $(GOTOOLDIR)/gosec
GOVULNCHECK_CMD := $(GOTOOLDIR)/govulncheck
# GOCOGNIT_CMD   := $(GOTOOLDIR)/gocognit

# === Phony Targets ===
.PHONY: all fmt vet lint test test-race bench cov coverage coverage-html \
        secure vulncheck \
        complexity doc-check \
        build clean tidy help install-tools check-tools

# === Default Target ===
# Add more checks to the default 'all' target for convenience
all: fmt vet lint secure vulncheck test ## Run fmt, vet, lint, security, vulncheck, test

# === Installation Targets ===
install-tools: install-golangci-lint install-gosec install-govulncheck ## Install required dev tools
	@echo "Tools installed. Ensure '$(GOTOOLDIR)' is in your PATH."

install-golangci-lint:
	@echo "==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

install-gosec:
	@echo "==> Installing gosec $(GOSEC_VERSION)..."
	@$(GO) install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)

install-govulncheck:
	@echo "==> Installing govulncheck $(GOVULNCHECK_VERSION)..."
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

# Optional: If using gocognit standalone
# install-gocognit:
#	@echo "==> Installing gocognit $(GOCOGNIT_VERSION)..."
#	@$(GO) install github.com/uudashr/gocognit/cmd/gocognit@$(GOCOGNIT_VERSION)


# === Check Targets ===
check-tools: check-golangci-lint check-gosec check-govulncheck ## Check if tools are installed

check-golangci-lint:
	@# ... (version check code from previous Makefile response) ...

check-gosec:
	@if ! command -v $(GOSEC_CMD) &> /dev/null; then \
		echo "Error: gosec not found. Run 'make install-tools'."; \
		exit 1; \
	fi
	@echo "gosec found." # Version check can be added if needed

check-govulncheck:
	@if ! command -v $(GOVULNCHECK_CMD) &> /dev/null; then \
		echo "Error: govulncheck not found. Run 'make install-tools'."; \
		exit 1; \
	fi
	@echo "govulncheck found."

# === Development Workflow Targets ===
fmt: ## Format Go source code
	@echo "==> Formatting code..."
	@$(GOFMT) $(shell find . -type f -name '*.go' -not -path "./vendor/*")
	@# Consider adding: $(GOLINT_CMD) run --fix --issues-exit-code 0 --enable gofumpt,goimports ./...
	@echo "Done."

vet: ## Run go vet static analysis
	@echo "==> Running go vet..."
	@$(GO) vet $(MODULE_PKGS)
	@echo "Done."

lint: check-tools ## Run golangci-lint using .golangci.yml config
	@echo "==> Running golangci-lint $(GOLANGCI_LINT_VERSION) with config..."
	@if [ ! -f .golangci.yml ]; then \
		echo "Warning: .golangci.yml not found. Using defaults."; \
	fi
	@$(GOLINT_CMD) run ./...
	@echo "Done."

# Complexity check is now primarily handled within 'make lint' via gocognit/gocyclo linters
# complexity: check-tools ## Run cognitive complexity checks (if using standalone gocognit)
#	@echo "==> Checking cognitive complexity..."
#	@$(GOCOGNIT_CMD) -min 25 $(MODULE_PKGS) # Adjust threshold as needed

secure: check-tools ## Run gosec security scanner
	@echo "==> Running gosec security scan..."
	@$(GOSEC_CMD) ./...
	@echo "Done."

vulncheck: check-tools ## Run govulncheck dependency vulnerability scanner
	@echo "==> Running govulncheck scan..."
	@$(GOVULNCHECK_CMD) ./...
	@echo "Done."

test: vet ## Run unit tests (short)
	@echo "==> Running short tests..."
	@$(GOTEST) -short -race $(MODULE_PKGS) # Include -race by default for tests
	@echo "Done."

test-all: vet ## Run all unit tests (including long-running)
	@echo "==> Running all tests with race detector..."
	@$(GOTEST) -race $(MODULE_PKGS)
	@echo "Done."

bench: ## Run benchmarks
	@echo "==> Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem $(MODULE_PKGS)
	@echo "Done."

cov: coverage ## Alias for coverage

coverage: ## Generate test coverage report (coverage.out)
	@echo "==> Generating coverage report..."
	@$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic $(MODULE_PKGS)
	@echo "Coverage profile generated: coverage.out"

coverage-html: coverage ## Generate and view HTML test coverage report
	@echo "==> Generating HTML coverage report..."
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "HTML report generated: coverage.html"
	@# Try to open the report automatically (platform specific)
	@if command -v xdg-open &> /dev/null; then xdg-open coverage.html; \
	 elif command -v open &> /dev/null; then open coverage.html; \
	 elif command -v start &> /dev/null; then start coverage.html; \
	 else echo "Please open coverage.html in your browser."; fi

doc-check: ## Check for undocumented exported symbols (Conceptual - requires script/tool)
	@echo "==> Checking documentation coverage (Conceptual)..."
	@echo "Note: Implement using godoc-coverage or a custom script."
	@# Example conceptual check:
	@# go list -f '{{if .IsExported}}{{.ImportPath}}: {{.Name}}{{end}}' ./... | grep -v '_test' | \
	@# while read symbol; do \
	@#   if ! godoc -ex $$symbol | grep -q "^$$"; then \
	@#     echo "Undocumented: $$symbol"; \
	@#     FAIL=1; \
	@#   fi; \
	@# done; \
	@# exit $$FAIL
	@echo "Done."


build: ## Build the project (if applicable)
	@echo "==> Building..."
	@# Example: $(GOBUILD) -o build/myexec ./cmd/myexec
	@echo "Done."

clean: ## Remove coverage files and optionally build artifacts
	@echo "==> Cleaning..."
	@rm -f coverage.out coverage.html
	@# rm -rf build/ # If build target creates files/dirs
	@echo "Done."


tidy: ## Tidy go module files
	@echo "==> Tidying go module files..."
	@$(GO) mod tidy
	@echo "Done."

# === Help Target ===
help: ## Display this help screen
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main Workflow Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -v 'install-' | grep -v 'check-' | grep -v 'build' | grep -v 'clean' | grep -v 'tidy' | grep -v 'help' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "Tool Management:"
	@grep -E '^(install-|check-)[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "Other Targets:"
	@grep -E '^(build|clean|tidy):.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'


# Default target executed when 'make' is run without arguments
.DEFAULT_GOAL := help