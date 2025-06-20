#!/bin/bash
# TinyServer Build Script

set -e

echo "TinyServer Build Script"
echo "======================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Build server
echo "Building server..."
if [ -f "cmd/server/main.go" ]; then
    go build -o cmd/server/server ./cmd/server
    echo "✓ Server built: cmd/server/server"
else
    echo "⚠ cmd/server/main.go not found, skipping server build"
fi

# Build client
echo "Building client..."
if [ -f "cmd/client/main.go" ]; then
    go build -o cmd/client/client ./cmd/client
    echo "✓ Client built: cmd/client/client"
else
    echo "⚠ cmd/client/main.go not found, skipping client build"
fi

# Build demos
echo "Building demos..."
for demo in demo/phase*; do
    if [ -d "$demo" ]; then
        demo_name=$(basename "$demo")
        
        # Build server if exists
        if [ -f "$demo/server/main.go" ]; then
            go build -o "$demo/server/server" "./$demo/server"
            echo "✓ $demo_name server built"
        fi
        
        # Build client if exists
        if [ -f "$demo/client/main.go" ]; then
            go build -o "$demo/client/client" "./$demo/client"
            echo "✓ $demo_name client built"
        fi
        
        # Build main if exists
        if [ -f "$demo/main.go" ]; then
            go build -o "$demo/main" "./$demo"
            echo "✓ $demo_name main built"
        fi
    fi
done

echo ""
echo "Build completed successfully!"