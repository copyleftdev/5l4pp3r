# Project variables
APP_NAME := 5l4pp3r
PKG := github.com/copyleftdev/$(APP_NAME)
CMD_DIR := ./cmd/$(APP_NAME)

# Tools
GO := go
GOLANGCI_LINT := golangci-lint

# Colors
COLOR_RESET := \033[0m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m
COLOR_RED := \033[31m

# Make is verbose in Linux by default. We don't want that.
.SILENT:

all: build

## build: Build the binary
build:
	@echo "$(COLOR_BLUE)==> Building $(APP_NAME)$(COLOR_RESET)"
	$(GO) build -o $(APP_NAME) $(CMD_DIR)
	@echo "$(COLOR_GREEN)Build complete!$(COLOR_RESET)"

## run: Run the application
run: build
	@echo "$(COLOR_BLUE)==> Running $(APP_NAME)...$(COLOR_RESET)"
	./$(APP_NAME)

## tidy: Cleanup and validate go.mod/go.sum
tidy:
	@echo "$(COLOR_BLUE)==> Tidying go.mod and go.sum$(COLOR_RESET)"
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)Modules tidied!$(COLOR_RESET)"

## test: Run unit tests
test:
	@echo "$(COLOR_BLUE)==> Running tests$(COLOR_RESET)"
	$(GO) test ./... -v

## clean: Remove build artifacts
clean:
	@echo "$(COLOR_BLUE)==> Cleaning up...$(COLOR_RESET)"
	rm -f $(APP_NAME)
	@echo "$(COLOR_GREEN)Cleanup complete!$(COLOR_RESET)"

## lint: Run golangci-lint (if installed)
lint:
	@echo "$(COLOR_BLUE)==> Running golangci-lint$(COLOR_RESET)"
	if ! command -v $(GOLANGCI_LINT) &> /dev/null; then \
		echo "$(COLOR_YELLOW)golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(COLOR_RESET)"; \
		exit 1; \
	fi
	$(GOLANGCI_LINT) run
	@echo "$(COLOR_GREEN)Linting complete!$(COLOR_RESET)"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "$(COLOR_BLUE)Available targets:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(COLOR_GREEN)%-10s$(COLOR_RESET) %s\n", $$1, $$2}'

# Annotate targets with `##` for help message
# (Already done above in comments, just run `make help` to see them)

.PHONY: all build run tidy test clean lint help
