# Justfile for teller - STX blockchain CLI tool

# Default recipe to run when just is called
default:
    @just --list

# Variables
BINARY_NAME := "teller"
MAIN_PATH := "cmd/teller/main.go"
BUILD_DIR := "bin"
DOCKER_IMAGE := "teller"
GO_VERSION := "1.24"

# Build the binary
build:
    @echo "Building {{BINARY_NAME}}..."
    go build -o {{BUILD_DIR}}/{{BINARY_NAME}} {{MAIN_PATH}}

# Build for multiple platforms
build-all:
    @echo "Building for multiple platforms..."
    mkdir -p {{BUILD_DIR}}
    GOOS=linux GOARCH=amd64 go build -o {{BUILD_DIR}}/{{BINARY_NAME}}-linux-amd64 {{MAIN_PATH}}
    GOOS=darwin GOARCH=amd64 go build -o {{BUILD_DIR}}/{{BINARY_NAME}}-darwin-amd64 {{MAIN_PATH}}
    GOOS=darwin GOARCH=arm64 go build -o {{BUILD_DIR}}/{{BINARY_NAME}}-darwin-arm64 {{MAIN_PATH}}
    GOOS=windows GOARCH=amd64 go build -o {{BUILD_DIR}}/{{BINARY_NAME}}-windows-amd64.exe {{MAIN_PATH}}

# Run the application
run *args:
    go run {{MAIN_PATH}} {{args}}

# Run with config initialization
run-with-config:
    @echo "Initializing config and running..."
    go run {{MAIN_PATH}} conf init
    go run {{MAIN_PATH}}

# Install dependencies
deps:
    @echo "Installing dependencies..."
    go mod download
    go mod tidy

# Update dependencies
deps-update:
    @echo "Updating dependencies..."
    go get -u ./...
    go mod tidy

# Run tests
test:
    @echo "Running tests..."
    go test ./...

# Run tests with coverage
test-coverage:
    @echo "Running tests with coverage..."
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
    @echo "Running tests with race detection..."
    go test -race ./...

# Benchmark tests
bench:
    @echo "Running benchmarks..."
    go test -bench=. ./...

# Format code
fmt:
    @echo "Formatting code..."
    go fmt ./...
    gofumpt -w .

# Lint code
lint:
    @echo "Linting code..."
    golangci-lint run

# Vet code
vet:
    @echo "Vetting code..."
    go vet ./...

# Run all checks (format, vet, lint, test)
check: fmt vet lint test

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf {{BUILD_DIR}}
    rm -f coverage.out coverage.html
    go clean

# Install the binary to GOPATH/bin
install:
    @echo "Installing {{BINARY_NAME}} to GOPATH/bin..."
    go install {{MAIN_PATH}}

# Docker build
docker-build:
    @echo "Building Docker image..."
    docker build -t {{DOCKER_IMAGE}} .

# Docker run
docker-run *args:
    @echo "Running Docker container..."
    docker run --rm -it {{DOCKER_IMAGE}} {{args}}

# Docker run with volume for config
docker-run-with-config *args:
    @echo "Running Docker container with config volume..."
    docker run --rm -it -v ~/.teller.yaml:/root/.teller.yaml {{DOCKER_IMAGE}} {{args}}

# Generate config file
config-init:
    @echo "Generating default config file..."
    ./{{BUILD_DIR}}/{{BINARY_NAME}} conf init

# Show version
version:
    @echo "Go version: $(go version)"
    @echo "Teller version:"
    ./{{BUILD_DIR}}/{{BINARY_NAME}} --version

# Show help
help:
    @echo "Running teller help..."
    ./{{BUILD_DIR}}/{{BINARY_NAME}} --help

# Development setup
dev-setup:
    @echo "Setting up development environment..."
    go mod download
    mkdir -p {{BUILD_DIR}}
    @echo "Installing development tools..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install mvdan.cc/gofumpt@latest
    @echo "Development setup complete!"

# Hot reload for development (requires entr)
dev-watch:
    @echo "Starting hot reload (requires 'entr' to be installed)..."
    find . -name "*.go" | entr -r just run

# Security audit
audit:
    @echo "Running security audit..."
    go list -json -deps ./... | nancy sleuth

# Show module info
mod-info:
    @echo "Module information:"
    go list -m all
    @echo "\nDirect dependencies:"
    go list -m -json | jq -r .Path

# Generate vendor directory
vendor:
    @echo "Creating vendor directory..."
    go mod vendor

# Verify dependencies
verify:
    @echo "Verifying dependencies..."
    go mod verify

# Performance profiling
profile:
    @echo "Running with CPU profiling..."
    go run {{MAIN_PATH}} -cpuprofile=cpu.prof
    @echo "Profile saved as cpu.prof"

# Memory profiling
profile-mem:
    @echo "Running with memory profiling..."
    go run {{MAIN_PATH}} -memprofile=mem.prof
    @echo "Profile saved as mem.prof"

# Release build (optimized)
release:
    @echo "Building release version..."
    mkdir -p {{BUILD_DIR}}
    go build -ldflags="-w -s" -o {{BUILD_DIR}}/{{BINARY_NAME}} {{MAIN_PATH}}
    @echo "Release binary created: {{BUILD_DIR}}/{{BINARY_NAME}}"

# Show project statistics
stats:
    @echo "Project statistics:"
    @echo "Lines of Go code:"
    find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1
    @echo "Number of Go files:"
    find . -name "*.go" -not -path "./vendor/*" | wc -l
    @echo "Dependencies:"
    go list -m all | wc -l 