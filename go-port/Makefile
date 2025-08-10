# WebGL Water Tutorial Go Port - Makefile

# Variables
BINARY_NAME=webgl-water-server
BINARY_PATH=./cmd/server
BUILD_DIR=./build
DOCKER_IMAGE=webgl-water-go
DOCKER_TAG=latest
GO_VERSION=1.21

# Build flags
BUILD_FLAGS=-ldflags="-s -w"
CGO_ENABLED=0

# Directories
ASSETS_DIR=./assets
WEB_DIR=./web
STATIC_DIR=./web/static
SHADER_DIR=./web/shaders

.PHONY: all build clean test run dev docker docker-build docker-run docker-dev help setup deps

# Default target
all: clean deps build

# Help target
help:
	@echo "WebGL Water Tutorial Go Port - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build        Build the application"
	@echo "  clean        Clean build artifacts"
	@echo "  test         Run tests"
	@echo "  run          Run the application locally"
	@echo "  dev          Run in development mode with hot reload"
	@echo "  setup        Setup development environment"
	@echo "  deps         Download dependencies"
	@echo "  docker       Build and run with Docker"
	@echo "  docker-build Build Docker image"
	@echo "  docker-run   Run Docker container"
	@echo "  docker-dev   Run development Docker setup"
	@echo "  lint         Run linter"
	@echo "  format       Format Go code"
	@echo "  assets       Copy assets from parent directory"
	@echo "  shaders      Validate shader files"
	@echo "  install      Install the binary"
	@echo "  uninstall    Uninstall the binary"
	@echo "  help         Show this help message"

# Setup development environment
setup: deps assets
	@echo "Setting up development environment..."
	@mkdir -p $(BUILD_DIR)
	@mkdir -p $(ASSETS_DIR)
	@mkdir -p $(STATIC_DIR)
	@mkdir -p $(SHADER_DIR)
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@echo "Development environment setup complete!"

# Download dependencies
deps:
	@echo "Downloading Go dependencies..."
	@go mod tidy
	@go mod download

# Build the application
build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all: deps
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	# Linux AMD64
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(BINARY_PATH)

	# Linux ARM64
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(BINARY_PATH)

	# macOS AMD64
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BINARY_PATH)

	# macOS ARM64
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(BINARY_PATH)

	# Windows AMD64
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(BINARY_PATH)

	@echo "Multi-platform build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Tests complete! Coverage report: coverage.html"

# Run tests with benchmarks
test-bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Run the application locally
run: build assets
	@echo "Starting $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode
dev:
	@echo "Starting development server with hot reload..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	@go install $(BINARY_PATH)

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)

# Copy assets from parent directory
assets:
	@echo "Copying assets..."
	@mkdir -p $(ASSETS_DIR)
	@if [ -f ../dudvmap.png ]; then cp ../dudvmap.png $(ASSETS_DIR)/; fi
	@if [ -f ../normalmap.png ]; then cp ../normalmap.png $(ASSETS_DIR)/; fi
	@if [ -f ../stone-texture.png ]; then cp ../stone-texture.png $(ASSETS_DIR)/; fi
	@echo "Assets copied!"

# Validate shader files
shaders:
	@echo "Validating shader files..."
	@for shader in $(SHADER_DIR)/*.glsl; do \
		if [ -f "$$shader" ]; then \
			echo "Checking $$shader..."; \
			# Basic validation - check if file is not empty and has proper extension
			if [ -s "$$shader" ] && [[ "$$shader" == *.glsl ]]; then \
				echo "  ✓ Valid"; \
			else \
				echo "  ✗ Invalid or empty"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "Shader validation complete!"

# Format Go code
format:
	@echo "Formatting Go code..."
	@gofmt -s -w .
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not found. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		goimports -w .; \
	fi

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Please install it:"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"; \
		exit 1; \
	fi

# Security scan
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# Docker targets
docker: docker-build docker-run

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: assets
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 \
		-v $(PWD)/assets:/app/assets \
		-v $(PWD)/web:/app/web \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-dev:
	@echo "Starting development environment with Docker..."
	@docker-compose --profile dev up --build

docker-prod:
	@echo "Starting production environment with Docker..."
	@docker-compose up --build -d

docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-clean:
	@echo "Cleaning Docker artifacts..."
	@docker-compose down --rmi all --volumes --remove-orphans
	@docker system prune -f

# Generate documentation
docs:
	@echo "Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "Documentation server: http://localhost:6060/pkg/github.com/webgl-water-go/go-port/"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Installing..."; \
		go install golang.org/x/tools/cmd/godoc@latest; \
		echo "Documentation server: http://localhost:6060/pkg/github.com/webgl-water-go/go-port/"; \
		godoc -http=:6060; \
	fi

# Generate mesh data from original files (if available)
generate-meshes:
	@echo "Generating mesh data..."
	@go run ./tools/mesh-converter/main.go -input ../meshes.bytes -output $(ASSETS_DIR)/meshes.json

# Benchmark performance
bench:
	@echo "Running performance benchmarks..."
	@go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./...
	@echo "Profiles generated: cpu.prof, mem.prof"

# Profile analysis
profile-cpu:
	@echo "Analyzing CPU profile..."
	@go tool pprof cpu.prof

profile-mem:
	@echo "Analyzing memory profile..."
	@go tool pprof mem.prof

# Release preparation
release: clean format lint test build-all
	@echo "Preparing release..."
	@mkdir -p $(BUILD_DIR)/release
	@cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip release/$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release packages created in $(BUILD_DIR)/release/"

# Quick development cycle
quick: format build run

# CI/CD pipeline simulation
ci: deps format lint test build

# Show project info
info:
	@echo "Project: WebGL Water Tutorial Go Port"
	@echo "Go version: $(shell go version)"
	@echo "Build directory: $(BUILD_DIR)"
	@echo "Binary name: $(BINARY_NAME)"
	@echo "Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
