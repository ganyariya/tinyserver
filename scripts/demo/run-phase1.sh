#!/bin/bash
# Phase 1 TCP Echo Server Demo Script

set -e

echo "TinyServer Phase 1 Demo: TCP Echo Server"
echo "========================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Get project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"
echo ""

# Build the demo binaries
echo "Building demo binaries..."
go build -o demo/phase1-tcp-echo/server/server ./demo/phase1-tcp-echo/server
go build -o demo/phase1-tcp-echo/client/client ./demo/phase1-tcp-echo/client
echo "✓ Binaries built successfully"
echo ""

# Start server in background
echo "Starting TCP Echo Server..."
./demo/phase1-tcp-echo/server/server -port 8080 &
SERVER_PID=$!

# Function to cleanup server on exit
cleanup() {
    echo ""
    echo "Cleaning up..."
    if kill -0 $SERVER_PID 2>/dev/null; then
        echo "Stopping server (PID: $SERVER_PID)..."
        kill $SERVER_PID
        wait $SERVER_PID 2>/dev/null || true
    fi
    echo "Cleanup completed"
}

# Set trap to cleanup on script exit
trap cleanup EXIT

# Wait for server to start
echo "Waiting for server to start..."
sleep 2

# Check if server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "Error: Server failed to start"
    exit 1
fi

echo "✓ Server started successfully (PID: $SERVER_PID)"
echo ""

# Test with multiple messages
echo "Testing TCP Echo Server with sample messages..."
echo ""

test_messages=(
    "Hello, TinyServer!"
    "This is my first TCP implementation!"
    "TCP Echo test successful!"
    "Goodbye from TinyServer!"
)

for i in "${!test_messages[@]}"; do
    message="${test_messages[$i]}"
    echo "Test $((i+1)): Sending message: \"$message\""
    
    # Send message and capture output
    output=$(./demo/phase1-tcp-echo/client/client -message "$message" 2>&1)
    
    # Check if the test was successful
    if echo "$output" | grep -q "✓ Echo successful!"; then
        echo "✓ Test $((i+1)) passed"
    else
        echo "✗ Test $((i+1)) failed"
        echo "Output: $output"
        exit 1
    fi
    
    echo ""
    sleep 1
done

echo "All tests passed! 🎉"
echo ""

# Test concurrent connections
echo "Testing concurrent connections..."
echo "Starting 3 concurrent clients..."

pids=()
for i in {1..3}; do
    (
        sleep 0.5
        ./demo/phase1-tcp-echo/client/client -message "Concurrent client $i" >/dev/null 2>&1
    ) &
    pids+=($!)
done

# Wait for all concurrent clients to finish
for pid in "${pids[@]}"; do
    wait $pid
done

echo "✓ Concurrent connections test passed"
echo ""

# Performance test
echo "Running performance test..."
echo "Sending 10 messages rapidly..."

start_time=$(date +%s.%N)
for i in {1..10}; do
    ./demo/phase1-tcp-echo/client/client -message "Performance test $i" >/dev/null 2>&1
done
end_time=$(date +%s.%N)

duration=$(echo "$end_time - $start_time" | bc)
echo "✓ Performance test completed in ${duration} seconds"
echo ""

echo "Demo Summary:"
echo "============="
echo "✓ Server startup: SUCCESS"
echo "✓ Basic echo functionality: SUCCESS"
echo "✓ Multiple messages: SUCCESS"
echo "✓ Concurrent connections: SUCCESS"
echo "✓ Performance test: SUCCESS"
echo ""

echo "Manual Testing:"
echo "==============="
echo "The server is still running. You can manually test it by running:"
echo "  ./demo/phase1-tcp-echo/client/client"
echo ""
echo "This will start interactive mode where you can type messages."
echo "Type 'quit' to exit the interactive client."
echo ""
echo "Press Enter to stop the demo and shutdown the server..."
read -r

echo ""
echo "Demo completed successfully! 🎉"
echo ""
echo "What was demonstrated:"
echo "- TCP server can accept multiple client connections"
echo "- Client can connect and send messages to server"
echo "- Server correctly echoes back received messages"
echo "- Proper connection handling and cleanup"
echo "- Support for both single-message and interactive modes"
echo ""
echo "Next step: Phase 2 - HTTP Protocol Implementation"