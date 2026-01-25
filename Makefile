# Multi-Agent Spec - Makefile
#
# This Makefile orchestrates the codegen pipeline:
#   Go types -> JSON Schema -> TypeScript/Zod

.PHONY: all generate generate-schema generate-typescript build test lint clean help

# Default target
all: generate build test

# Generate all schemas (Go -> JSON Schema -> TypeScript)
generate: generate-schema generate-typescript

# Generate JSON Schemas from Go types
generate-schema:
	@echo "Generating JSON Schemas from Go types..."
	cd tools/generate && go run .

# Generate TypeScript/Zod from JSON Schemas
generate-typescript:
	@echo "Generating TypeScript/Zod schemas from JSON Schema..."
	cd sdk/typescript && npm run generate

# Build the TypeScript SDK
build:
	@echo "Building TypeScript SDK..."
	cd sdk/typescript && npm run build

# Run all tests
test: test-go test-typescript

test-go:
	@echo "Running Go SDK tests..."
	cd sdk/go && go test ./...

test-typescript:
	@echo "Running TypeScript SDK tests..."
	cd sdk/typescript && npm test

# Lint JSON Schemas for Go-friendliness
lint-schema:
	@echo "Linting JSON Schemas..."
	@if command -v schemago >/dev/null 2>&1; then \
		schemago lint schema/agent/agent.schema.json && \
		schemago lint schema/orchestration/team.schema.json && \
		schemago lint schema/deployment/deployment.schema.json; \
	else \
		echo "schemago not found, skipping schema lint"; \
	fi

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf sdk/typescript/dist
	rm -rf sdk/typescript/src/generated

# Install dependencies
install:
	@echo "Installing dependencies..."
	cd sdk/typescript && npm install
	cd tools/generate && go mod download

# Help
help:
	@echo "Multi-Agent Spec Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all              Generate, build, and test (default)"
	@echo "  generate         Generate all schemas (Go -> JSON Schema -> TypeScript)"
	@echo "  generate-schema  Generate JSON Schemas from Go types"
	@echo "  generate-typescript  Generate TypeScript/Zod from JSON Schemas"
	@echo "  build            Build the TypeScript SDK"
	@echo "  test             Run all tests (Go + TypeScript)"
	@echo "  test-go          Run Go SDK tests"
	@echo "  test-typescript  Run TypeScript SDK tests"
	@echo "  lint-schema      Lint JSON Schemas for Go-friendliness"
	@echo "  clean            Clean generated files"
	@echo "  install          Install dependencies"
	@echo "  help             Show this help message"
