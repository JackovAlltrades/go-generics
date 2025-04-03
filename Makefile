# Makefile for the go-generics project

# === Variables ===

# Project name (optional, can be used for build output)
PROJECT_NAME := go-generics

# Go commands
GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOLINT := $(GO) lint  # Note: golangci-lint is generally preferred over 'go lint'

# Tools directory
TOOLS_DIR := $(shell pwd)/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin

# Ensure tools are runnable from ./tools/bin
export PATH := $(TOOLS_BIN_DIR):$(PATH)

# Tool versions
GOLANGCILINT_VERSION := v1.64.8

# Tool binaries
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint

# Default target executed when you just run 'make'
.DEFAULT_GOAL := help

# Phony targets (targets that don't represent actual files)
.PHONY: all build test lint install-tools clean help

# === Targets ===

all: build test lint ## Build, test, and lint the project

build: ## Build the project (adjust if it's a library vs application)
	@echo "==> Building..."
	$(GOBUILD) ./...

test: ## Run unit tests with race detector
	@echo "==> Running tests..."
	$(GOTEST) -v -race ./...

lint: install-golangci-lint ## Run golangci-lint
	@echo "==> Linting..."
	$(GOLANGCI_LINT) run ./...

# --- Tool Installation ---

# This acts as a dependency check. If the binary exists, the target is satisfied.
$(GOLANGCI_LINT):
	@echo "==> Installing golangci-lint $(GOLANGCILINT_VERSION)..."
	@mkdir -p $(TOOLS_BIN_DIR)
	GOBIN=$(TOOLS_BIN_DIR) $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

install-golangci-lint: $(GOLANGCI_LINT) ## Install golangci-lint locally

install-tools: install-golangci-lint ## Install all required tools

# --- Cleanup ---

clean: ## Remove tools directory and potentially build artifacts
	@echo "==> Cleaning..."
	@rm -rf $(TOOLS_DIR)
	# Add other clean targets if needed, e.g., rm -f bin/$(PROJECT_NAME)

# --- Help ---

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
