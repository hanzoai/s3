# S3 API Test Makefile
# This Makefile provides comprehensive targets for running S3 versioning tests

.PHONY: all test help build-s3 setup-server start-server stop-server test-versioning test-versioning-quick test-versioning-regressions test-versioning-comprehensive test-all clean logs check-deps

# Configuration
WEED_BINARY := ../../../s3/s3_binary
S3_PORT := 8333
MASTER_PORT := 9333
VOLUME_PORT := 8080
FILER_PORT := 8888
TEST_TIMEOUT := 10m
# Run every test in the suite by default. "TestVersioning" used to leave whole
# tests (bucket creation, suspended delete, directory/version listing) ungated.
# Opt-in stress tests are matched by "." but self-skip without ENABLE_STRESS_TESTS
# (they run with it set in their own targets), so they don't actually execute here.
TEST_PATTERN := .

.DEFAULT_GOAL := help

all: test

test: test-with-server

# Default target
help:
	@echo "S3 API Test Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  help                    - Show this help message"
	@echo "  build-s3             - Build the Hanzo binary"
	@echo "  check-deps             - Check dependencies and build binary if needed"
	@echo "  start-server           - Start Hanzo server for testing"
	@echo "  start-server-simple    - Start server without process cleanup (for CI)"
	@echo "  stop-server            - Stop Hanzo server"
	@echo "  test-versioning        - Run all versioning tests"
	@echo "  test-versioning-quick  - Run core versioning tests only"
	@echo "  test-versioning-regressions - Run targeted versioning regression tests"
	@echo "  test-versioning-simple - Run tests without server management"
	@echo "  test-versioning-comprehensive - Run comprehensive versioning tests"
	@echo "  test-all               - Run all S3 API tests"
	@echo "  test-with-server       - Start server, run tests, stop server"
	@echo "  logs                   - Show server logs"
	@echo "  clean                  - Clean up test artifacts and stop server"
	@echo "  health-check           - Check if server is accessible"
	@echo ""
	@echo "Configuration:"
	@echo "  S3_PORT=${S3_PORT}"
	@echo "  TEST_TIMEOUT=${TEST_TIMEOUT}"

# Check dependencies
# Build the Hanzo binary
build-s3:
	@echo "Building Hanzo binary..."
	@cd ../../../s3 && go build -o s3_binary .
	@chmod +x $(WEED_BINARY)
	@echo "OK Hanzo binary built at $(WEED_BINARY)"

check-deps: build-s3
	@echo "Checking dependencies..."
	@echo "DEBUG: Checking Go installation..."
	@command -v go >/dev/null 2>&1 || (echo "Go is required but not installed" && exit 1)
	@echo "DEBUG: Go version: $$(go version)"
	@echo "🔍 DEBUG: Checking binary at $(WEED_BINARY)..."
	@test -f $(WEED_BINARY) || (echo "Hanzo binary not found at $(WEED_BINARY)" && exit 1)
	@echo "🔍 DEBUG: Binary size: $$(ls -lh $(WEED_BINARY) | awk '{print $$5}')"
	@echo "🔍 DEBUG: Binary permissions: $$(ls -la $(WEED_BINARY) | awk '{print $$1}')"
	@echo "🔍 DEBUG: Checking Go module dependencies..."
	@go list -m github.com/aws/aws-sdk-go-v2 >/dev/null 2>&1 || (echo "AWS SDK Go v2 not found. Run 'go mod tidy'." && exit 1)
	@go list -m github.com/stretchr/testify >/dev/null 2>&1 || (echo "Testify not found. Run 'go mod tidy'." && exit 1)
	@echo "OK All dependencies are available"

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
	@echo "🔍 DEBUG: Checking config file at ../../../docker/compose/s3.json"
	@ls -la ../../../docker/compose/s3.json || echo "⚠️  Config file not found, continuing without it"
	@echo "🔍 DEBUG: Creating volume directory..."
	@mkdir -p ./test-volume-data
	@echo "🔍 DEBUG: Launching Hanzo server in background..."
	@echo "🔍 DEBUG: Command: $(WEED_BINARY) mini -dir=./test-volume-data -s3.port=$(S3_PORT) -s3.config=../../../docker/compose/s3.json"
	@$(WEED_BINARY) mini \
		-dir=./test-volume-data \
		-s3.port=$(S3_PORT) \
		-s3.config=../../../docker/compose/s3.json \
		> s3-test.log 2>&1 & echo $$! > s3-server.pid
	@echo "🔍 DEBUG: Server PID: $$(cat s3-server.pid 2>/dev/null || echo 'PID file not found')"
	@echo "🔍 DEBUG: Checking if PID is still running..."
	@sleep 2
	@if [ -f s3-server.pid ]; then \
		SERVER_PID=$$(cat s3-server.pid); \
		ps -p $$SERVER_PID || echo "⚠️  Server PID $$SERVER_PID not found after 2 seconds"; \
	else \
		echo "⚠️  PID file not found"; \
	fi
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
			netstat -tlnp 2>/dev/null | grep 9333 || echo "Port 9333 not bound"; \
			netstat -tlnp 2>/dev/null | grep 8080 || echo "Port 8080 not bound"; \
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
	netstat -tlnp 2>/dev/null | grep -E "(8333|9333|8080)" || echo "No ports bound"; \
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

# Core versioning tests (equivalent to Python s3tests)
test-versioning-quick: check-deps
	@echo "Running core S3 versioning tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestBucketListReturnDataVersioning|TestVersioningBasicWorkflow|TestVersioningDeleteMarkers" .
	@echo "✅ Core versioning tests completed"

# Targeted regressions for conditional mutation, multipart, copy, and suspended delete flows.
test-versioning-regressions: check-deps
	@echo "Running targeted S3 versioning regression tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "TestVersioningCompleteMultipartUploadIsIdempotent|TestVersioningSelfCopyMetadataReplaceCreatesNewVersion|TestVersioningSelfCopyMetadataReplaceSuspendedKeepsNullVersion|TestSuspendedDeleteCreatesDeleteMarker" .
	@echo "✅ Targeted versioning regression tests completed"

# All versioning tests
test-versioning: check-deps
	@echo "Running all S3 versioning tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" .
	@echo "✅ All versioning tests completed"

# Comprehensive versioning tests (including edge cases)
test-versioning-comprehensive: check-deps
	@echo "Running comprehensive S3 versioning tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" . -count=1
	@echo "✅ Comprehensive versioning tests completed"

# All S3 API tests
test-all: check-deps
	@echo "Running all S3 API tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) ./...
	@echo "✅ All S3 API tests completed"

# Run tests with automatic server management
test-with-server: start-server
	@echo "🔍 DEBUG: Server started successfully, now running versioning tests..."
	@echo "🔍 DEBUG: Test pattern: $(TEST_PATTERN)"
	@echo "🔍 DEBUG: Test timeout: $(TEST_TIMEOUT)"
	@echo "Running versioning tests with managed server..."
	@trap "$(MAKE) stop-server" EXIT; \
	$(MAKE) test-versioning || (echo "❌ Tests failed, showing server logs:" && echo "=== Last 50 lines of server logs ===" && tail -50 s3-test.log && echo "=== End of server logs ===" && exit 1)
	@$(MAKE) stop-server
	@echo "✅ Tests completed and server stopped"

# Test with different configurations
test-versioning-with-configs: check-deps
	@echo "Testing with different S3 configurations..."
	@echo "Testing with empty folder allowed..."
	@$(WEED_BINARY) mini -s3.port=$(S3_PORT) > s3-test-config1.log 2>&1 & echo $$! > s3-config1.pid
	@sleep 5
	@go test -v -timeout=5m -run "TestVersioningBasicWorkflow" . || true
	@if [ -f s3-config1.pid ]; then kill -TERM $$(cat s3-config1.pid) 2>/dev/null || true; rm -f s3-config1.pid; fi
	@sleep 2
	@echo "Testing with delete bucket not empty disabled..."
	@$(WEED_BINARY) mini -s3.port=$(S3_PORT) -s3.allowDeleteBucketNotEmpty=false > s3-test-config2.log 2>&1 & echo $$! > s3-config2.pid
	@sleep 5
	@go test -v -timeout=5m -run "TestVersioningBasicWorkflow" . || true
	@if [ -f s3-config2.pid ]; then kill -TERM $$(cat s3-config2.pid) 2>/dev/null || true; rm -f s3-config2.pid; fi
	@echo "✅ Configuration tests completed"

# Performance/stress testing
test-versioning-stress: check-deps
	@echo "Running stress tests for versioning..."
	@go test -v -timeout=20m -run "TestVersioningConcurrentOperations" . -count=5
	@echo "✅ Stress tests completed"

# Pagination stress testing (tests >1000 versions)
test-versioning-pagination-stress: check-deps
	@echo "Running pagination stress tests (>1000 versions)..."
	@echo "This test creates 1500+ versions and may take several minutes..."
	@ENABLE_STRESS_TESTS=true go test -v -timeout=30m -run "TestVersioningPagination" .
	@echo "✅ Pagination stress tests completed"

# Generate test reports
test-report: check-deps
	@echo "Generating test reports..."
	@mkdir -p reports
	@go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" . -json > reports/test-results.json 2>&1 || true
	@go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" . -coverprofile=reports/coverage.out 2>&1 || true
	@go tool cover -html=reports/coverage.out -o reports/coverage.html 2>/dev/null || true
	@echo "✅ Test reports generated in reports/ directory"

# Clean up test artifacts
clean:
	@echo "Cleaning up test artifacts..."
	@$(MAKE) stop-server
	@rm -f s3-test*.log s3-server.pid s3-config*.pid
	@rm -rf reports/
	@rm -rf test-volume-data/
	@go clean -testcache
	@echo "✅ Cleanup completed"

# Debug mode - start server with verbose logging
debug-server:
	@echo "Starting Hanzo server in debug mode..."
	@$(MAKE) stop-server
	@mkdir -p ./test-volume-data
	@$(WEED_BINARY) mini \
		-debug \
		-s3.port=$(S3_PORT) \
		-s3.allowDeleteBucketNotEmpty=true \
		-s3.config=../../../docker/compose/s3.json \
		-filer.maxMB=16 \
		-master.peers=none \
		-dir=./test-volume-data \
		-metricsPort=9324

# Run a single test for debugging
debug-test: check-deps
	@echo "Running single test for debugging..."
	@go test -v -timeout=5m -run "TestBucketListReturnDataVersioning" . -count=1

# Continuous testing (re-run tests on file changes)
watch-tests:
	@echo "Starting continuous testing (requires 'entr' command)..."
	@command -v entr >/dev/null 2>&1 || (echo "Install 'entr' for file watching: brew install entr (macOS) or apt-get install entr (Linux)" && exit 1)
	@find . -name "*.go" | entr -c $(MAKE) test-versioning-quick

# Install missing Go dependencies
install-deps:
	@echo "Installing Go dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Validate test configuration
validate-config:
	@echo "Validating test configuration..."
	@test -f test_config.json || (echo "❌ test_config.json not found" && exit 1)
	@python3 -m json.tool test_config.json > /dev/null 2>&1 || (echo "❌ test_config.json is not valid JSON" && exit 1)
	@echo "✅ Configuration is valid"

# Quick health check
health-check:
	@echo "Running health check..."
	@curl -s http://localhost:$(S3_PORT) >/dev/null 2>&1 && echo "✅ S3 API is accessible" || echo "❌ S3 API is not accessible"
	@curl -s http://localhost:9324/metrics >/dev/null 2>&1 && echo "✅ Metrics endpoint is accessible" || echo "❌ Metrics endpoint is not accessible"

# Simple server start without process cleanup (for CI troubleshooting)
start-server-simple: check-deps
	@echo "Starting Hanzo server (simple mode)..."
	@$(WEED_BINARY) mini \
		-dir=$(SERVER_DIR) \
		-s3 \
		-s3.port=$(S3_PORT) \
		-s3.config=$(S3_CONFIG) \
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
	echo "S3 server failed to start"; \
	exit 1

# Simple test run without server management
test-versioning-simple: check-deps
	@echo "Running versioning tests (assuming server is already running)..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" .
	@echo "✅ Tests completed"

# Force cleanup all s3 processes (use with caution)
force-cleanup:
	@echo "⚠️  Force cleaning up all s3 processes..."
	@echo "This will attempt to kill ALL s3 processes on the system"
	@ps aux | grep s3 | grep -v grep || echo "No s3 processes found"
	@killall -TERM s3_binary 2>/dev/null || echo "No s3_binary processes to terminate"
	@sleep 2
	@killall -KILL s3_binary 2>/dev/null || echo "No s3_binary processes to kill"
	@rm -f s3-server.pid s3-config*.pid
	@echo "✅ Force cleanup completed"

# Compare with Python s3tests (if available)
compare-python-tests:
	@echo "Comparing Go tests with Python s3tests..."
	@echo "Go test: TestBucketListReturnDataVersioning"
	@echo "Python equivalent: test_bucket_list_return_data_versioning"
	@echo ""
	@echo "Running Go version..."
	@time go test -v -run "TestBucketListReturnDataVersioning" . 2>&1 | grep -E "(PASS|FAIL|took)" 
