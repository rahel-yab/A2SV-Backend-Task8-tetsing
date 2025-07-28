.PHONY: test test-coverage lint vet fmt clean build run

# Default target
all: test lint vet fmt

# Run all tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -v -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Check test coverage threshold (80%)
test-coverage-check:
	go test ./... -v -coverprofile=coverage.out -covermode=atomic
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Test coverage: $$COVERAGE%"; \
	if [ $$(echo "$$COVERAGE < 80" | bc -l) -eq 1 ]; then \
		echo "❌ Test coverage is below 80% threshold"; \
		exit 1; \
	else \
		echo "✅ Test coverage meets 80% threshold"; \
	fi

# Run linting
lint:
	go install golang.org/x/lint/golint@latest
	golint ./...

# Run vet
vet:
	go vet ./...

# Check formatting
fmt:
	@if [ "$(shell gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "❌ Code is not properly formatted"; \
		gofmt -s -d .; \
		exit 1; \
	else \
		echo "✅ Code is properly formatted"; \
	fi

# Format code
fmt-fix:
	gofmt -s -w .

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	go clean

# Build the application
build:
	go build -o task-manager-api ./Delivery

# Run the application
run:
	go run ./Delivery

# Install dependencies
deps:
	go mod download
	go mod tidy

# CI pipeline (runs all checks)
ci: deps test-coverage-check lint vet fmt

# Development setup
setup: deps fmt-fix 