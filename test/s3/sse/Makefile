# Makefile for S3 SSE Integration Tests
# This Makefile provides targets for running comprehensive S3 Server-Side Encryption tests

# Default values
SEAWEEDFS_BINARY ?= s3
S3_PORT ?= 8333
FILER_PORT ?= 8888
VOLUME_PORT ?= 8080
MASTER_PORT ?= 9333
TEST_TIMEOUT ?= 15m
BUCKET_PREFIX ?= test-sse-
ACCESS_KEY ?= some_access_key1
SECRET_KEY ?= some_secret_key1
VOLUME_MAX_SIZE_MB ?= 50
VOLUME_MAX_COUNT ?= 100

# SSE-KMS configuration
KMS_KEY_ID ?= test-key-123
KMS_TYPE ?= local

# Test directory
TEST_DIR := $(shell pwd)
SEAWEEDFS_ROOT := $(shell cd ../../../ && pwd)

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m # No Color

.PHONY: all test clean start-hanzo stop-hanzo stop-hanzo-safe start-hanzo-ci check-binary build-s3 help help-extended test-with-server test-quick-with-server test-metadata-persistence

all: test-basic

# Build Hanzo binary (GitHub Actions compatible)
build-s3:
	@echo "Building Hanzo binary..."
	@cd $(SEAWEEDFS_ROOT)/s3 && go install -buildvcs=false
	@echo "✅ Hanzo binary built successfully"

help:
	@echo "Hanzo S3 SSE Integration Tests"
	@echo ""
	@echo "Available targets:"
	@echo "  test-basic    - Run basic S3 put/get tests first"
	@echo "  test          - Run all S3 SSE integration tests"
	@echo "  test-ssec     - Run SSE-C tests only"
	@echo "  test-ssekms   - Run SSE-KMS tests only"
	@echo "  test-copy     - Run SSE copy operation tests"
	@echo "  test-multipart - Run SSE multipart upload tests"
	@echo "  test-errors   - Run SSE error condition tests"
	@echo "  benchmark     - Run SSE performance benchmarks"
	@echo "  start-hanzo - Start Hanzo server for testing"
	@echo "  stop-hanzo - Stop Hanzo server"
	@echo "  clean         - Clean up test artifacts"
	@echo "  check-binary  - Check if Hanzo binary exists"
	@echo ""
	@echo "Configuration:"
	@echo "  SEAWEEDFS_BINARY=$(SEAWEEDFS_BINARY)"
	@echo "  S3_PORT=$(S3_PORT)"
	@echo "  FILER_PORT=$(FILER_PORT)"
	@echo "  VOLUME_PORT=$(VOLUME_PORT)"
	@echo "  MASTER_PORT=$(MASTER_PORT)"
	@echo "  TEST_TIMEOUT=$(TEST_TIMEOUT)"
	@echo "  VOLUME_MAX_SIZE_MB=$(VOLUME_MAX_SIZE_MB)"

check-binary:
	@if ! command -v $(SEAWEEDFS_BINARY) > /dev/null 2>&1; then \
		echo "$(RED)Error: Hanzo binary '$(SEAWEEDFS_BINARY)' not found in PATH$(NC)"; \
		echo "Please build Hanzo first by running 'make' in the root directory"; \
		exit 1; \
	fi
	@echo "$(GREEN)Hanzo binary found: $$(which $(SEAWEEDFS_BINARY))$(NC)"

start-hanzo: check-binary
	@echo "$(YELLOW)Starting Hanzo server for SSE testing...$(NC)"
	@# Use port-based cleanup for consistency and safety
	@echo "Cleaning up any existing processes..."
	@lsof -ti :$(MASTER_PORT) | xargs -r kill -TERM || true
	@lsof -ti :$(VOLUME_PORT) | xargs -r kill -TERM || true
	@lsof -ti :$(FILER_PORT) | xargs -r kill -TERM || true
	@lsof -ti :$(S3_PORT) | xargs -r kill -TERM || true
	@sleep 2
	
	# Create necessary directories
	@mkdir -p /tmp/hanzo-test-sse
	
	# Create S3 configuration with SSE-KMS support
	@printf '{"identities":[{"name":"%s","credentials":[{"accessKey":"%s","secretKey":"%s"}],"actions":["Admin","Read","Write"]}],"kms":{"type":"%s","configs":{"keyId":"%s","encryptionContext":{},"bucketKey":false}}}' "$(ACCESS_KEY)" "$(ACCESS_KEY)" "$(SECRET_KEY)" "$(KMS_TYPE)" "$(KMS_KEY_ID)" > /tmp/hanzo-sse-s3.json
	
	# Start s3 mini (WEED_S3_SSE_KEY enables SSE-S3 encryption)
	@AWS_ACCESS_KEY_ID=$(ACCESS_KEY) AWS_SECRET_ACCESS_KEY=$(SECRET_KEY) WEED_S3_SSE_KEY=test-sse-s3-key $(SEAWEEDFS_BINARY) mini \
		-dir=/tmp/hanzo-test-sse \
		-s3.port=$(S3_PORT) \
		-s3.config=/tmp/hanzo-sse-s3.json \
		> /tmp/hanzo-sse-mini.log 2>&1 & echo $$! > /tmp/s3-mini.pid

	@echo "Checking S3 service is ready..."
	@for i in $$(seq 1 30); do \
		if curl -s http://127.0.0.1:$(S3_PORT) > /dev/null 2>&1; then \
			echo "✅ S3 service is ready"; \
			break; \
		fi; \
		sleep 1; \
	done

stop-hanzo:
	@echo "$(YELLOW)Stopping Hanzo server...$(NC)"
	@# Use port-based cleanup for consistency and safety
	@if [ -f /tmp/s3-mini.pid ]; then \
		echo "Stopping s3 mini..."; \
		kill $$(cat /tmp/s3-mini.pid) || true; \
		rm -f /tmp/s3-mini.pid; \
	fi
	@lsof -ti :$(MASTER_PORT) | xargs -r kill -TERM || true
	@lsof -ti :$(VOLUME_PORT) | xargs -r kill -TERM || true
	@lsof -ti :$(FILER_PORT) | xargs -r kill -TERM || true
	@lsof -ti :$(S3_PORT) | xargs -r kill -TERM || true
	@sleep 2
	@echo "$(GREEN)Hanzo server stopped$(NC)"

# CI-safe server stop that's more conservative
stop-hanzo-safe:
	@echo "$(YELLOW)Safely stopping Hanzo server...$(NC)"
	@# Use port-based cleanup which is safer in CI
	@if command -v lsof >/dev/null 2>&1; then \
		echo "Using lsof for port-based cleanup..."; \
		lsof -ti :$(MASTER_PORT) 2>/dev/null | head -5 | while read pid; do kill -TERM $$pid 2>/dev/null || true; done; \
		lsof -ti :$(VOLUME_PORT) 2>/dev/null | head -5 | while read pid; do kill -TERM $$pid 2>/dev/null || true; done; \
		lsof -ti :$(FILER_PORT) 2>/dev/null | head -5 | while read pid; do kill -TERM $$pid 2>/dev/null || true; done; \
		lsof -ti :$(S3_PORT) 2>/dev/null | head -5 | while read pid; do kill -TERM $$pid 2>/dev/null || true; done; \
	else \
		echo "lsof not available, using netstat approach..."; \
		netstat -tlnp 2>/dev/null | grep :$(MASTER_PORT) | awk '{print $$7}' | cut -d/ -f1 | head -5 | while read pid; do [ "$$pid" != "-" ] && kill -TERM $$pid 2>/dev/null || true; done; \
		netstat -tlnp 2>/dev/null | grep :$(VOLUME_PORT) | awk '{print $$7}' | cut -d/ -f1 | head -5 | while read pid; do [ "$$pid" != "-" ] && kill -TERM $$pid 2>/dev/null || true; done; \
		netstat -tlnp 2>/dev/null | grep :$(FILER_PORT) | awk '{print $$7}' | cut -d/ -f1 | head -5 | while read pid; do [ "$$pid" != "-" ] && kill -TERM $$pid 2>/dev/null || true; done; \
		netstat -tlnp 2>/dev/null | grep :$(S3_PORT) | awk '{print $$7}' | cut -d/ -f1 | head -5 | while read pid; do [ "$$pid" != "-" ] && kill -TERM $$pid 2>/dev/null || true; done; \
	fi
	@sleep 2
	@echo "$(GREEN)Hanzo server safely stopped$(NC)"

clean:
	@echo "$(YELLOW)Cleaning up SSE test artifacts...$(NC)"
	@rm -rf /tmp/hanzo-test-sse-*
	@rm -f /tmp/hanzo-sse-*.log
	@rm -f /tmp/hanzo-sse-s3.json
	@echo "$(GREEN)SSE test cleanup completed$(NC)"

test-basic: check-binary
	@echo "$(YELLOW)Running basic S3 SSE integration tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting basic SSE tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSECIntegrationBasic|TestSSEKMSIntegrationBasic" ./test/s3/sse || (echo "$(RED)Basic SSE tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)Basic SSE tests completed successfully!$(NC)"

test: test-basic
	@echo "$(YELLOW)Running all S3 SSE integration tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting comprehensive SSE tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSE.*Integration" ./test/s3/sse || (echo "$(RED)SSE tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)All SSE integration tests completed successfully!$(NC)"

test-ssec: check-binary
	@echo "$(YELLOW)Running SSE-C integration tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting SSE-C tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSEC.*Integration" ./test/s3/sse || (echo "$(RED)SSE-C tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)SSE-C tests completed successfully!$(NC)"

test-ssekms: check-binary
	@echo "$(YELLOW)Running SSE-KMS integration tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting SSE-KMS tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSEKMS.*Integration" ./test/s3/sse || (echo "$(RED)SSE-KMS tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)SSE-KMS tests completed successfully!$(NC)"

test-copy: check-binary
	@echo "$(YELLOW)Running SSE copy operation tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting SSE copy tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run ".*CopyIntegration" ./test/s3/sse || (echo "$(RED)SSE copy tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)SSE copy tests completed successfully!$(NC)"

test-multipart: check-binary
	@echo "$(YELLOW)Running SSE multipart upload tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting SSE multipart tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSEMultipartUploadIntegration" ./test/s3/sse || (echo "$(RED)SSE multipart tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)SSE multipart tests completed successfully!$(NC)"

test-errors: check-binary
	@echo "$(YELLOW)Running SSE error condition tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting SSE error tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSEErrorConditions" ./test/s3/sse || (echo "$(RED)SSE error tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)SSE error tests completed successfully!$(NC)"

test-quick: check-binary
	@echo "$(YELLOW)Running quick SSE tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting quick SSE tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=5m -run "TestSSECIntegrationBasic|TestSSEKMSIntegrationBasic" ./test/s3/sse || (echo "$(RED)Quick SSE tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)Quick SSE tests completed successfully!$(NC)"

benchmark: check-binary
	@echo "$(YELLOW)Running SSE performance benchmarks...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Starting SSE benchmarks...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=30m -bench=. -run=Benchmark ./test/s3/sse || (echo "$(RED)SSE benchmarks failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)SSE benchmarks completed!$(NC)"

# Debug targets
debug-logs:
	@echo "$(YELLOW)=== Master Log ===$(NC)"
	@tail -n 50 /tmp/hanzo-sse-master.log || echo "No master log found"
	@echo "$(YELLOW)=== Volume Log ===$(NC)"
	@tail -n 50 /tmp/hanzo-sse-volume.log || echo "No volume log found"
	@echo "$(YELLOW)=== Filer Log ===$(NC)"
	@tail -n 50 /tmp/hanzo-sse-filer.log || echo "No filer log found"
	@echo "$(YELLOW)=== S3 Log ===$(NC)"
	@tail -n 50 /tmp/hanzo-sse-s3.log || echo "No S3 log found"

debug-status:
	@echo "$(YELLOW)=== Process Status ===$(NC)"
	@ps aux | grep -E "(s3|hanzo)" | grep -v grep || echo "No Hanzo processes found"
	@echo "$(YELLOW)=== Port Status ===$(NC)"
	@netstat -an | grep -E "($(MASTER_PORT)|$(VOLUME_PORT)|$(FILER_PORT)|$(S3_PORT))" || echo "No ports in use"

# Manual test targets for development
manual-start: start-hanzo
	@echo "$(GREEN)Hanzo with SSE support is now running for manual testing$(NC)"
	@echo "You can now run SSE tests manually or use S3 clients to test SSE functionality"
	@echo "Run 'make manual-stop' when finished"

manual-stop: stop-hanzo clean

# CI/CD targets
ci-test: test-quick

# Stress test
stress: check-binary
	@echo "$(YELLOW)Running SSE stress tests...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=60m -run="TestSSE.*Integration" -count=5 ./test/s3/sse || (echo "$(RED)SSE stress tests failed$(NC)" && $(MAKE) stop-hanzo-safe && exit 1)
	@$(MAKE) stop-hanzo-safe
	@echo "$(GREEN)SSE stress tests completed!$(NC)"

# Performance test with various data sizes
perf: check-binary
	@echo "$(YELLOW)Running SSE performance tests with various data sizes...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=60m -run=".*VariousDataSizes" ./test/s3/sse || (echo "$(RED)SSE performance tests failed$(NC)" && $(MAKE) -C $(TEST_DIR) stop-hanzo-safe && exit 1)
	@$(MAKE) -C $(TEST_DIR) stop-hanzo-safe
	@echo "$(GREEN)SSE performance tests completed!$(NC)"

# Test specific scenarios that would catch the metadata bug
test-metadata-persistence: check-binary
	@echo "$(YELLOW)Running SSE metadata persistence tests (would catch filer metadata bugs)...$(NC)"
	@$(MAKE) start-hanzo-ci
	@sleep 5
	@echo "$(GREEN)Testing that SSE metadata survives full PUT/GET cycle...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSECIntegrationBasic" ./test/s3/sse || (echo "$(RED)SSE metadata persistence tests failed$(NC)" && $(MAKE) -C $(TEST_DIR) stop-hanzo-safe && exit 1)
	@$(MAKE) -C $(TEST_DIR) stop-hanzo-safe
	@echo "$(GREEN)SSE metadata persistence tests completed successfully!$(NC)"
	@echo "$(GREEN)✅ These tests would have caught the filer metadata storage bug!$(NC)"

# GitHub Actions compatible test-with-server target that handles server lifecycle
test-with-server: build-s3
	@echo "🚀 Starting SSE integration tests with automated server management..."
	@echo "Starting Hanzo cluster..."
	@# Use the CI-safe startup directly without aggressive cleanup
	@if $(MAKE) start-hanzo-ci > s3-test.log 2>&1; then \
		echo "✅ Hanzo cluster started successfully"; \
		echo "Running SSE integration tests..."; \
		trap '$(MAKE) -C $(TEST_DIR) stop-hanzo-safe || true' EXIT; \
		if [ -n "$(TEST_PATTERN)" ]; then \
			echo "🔍 Running tests matching pattern: $(TEST_PATTERN)"; \
			cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" ./test/s3/sse || exit 1; \
		else \
			echo "🔍 Running all SSE integration tests"; \
			cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSE.*Integration" ./test/s3/sse || exit 1; \
		fi; \
		echo "✅ All tests completed successfully"; \
		$(MAKE) -C $(TEST_DIR) stop-hanzo-safe || true; \
	else \
		echo "❌ Failed to start Hanzo cluster"; \
		echo "=== Server startup logs ==="; \
		tail -100 s3-test.log 2>/dev/null || echo "No startup log available"; \
		echo "=== System information ==="; \
		ps aux | grep -E "s3|make" | grep -v grep || echo "No relevant processes found"; \
		exit 1; \
	fi

# CI-safe server startup that avoids process conflicts
start-hanzo-ci: check-binary
	@echo "$(YELLOW)Starting Hanzo server for CI testing...$(NC)"
	
	# Create necessary directories
	@mkdir -p /tmp/hanzo-test-sse
	
	# Clean up any old server logs
	@rm -f /tmp/hanzo-sse-*.log || true
	
	# Create S3 JSON configuration with KMS (Local provider) and basic identity for embedded S3
	@sed -e 's/ACCESS_KEY_PLACEHOLDER/$(ACCESS_KEY)/g' \
	     -e 's/SECRET_KEY_PLACEHOLDER/$(SECRET_KEY)/g' \
	     s3-config-template.json > /tmp/hanzo-s3.json
	
	# Start s3 mini with embedded S3 using the JSON config (with verbose logging)
	# WEED_S3_SSE_KEY enables SSE-S3 encryption for testing (KEK derived via HKDF)
	@echo "Starting s3 mini with embedded S3..."
	@AWS_ACCESS_KEY_ID=$(ACCESS_KEY) AWS_SECRET_ACCESS_KEY=$(SECRET_KEY) WEED_S3_SSE_KEY=test-sse-s3-key GLOG_v=4 $(SEAWEEDFS_BINARY) mini \
		-dir=/tmp/hanzo-test-sse \
		-s3.port=$(S3_PORT) \
		-s3.config=/tmp/hanzo-s3.json \
		-ip=127.0.0.1 \
		> /tmp/hanzo-sse-mini.log 2>&1 & echo $$! > /tmp/s3-mini.pid
	
	@echo "Checking S3 service is ready..."
	@for i in $$(seq 1 30); do \
		if curl -s http://127.0.0.1:$(S3_PORT) > /dev/null 2>&1; then \
			echo "✅ S3 service is ready"; \
			break; \
		fi; \
		sleep 1; \
	done

# GitHub Actions compatible quick test subset
test-quick-with-server: build-s3
	@echo "🚀 Starting quick SSE tests with automated server management..."
	@trap 'make stop-hanzo-safe || true' EXIT; \
	echo "Starting Hanzo cluster..."; \
	if make start-hanzo-ci > s3-test.log 2>&1; then \
		echo "✅ Hanzo cluster started successfully"; \
		echo "Running quick SSE integration tests..."; \
		cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestSSECIntegrationBasic|TestSSEKMSIntegrationBasic|TestSimpleSSECIntegration" ./test/s3/sse || exit 1; \
		echo "✅ Quick tests completed successfully"; \
		make stop-hanzo-safe || true; \
	else \
		echo "❌ Failed to start Hanzo cluster"; \
		echo "=== Server startup logs ==="; \
		tail -50 s3-test.log; \
		exit 1; \
	fi

# Help target - extended version
help-extended:
	@echo "Available targets:"
	@echo "  test                    - Run all SSE integration tests (requires running server)"
	@echo "  test-with-server        - Run all tests with automatic server management (GitHub Actions compatible)"
	@echo "  test-quick-with-server  - Run quick tests with automatic server management"
	@echo "  test-ssec               - Run only SSE-C tests"  
	@echo "  test-ssekms             - Run only SSE-KMS tests"
	@echo "  test-copy               - Run only copy operation tests"
	@echo "  test-multipart          - Run only multipart upload tests"
	@echo "  benchmark               - Run performance benchmarks"
	@echo "  perf                    - Run performance tests with various data sizes"
	@echo "  test-metadata-persistence - Test metadata persistence (catches filer bugs)"
	@echo "  build-s3              - Build Hanzo binary"
	@echo "  check-binary            - Check if Hanzo binary exists"
	@echo "  start-hanzo         - Start Hanzo cluster"
	@echo "  start-hanzo-ci      - Start Hanzo cluster (CI-safe version)"
	@echo "  stop-hanzo          - Stop Hanzo cluster"
	@echo "  stop-hanzo-safe     - Stop Hanzo cluster (CI-safe version)"
	@echo "  clean                   - Clean up test artifacts"
	@echo "  debug-logs              - Show recent logs from all services"
	@echo ""
	@echo "Environment Variables:"
	@echo "  ACCESS_KEY              - S3 access key (default: some_access_key1)"
	@echo "  SECRET_KEY              - S3 secret key (default: some_secret_key1)"
	@echo "  KMS_KEY_ID              - KMS key ID for SSE-KMS (default: test-key-123)"
	@echo "  KMS_TYPE                - KMS type (default: local)"
	@echo "  VOLUME_MAX_SIZE_MB      - Volume maximum size in MB (default: 50)"
	@echo "  TEST_TIMEOUT            - Test timeout (default: 15m)"

# Volume encryption integration tests
test-volume-encryption: build-s3
	@echo "🚀 Starting S3 volume encryption integration tests..."
	@echo "Starting Hanzo cluster with volume encryption enabled..."
	@# Start server with -s3.encryptVolumeData flag
	@mkdir -p /tmp/hanzo-test-sse
	@rm -f /tmp/hanzo-sse-*.log || true
	@sed -e 's/ACCESS_KEY_PLACEHOLDER/$(ACCESS_KEY)/g' \
	     -e 's/SECRET_KEY_PLACEHOLDER/$(SECRET_KEY)/g' \
	     s3-config-template.json > /tmp/hanzo-s3.json
	@echo "Starting s3 mini with S3 volume encryption..."
	@AWS_ACCESS_KEY_ID=$(ACCESS_KEY) AWS_SECRET_ACCESS_KEY=$(SECRET_KEY) WEED_S3_SSE_KEY=test-sse-s3-key GLOG_v=4 $(SEAWEEDFS_BINARY) mini \
		-dir=/tmp/hanzo-test-sse \
		-s3.port=$(S3_PORT) \
		-s3.config=/tmp/hanzo-s3.json \
		-s3.encryptVolumeData \
		-ip=127.0.0.1 \
		> /tmp/hanzo-sse-mini.log 2>&1 & echo $$! > /tmp/s3-mini.pid
	@echo "Checking S3 service is ready..."
	@for i in $$(seq 1 30); do \
		if curl -s http://127.0.0.1:$(S3_PORT) > /dev/null 2>&1; then \
			echo "✅ S3 service is ready"; \
			break; \
		fi; \
		sleep 1; \
	done
	@echo "Running volume encryption integration tests..."
	@trap '$(MAKE) -C $(TEST_DIR) stop-hanzo-safe || true' EXIT; \
	cd $(SEAWEEDFS_ROOT) && go test -v -tags=integration -timeout=10m -run "TestS3VolumeEncryption" ./test/s3/sse || exit 1; \
	echo "✅ Volume encryption tests completed successfully"; \
	$(MAKE) -C $(TEST_DIR) stop-hanzo-safe || true

