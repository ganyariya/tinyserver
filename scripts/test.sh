#!/bin/bash
# TinyServer Test Script

set -e

echo "TinyServer Test Runner"
echo "====================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Run tests
echo "Running tests..."
go test -v ./...

# Run tests with race detection
echo ""
echo "Running tests with race detection..."
go test -race ./...

# Generate coverage report if requested
if [ "$1" = "--coverage" ]; then
    echo ""
    echo "Generating coverage report..."
    go test -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated: coverage.html"
fi

echo ""
echo "All tests passed!"