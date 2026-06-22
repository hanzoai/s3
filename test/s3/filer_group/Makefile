# S3 Filer Group Test Makefile
# This Makefile provides targets for running S3 filer group integration tests
# These tests verify that S3 bucket operations work correctly when a filer group is configured.

.PHONY: help build-s3 check-deps start-server stop-server test test-with-server clean logs health-check

# Configuration
WEED_BINARY := ../../../s3/s3_binary
S3_PORT := 8333
MASTER_PORT := 9333
MASTER_GRPC_PORT := 19333
VOLUME_PORT := 8080
FILER_PORT := 8888
TEST_TIMEOUT := 10m
TEST_PATTERN := Test

# Filer group configuration
FILER_GROUP := testgroup
SERVER_DIR := ./test-volume-data
S3_CONFIG := ../../../docker/compose/s3.json

# Default target
help:
	@echo "S3 Filer Group Test Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  help                    - Show this help message"
	@echo "  build-s3             - Build the Hanzo binary"
	@echo "  check-deps             - Check dependencies and build binary if needed"
	@echo "  start-server           - Start Hanzo server with filer group"
	@echo "  stop-server            - Stop Hanzo server"
	@echo "  test                   - Run filer group tests (server must be running)"
	@echo "  test-with-server       - Start server, run tests, stop server"
	@echo "  logs                   - Show server logs"
	@echo "  clean                  - Clean up test artifacts and stop server"
	@echo "  health-check           - Check if server is accessible"
	@echo ""
	@echo "Configuration:"
	@echo "  S3_PORT=${S3_PORT}"
	@echo "  MASTER_GRPC_PORT=${MASTER_GRPC_PORT}"
	@echo "  FILER_GROUP=${FILER_GROUP}"
	@echo "  TEST_TIMEOUT=${TEST_TIMEOUT}"

# Build the Hanzo binary
build-s3:
	@echo "Building Hanzo binary..."
	@cd ../../../s3 && go build -o s3_binary .
	@chmod +x $(WEED_BINARY)
	@echo "OK Hanzo binary built at $(WEED_BINARY)"

check-deps: build-s3
	@echo "Checking dependencies..."
	@command -v go >/dev/null 2>&1 || (echo "Go is required but not installed" && exit 1)
	@echo "Go version: $$(go version)"
	@test -f $(WEED_BINARY) || (echo "Hanzo binary not found at $(WEED_BINARY)" && exit 1)
	@go list -m github.com/aws/aws-sdk-go-v2 >/dev/null 2>&1 || (echo "AWS SDK Go v2 not found. Run 'go mod tidy'." && exit 1)
	@go list -m github.com/stretchr/testify >/dev/null 2>&1 || (echo "Testify not found. Run 'go mod tidy'." && exit 1)
	@echo "OK All dependencies are available"

# Start Hanzo server with filer group configured
start-server: check-deps
	@echo "Starting Hanzo server with filer group: $(FILER_GROUP)..."
	@rm -f s3-server.pid
	@mkdir -p $(SERVER_DIR)
	@if netstat -tlnp 2>/dev/null | grep $(S3_PORT) >/dev/null; then \
		echo "WARNING: Port $(S3_PORT) is already in use"; \
		exit 1; \
	fi
	@echo "Launching Hanzo server with filer group $(FILER_GROUP)..."
	@export AWS_ACCESS_KEY_ID=some_access_key1 && \
	export AWS_SECRET_ACCESS_KEY=some_secret_key1 && \
	$(WEED_BINARY) mini \
		-debug \
		-dir=$(SERVER_DIR) \
		-s3.port=$(S3_PORT) \
		-s3.config=$(S3_CONFIG) \
		-filer.filerGroup=$(FILER_GROUP) \
		> s3-server.log 2>&1 & \
		echo $$! > s3-server.pid
	
	@echo "Waiting for S3 server to be ready..."
	@for i in $$(seq 1 30); do \
		if echo | nc -z localhost $(S3_PORT); then \
			echo "S3 server is ready!"; \
			exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "❌ Server failed to start within 30 seconds"; \
	if [ -f s3-server.log ]; then \
		cat s3-server.log; \
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
		echo "No PID file found"; \
	fi
	@echo "✅ Hanzo server stopped"

# Show server logs
logs:
	@if test -f s3-server.log; then \
		echo "=== Hanzo Server Logs ==="; \
		tail -f s3-server.log; \
	else \
		echo "No log file found. Server may not be running."; \
	fi

# Run filer group tests (assumes server is already running)
test: check-deps
	@echo "Running filer group S3 tests..."
	@FILER_GROUP=$(FILER_GROUP) S3_ENDPOINT=http://localhost:$(S3_PORT) MASTER_ADDRESS=localhost:$(MASTER_GRPC_PORT) \
		go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" .
	@echo "✅ Filer group tests completed"

# Run tests with automatic server management
test-with-server: start-server
	@echo "Server started successfully, running filer group tests..."
	@echo "Test pattern: $(TEST_PATTERN)"
	@echo "Test timeout: $(TEST_TIMEOUT)"
	@trap "$(MAKE) stop-server" EXIT; \
	$(MAKE) test || (echo "❌ Tests failed, showing server logs:" && echo "=== Last 50 lines of server logs ===" && tail -50 s3-server.log && echo "=== End of server logs ===" && exit 1)
	@$(MAKE) stop-server
	@echo "✅ Tests completed and server stopped"

# Clean up test artifacts
clean:
	@echo "Cleaning up test artifacts..."
	@$(MAKE) stop-server
	@rm -f s3-server.log s3-test*.log s3-server.pid
	@rm -rf test-volume-data/
	@go clean -testcache
	@echo "✅ Cleanup completed"

# Quick health check
health-check:
	@echo "Running health check..."
	@curl -s http://localhost:$(S3_PORT) >/dev/null 2>&1 && echo "✅ S3 API is accessible" || echo "❌ S3 API is not accessible"
	@curl -s http://localhost:9325/metrics >/dev/null 2>&1 && echo "✅ Metrics endpoint is accessible" || echo "❌ Metrics endpoint is not accessible"
