# TinyServer Makefile
# Educational TCP/IP and HTTP implementation project

.PHONY: all build test lint clean demo-phase1 demo-phase2 demo-phase3 demo-phase4 demo-phase5 help

# Default target
all: test lint build

# Build targets
build: build-server build-client

build-server:
	@echo "Building server..."
	go build -o cmd/server/server ./cmd/server

build-client:
	@echo "Building client..."
	go build -o cmd/client/client ./cmd/client

# Test targets
test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint and format
lint:
	@echo "Running linters..."
	golangci-lint run

format:
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

vet:
	@echo "Running go vet..."
	go vet ./...

# Demo targets
demo-phase1:
	@echo "Running Phase 1 Demo: TCP Echo Server"
	@./scripts/demo/run-phase1.sh

demo-phase2:
	@echo "Running Phase 2 Demo: HTTP Parser"
	@./scripts/demo/run-phase2.sh

demo-phase3:
	@echo "Running Phase 3 Demo: Simple Server"
	@./scripts/demo/run-phase3.sh

demo-phase4:
	@echo "Running Phase 4 Demo: HTTP Client"
	@./scripts/demo/run-phase4.sh

demo-phase5:
	@echo "Running Phase 5 Demo: Full Stack Chat"
	@./scripts/demo/run-phase5.sh

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	go mod tidy
	go mod download

clean:
	@echo "Cleaning build artifacts..."
	rm -f cmd/server/server
	rm -f cmd/client/client
	rm -f demo/*/server
	rm -f demo/*/client
	rm -f demo/*/main
	rm -f coverage.out coverage.html
	rm -f *.prof *.pprof

# Build demo binaries
build-demos:
	@echo "Building demo binaries..."
	@for demo in demo/phase*; do \
		if [ -f "$$demo/server/main.go" ]; then \
			echo "Building $$demo/server..."; \
			go build -o "$$demo/server/server" "./$$demo/server"; \
		fi; \
		if [ -f "$$demo/client/main.go" ]; then \
			echo "Building $$demo/client..."; \
			go build -o "$$demo/client/client" "./$$demo/client"; \
		fi; \
		if [ -f "$$demo/main.go" ]; then \
			echo "Building $$demo/main..."; \
			go build -o "$$demo/main" "./$$demo"; \
		fi; \
	done

# Check dependencies
deps-check:
	@echo "Checking dependencies..."
	go mod verify
	go mod tidy

# Quick development cycle
quick: format vet test

# Full quality check
quality: clean format vet lint test test-coverage

help:
	@echo "TinyServer Makefile Commands:"
	@echo ""
	@echo "Build commands:"
	@echo "  all          - Run tests, lint, and build"
	@echo "  build        - Build server and client"
	@echo "  build-server - Build main server binary"
	@echo "  build-client - Build main client binary"
	@echo "  build-demos  - Build all demo binaries"
	@echo ""
	@echo "Test commands:"
	@echo "  test         - Run all tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo ""
	@echo "Code quality:"
	@echo "  lint         - Run linters"
	@echo "  format       - Format code with gofmt and goimports"
	@echo "  vet          - Run go vet"
	@echo "  quick        - Format, vet, and test"
	@echo "  quality      - Full quality check with coverage"
	@echo ""
	@echo "Demo commands:"
	@echo "  demo-phase1  - Run TCP Echo Server demo"
	@echo "  demo-phase2  - Run HTTP Parser demo"
	@echo "  demo-phase3  - Run Simple Server demo"
	@echo "  demo-phase4  - Run HTTP Client demo"
	@echo "  demo-phase5  - Run Full Stack Chat demo"
	@echo ""
	@echo "Utility commands:"
	@echo "  dev-setup    - Setup development environment"
	@echo "  deps-check   - Check and tidy dependencies"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"