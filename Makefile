.PHONY: help build run lint lint-fix fmt fmt-check test clean install install-tools

# Default target
help:
	@echo "Available commands:"
	@echo "  make build        - Build the focus binary"
	@echo "  make install      - Build and install the focus binary to /usr/local/bin"
	@echo "  make run          - Run the application"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make lint-fix     - Run golangci-lint with auto-fix"
	@echo "  make fmt          - Format code with gofmt and goimports"
	@echo "  make fmt-check    - Check code formatting without making changes"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make install-tools - Install required development tools"

# Build the binary
build:
	@echo "🔨 Building lolcli..."
	@go build -o bin/lolcli ./cmd/lol-agent

# Install the binary to /usr/local/bin
install: build
	@echo "📦 Installing lol to /usr/local/bin..."
	@sudo cp bin/lolcli /usr/local/bin/lolcli
	@echo "✅ Installation complete! You can now run 'lolcli' from anywhere."

# Run the application
run: build
	@./bin/lolcli start

# Run linter
lint:
	@echo "🔍 Running golangci-lint..."
	@golangci-lint run --config .golangci.yml

# Run linter with auto-fix
lint-fix:
	@echo "🔧 Running golangci-lint with auto-fix..."
	@golangci-lint run --config .golangci.yml --fix

# Format code
fmt:
	@echo "✨ Formatting code..."
	@gofmt -w .
	@if command -v goimports &> /dev/null; then \
		goimports -w -local go-focus-tracker .; \
	else \
		echo "⚠️  goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi
	@echo "✅ Code formatted!"

# Check formatting without making changes
fmt-check:
	@echo "🔍 Checking code formatting..."
	@UNFORMATTED=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "❌ The following files need formatting:"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi
	@echo "✅ All files are properly formatted!"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@go clean

# Start dev server
dev:
	@echo "🔧 Starting dev server..."
	@go run ./cmd/lol-agent start

# Install development tools
install-tools:
	@echo "📦 Installing development tools..."
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint &> /dev/null; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	else \
		echo "✅ golangci-lint already installed"; \
	fi
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ All tools installed!"