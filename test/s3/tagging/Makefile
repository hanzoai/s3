# S3 Object Tagging Tests
# Tests for GitHub issue #7589: S3 object Tags query comes back empty

.PHONY: help build-s3 setup-server start-server stop-server test-tagging test-tagging-quick test-tagging-comprehensive test-all clean logs check-deps test-with-server

# Configuration
WEED_BINARY := ../../../s3/s3_binary
S3_PORT := 8006
MASTER_PORT := 9338
VOLUME_PORT := 8085
FILER_PORT := 8893
TEST_TIMEOUT := 10m
TEST_PATTERN := TestObjectTaggingOnUpload|TestPutObjectTaggingAPI|TestDeleteObjectTagging|TestTag

# Default target
help:
	@echo "S3 Object Tagging Tests Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  help                       - Show this help message"
	@echo "  build-s3                - Build the Hanzo binary"
	@echo "  check-deps                - Check dependencies and build binary if needed"
	@echo "  start-server              - Start Hanzo server for testing"
	@echo "  stop-server               - Stop Hanzo server"
	@echo "  test-tagging              - Run all tagging tests"
	@echo "  test-tagging-quick        - Run core tagging tests only"
	@echo "  test-tagging-comprehensive - Run comprehensive tagging tests"
	@echo "  test-with-server          - Start server, run tests, stop server"
	@echo "  logs                      - Show server logs"
	@echo "  clean                     - Clean up test artifacts and stop server"
	@echo "  health-check              - Check if server is accessible"
	@echo ""
	@echo "Configuration:"
	@echo "  S3_PORT=${S3_PORT}"
	@echo "  TEST_TIMEOUT=${TEST_TIMEOUT}"

# Build the Hanzo binary
build-s3:
	@echo "Building Hanzo binary..."
	@cd ../../../s3 && go build -o s3_binary .
	@chmod +x $(WEED_BINARY)
	@echo "✅ Hanzo binary built at $(WEED_BINARY)"

check-deps: build-s3
	@echo "Checking dependencies..."
	@echo "🔍 DEBUG: Checking Go installation..."
	@command -v go >/dev/null 2>&1 || (echo "Go is required but not installed" && exit 1)
	@echo "🔍 DEBUG: Go version: $$(go version)"
	@echo "🔍 DEBUG: Checking binary at $(WEED_BINARY)..."
	@test -f $(WEED_BINARY) || (echo "Hanzo binary not found at $(WEED_BINARY)" && exit 1)
	@echo "🔍 DEBUG: Binary size: $$(ls -lh $(WEED_BINARY) | awk '{print $$5}')"
	@echo "🔍 DEBUG: Binary permissions: $$(ls -la $(WEED_BINARY) | awk '{print $$1}')"
	@echo "🔍 DEBUG: Checking Go module dependencies..."
	@go list -m github.com/aws/aws-sdk-go-v2 >/dev/null 2>&1 || (echo "AWS SDK Go v2 not found. Run 'go mod tidy'." && exit 1)
	@go list -m github.com/stretchr/testify >/dev/null 2>&1 || (echo "Testify not found. Run 'go mod tidy'." && exit 1)
	@echo "✅ All dependencies are available"

# Start Hanzo server for testing
start-server: check-deps
	@echo "Starting Hanzo server..."
	@echo "🔍 DEBUG: Current working directory: $$(pwd)"
	@echo "🔍 DEBUG: Checking for existing s3 processes..."
	@ps aux | grep s3 | grep -v grep || echo "No existing s3 processes found"
	@echo "🔍 DEBUG: Cleaning up any existing PID file..."
	@rm -f s3-server.pid
	@echo "🔍 DEBUG: Checking for port conflicts..."
	@if netstat -tlnp 2>/dev/null | grep $(S3_PORT) >/dev/null; then \
		echo "⚠️  Port $(S3_PORT) is already in use, trying to find the process..."; \
		netstat -tlnp 2>/dev/null | grep $(S3_PORT) || true; \
	else \
		echo "✅ Port $(S3_PORT) is available"; \
	fi
	@echo "🔍 DEBUG: Checking binary at $(WEED_BINARY)"
	@ls -la $(WEED_BINARY) || (echo "❌ Binary not found!" && exit 1)
	@echo "🔍 DEBUG: Creating volume directory..."
	@mkdir -p ./test-volume-data
	@echo "🔍 DEBUG: Launching Hanzo server in background..."
	@echo "🔍 DEBUG: Command: $(WEED_BINARY) mini -dir=./test-volume-data -s3.port=$(S3_PORT)"
	@$(WEED_BINARY) mini \
		-filer.maxMB=64 \
		-dir=./test-volume-data \
		-master.raftHashicorp \
		-s3.port=$(S3_PORT) \
		> s3-test.log 2>&1 & \
		echo $$! > s3-server.pid
	@echo "🔍 DEBUG: Waiting for server to start (up to 90 seconds)..."
	@for i in $$(seq 1 90); do \
		echo "🔍 DEBUG: Attempt $$i/90 - checking port $(S3_PORT)"; \
		if curl -s http://localhost:$(S3_PORT) >/dev/null 2>&1; then \
			echo "✅ Hanzo server started successfully on port $(S3_PORT) after $$i seconds"; \
			exit 0; \
		fi; \
		if [ $$i -eq 5 ]; then \
			echo "🔍 DEBUG: After 5 seconds, checking process and logs..."; \
			ps aux | grep s3 | grep -v grep || echo "No s3 processes found"; \
			if [ -f s3-test.log ]; then \
				echo "=== First server logs ==="; \
				head -20 s3-test.log; \
			fi; \
		fi; \
		if [ $$i -eq 15 ]; then \
			echo "🔍 DEBUG: After 15 seconds, checking port bindings..."; \
			netstat -tlnp 2>/dev/null | grep $(S3_PORT) || echo "Port $(S3_PORT) not bound"; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "⚠️  Server taking longer than expected (30s), checking logs..."; \
			if [ -f s3-test.log ]; then \
				echo "=== Recent server logs ==="; \
				tail -20 s3-test.log; \
			fi; \
		fi; \
		sleep 1; \
	done; \
	echo "❌ Server failed to start within 90 seconds"; \
	echo "🔍 DEBUG: Final process check:"; \
	ps aux | grep s3 | grep -v grep || echo "No s3 processes found"; \
	echo "🔍 DEBUG: Final port check:"; \
	netstat -tlnp 2>/dev/null | grep -E "($(S3_PORT))" || echo "No ports bound"; \
	echo "=== Full server logs ==="; \
	if [ -f s3-test.log ]; then \
		cat s3-test.log; \
	else \
		echo "No log file found"; \
	fi; \
	exit 1

# Stop Hanzo server
stop-server:
	@echo "Stopping Hanzo server..."
	@if [ -f s3-server.pid ]; then \
		SERVER_PID=$$(cat s3-server.pid); \
		echo "Killing server PID $$SERVER_PID"; \
		if ps -p $$SERVER_PID >/dev/null 2>&1; then \
			kill -TERM $$SERVER_PID 2>/dev/null || true; \
			sleep 2; \
			if ps -p $$SERVER_PID >/dev/null 2>&1; then \
				echo "Process still running, sending KILL signal..."; \
				kill -KILL $$SERVER_PID 2>/dev/null || true; \
				sleep 1; \
			fi; \
		else \
			echo "Process $$SERVER_PID not found (already stopped)"; \
		fi; \
		rm -f s3-server.pid; \
	else \
		echo "No PID file found, checking for running processes..."; \
		echo "⚠️  Skipping automatic process cleanup to avoid CI issues"; \
		echo "Note: Any remaining s3 processes should be cleaned up by the CI environment"; \
	fi
	@echo "✅ Hanzo server stopped"

# Show server logs
logs:
	@if test -f s3-test.log; then \
		echo "=== Hanzo Server Logs ==="; \
		tail -f s3-test.log; \
	else \
		echo "No log file found. Server may not be running."; \
	fi

# Core tagging tests (basic functionality)
test-tagging-quick: check-deps
	@echo "Running core tagging tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestObjectTaggingOnUpload|TestPutObjectTaggingAPI" .
	@echo "✅ Core tagging tests completed"

# All tagging tests (comprehensive)
test-tagging: check-deps
	@echo "Running all tagging tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" .
	@echo "✅ All tagging tests completed"

# Comprehensive tagging tests (all features)
test-tagging-comprehensive: check-deps
	@echo "Running comprehensive tagging tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) .
	@echo "✅ Comprehensive tagging tests completed"

# All tests without server management
test-tagging-simple: check-deps
	@echo "Running tagging tests (assuming server is already running)..."
	@go test -v -timeout=$(TEST_TIMEOUT) .
	@echo "✅ All tagging tests completed"

# Start server, run tests, stop server
test-with-server: start-server
	@echo "Running tagging tests with managed server..."
	@sleep 5  # Give server time to fully start
	@make test-tagging-comprehensive || (echo "Tests failed, stopping server..." && make stop-server && exit 1)
	@make stop-server
	@echo "✅ All tests completed with managed server"

# Health check
health-check:
	@echo "Checking server health..."
	@if curl -s http://localhost:$(S3_PORT) >/dev/null 2>&1; then \
		echo "✅ Server is accessible on port $(S3_PORT)"; \
	else \
		echo "❌ Server is not accessible on port $(S3_PORT)"; \
		exit 1; \
	fi

# Clean up
clean:
	@echo "Cleaning up test artifacts..."
	@make stop-server
	@rm -f s3-test.log
	@rm -f s3-server.pid
	@rm -rf ./test-volume-data
	@rm -f tagging.test
	@go clean -testcache
	@echo "✅ Cleanup completed"

# Individual test targets for specific functionality
test-upload:
	@echo "Running upload tagging tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestObjectTaggingOnUpload" .

test-special-chars:
	@echo "Running special characters tagging tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestObjectTaggingOnUploadWithSpecialCharacters" .

test-api:
	@echo "Running API tagging tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestPutObjectTaggingAPI" .

test-get:
	@echo "Running get tagging tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestGetObjectTaggingAPI" .

test-delete:
	@echo "Running delete tagging tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestDeleteObjectTaggingAPI" .

# Development targets
dev-start: start-server
	@echo "Development server started. Access S3 API at http://localhost:$(S3_PORT)"
	@echo "To stop: make stop-server"

dev-test: check-deps
	@echo "Running tests in development mode..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestObjectTaggingOnUpload" .

# CI targets
ci-test: check-deps
	@echo "Running tests in CI mode..."
	@go test -v -timeout=$(TEST_TIMEOUT) -race .

# All targets
test-all: test-tagging test-tagging-comprehensive
	@echo "✅ All tagging tests completed"

# Benchmark targets
benchmark-tagging:
	@echo "Running tagging performance benchmarks..."
	@go test -v -timeout=$(TEST_TIMEOUT) -bench=. -benchmem .

# Coverage targets
coverage:
	@echo "Running tests with coverage..."
	@go test -v -timeout=$(TEST_TIMEOUT) -coverprofile=coverage.out .
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format and lint
fmt:
	@echo "Formatting Go code..."
	@go fmt .

lint:
	@echo "Running linter..."
	@golint . || echo "golint not available, skipping..."

# Install dependencies for development
install-deps:
	@echo "Installing Go dependencies..."
	@go mod tidy
	@go mod download

# Show current configuration
show-config:
	@echo "Current configuration:"
	@echo "  WEED_BINARY: $(WEED_BINARY)"
	@echo "  S3_PORT: $(S3_PORT)"
	@echo "  TEST_TIMEOUT: $(TEST_TIMEOUT)"
	@echo "  TEST_PATTERN: $(TEST_PATTERN)"

# Legacy targets for backward compatibility
test: test-with-server
test-verbose: test-tagging-comprehensive
test-single: test-upload
test-clean: clean
build: check-deps
setup: check-deps
