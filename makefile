# Variables
BINARY_NAME=main
BUILD_DIR=bin

.PHONY: check

# Default target runs everything to ensure quality before building
check: fmt vet

## vet: Examines Go source code and reports suspicious constructs
vet:
	@echo "Vetting code..."
	go vet ./...

## fmt: Formats all Go source files according to standard styles
fmt:
	@echo "Formatting code..."
	go fmt ./...
