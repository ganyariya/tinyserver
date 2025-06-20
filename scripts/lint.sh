#!/bin/bash
# TinyServer Lint Script

set -e

echo "TinyServer Lint Script"
echo "====================="

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "Warning: golangci-lint is not installed"
    echo "Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
fi

# Format code
echo "Formatting code..."
gofmt -s -w .

# Run goimports if available
if command -v goimports &> /dev/null; then
    echo "Running goimports..."
    goimports -w .
fi

# Run go vet
echo "Running go vet..."
go vet ./...

# Run golangci-lint
echo "Running golangci-lint..."
golangci-lint run

echo ""
echo "Linting completed successfully!"