#!/bin/bash

# Phase 2: HTTP Parser Demo Script
# This script runs the HTTP parser and analyzer demo

set -e

echo "==========================================="
echo "Phase 2: HTTP Parser & Analyzer Demo"
echo "==========================================="
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "demo/phase2-http-parser" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    echo "   Expected: /path/to/tinyserver/"
    echo "   Current:  $(pwd)"
    exit 1
fi

echo "🔨 Building project..."
if ! go build ./...; then
    echo "❌ Build failed! Please fix compilation errors first."
    exit 1
fi

echo "✅ Build successful!"
echo ""

echo "🚀 Starting HTTP Parser Demo..."
echo ""

# Run the main demo
go run demo/phase2-http-parser/main.go

echo ""
echo "==========================================="
echo "✨ Demo Complete!"
echo "==========================================="
echo ""
echo "What you just saw:"
echo "  ✓ HTTP request parsing from raw text"
echo "  ✓ Query parameter extraction"
echo "  ✓ Header analysis and validation"
echo "  ✓ Request body handling"
echo "  ✓ HTTP response generation"
echo "  ✓ Error handling for malformed requests"
echo ""
echo "📁 Sample files available at:"
echo "   demo/phase2-http-parser/samples/"
echo ""
echo "🎓 Next Phase: Phase 3 - Simple HTTP Server"
echo "   Run: ./scripts/demo/run-phase3.sh (when ready)"
echo ""