# Makefile for TUS Protocol Integration Tests
# This Makefile provides targets for running TUS (resumable upload) integration tests

# Default values
SEAWEEDFS_BINARY ?= s3
FILER_PORT ?= 18888
VOLUME_PORT ?= 18080
MASTER_PORT ?= 19333
TEST_TIMEOUT ?= 10m
VOLUME_MAX_SIZE_MB ?= 50
VOLUME_MAX_COUNT ?= 100

# Test directory
TEST_DIR := $(shell pwd)
SEAWEEDFS_ROOT := $(shell cd ../.. && pwd)

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m # No Color

.PHONY: all test clean start-hanzo stop-hanzo check-binary build-s3 help test-basic test-chunked test-resume test-errors test-with-server

all: test

# Build Hanzo binary
build-s3:
	@echo "Building Hanzo binary..."
	@cd $(SEAWEEDFS_ROOT)/s3 && go build -o s3
	@echo "$(GREEN)Hanzo binary built successfully$(NC)"

help:
	@echo "Hanzo TUS Protocol Integration Tests"
	@echo ""
	@echo "Available targets:"
	@echo "  test              - Run all TUS integration tests"
	@echo "  test-basic        - Run basic TUS upload tests"
	@echo "  test-chunked      - Run chunked upload tests"
	@echo "  test-resume       - Run upload resume tests"
	@echo "  test-errors       - Run error handling tests"
	@echo "  test-with-server  - Run tests with automatic server management"
	@echo "  start-hanzo   - Start Hanzo server for testing"
	@echo "  stop-hanzo    - Stop Hanzo server"
	@echo "  clean             - Clean up test artifacts"
	@echo "  check-binary      - Check if Hanzo binary exists"
	@echo "  build-s3        - Build Hanzo binary"
	@echo ""
	@echo "Configuration:"
	@echo "  SEAWEEDFS_BINARY=$(SEAWEEDFS_BINARY)"
	@echo "  FILER_PORT=$(FILER_PORT)"
	@echo "  VOLUME_PORT=$(VOLUME_PORT)"
	@echo "  MASTER_PORT=$(MASTER_PORT)"
	@echo "  TEST_TIMEOUT=$(TEST_TIMEOUT)"

check-binary:
	@if ! command -v $(SEAWEEDFS_BINARY) > /dev/null 2>&1 && [ ! -f "$(SEAWEEDFS_ROOT)/s3/s3" ]; then \
		echo "$(RED)Error: Hanzo binary not found$(NC)"; \
		echo "Please build Hanzo first: make build-s3"; \
		exit 1; \
	fi
	@echo "$(GREEN)Hanzo binary found$(NC)"

start-hanzo: check-binary
	@echo "$(YELLOW)Starting Hanzo server for TUS testing...$(NC)"
	@# Clean up any existing processes on our test ports
	@lsof -ti :$(MASTER_PORT) | grep -v '^$$' | xargs kill -TERM 2>/dev/null || true
	@lsof -ti :$(VOLUME_PORT) | grep -v '^$$' | xargs kill -TERM 2>/dev/null || true
	@lsof -ti :$(FILER_PORT) | grep -v '^$$' | xargs kill -TERM 2>/dev/null || true
	@sleep 2
	
	# Create necessary directories
	@mkdir -p /tmp/hanzo-test-tus-master
	@mkdir -p /tmp/hanzo-test-tus-volume
	@mkdir -p /tmp/hanzo-test-tus-filer
	
	# Start master server (use freshly built binary)
	@echo "Starting master server..."
	@nohup $(SEAWEEDFS_ROOT)/s3/s3 master \
		-port=$(MASTER_PORT) \
		-mdir=/tmp/hanzo-test-tus-master \
		-volumeSizeLimitMB=$(VOLUME_MAX_SIZE_MB) \
		-ip=127.0.0.1 \
		> /tmp/hanzo-tus-master.log 2>&1 &
	@sleep 3
	
	# Start volume server
	@echo "Starting volume server..."
	@nohup $(SEAWEEDFS_ROOT)/s3/s3 volume \
		-port=$(VOLUME_PORT) \
		-mserver=127.0.0.1:$(MASTER_PORT) \
		-dir=/tmp/hanzo-test-tus-volume \
		-max=$(VOLUME_MAX_COUNT) \
		-ip=127.0.0.1 \
		> /tmp/hanzo-tus-volume.log 2>&1 &
	@sleep 3
	
	# Start filer server with TUS enabled (default tusBasePath is .tus)
	@echo "Starting filer server..."
	@nohup $(SEAWEEDFS_ROOT)/s3/s3 filer \
		-port=$(FILER_PORT) \
		-master=127.0.0.1:$(MASTER_PORT) \
		-ip=127.0.0.1 \
		> /tmp/hanzo-tus-filer.log 2>&1 &
	@sleep 5
	
	# Wait for filer to be ready
	@echo "$(YELLOW)Waiting for filer to be ready...$(NC)"
	@for i in $$(seq 1 30); do \
		if curl -s -f http://127.0.0.1:$(FILER_PORT)/ > /dev/null 2>&1; then \
			echo "$(GREEN)Filer is ready$(NC)"; \
			break; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "$(RED)Filer failed to start within 30 seconds$(NC)"; \
			$(MAKE) debug-logs; \
			exit 1; \
		fi; \
		echo "Waiting for filer... ($$i/30)"; \
		sleep 1; \
	done
	
	@echo "$(GREEN)Hanzo server started successfully for TUS testing$(NC)"
	@echo "Master: http://localhost:$(MASTER_PORT)"
	@echo "Volume: http://localhost:$(VOLUME_PORT)"
	@echo "Filer: http://localhost:$(FILER_PORT)"
	@echo "TUS Endpoint: http://localhost:$(FILER_PORT)/.tus/"

stop-hanzo:
	@echo "$(YELLOW)Stopping Hanzo server...$(NC)"
	@lsof -ti :$(MASTER_PORT) | grep -v '^$$' | xargs kill -TERM 2>/dev/null || true
	@lsof -ti :$(VOLUME_PORT) | grep -v '^$$' | xargs kill -TERM 2>/dev/null || true
	@lsof -ti :$(FILER_PORT) | grep -v '^$$' | xargs kill -TERM 2>/dev/null || true
	@sleep 2
	@echo "$(GREEN)Hanzo server stopped$(NC)"

clean:
	@echo "$(YELLOW)Cleaning up TUS test artifacts...$(NC)"
	@rm -rf /tmp/hanzo-test-tus-*
	@rm -f /tmp/hanzo-tus-*.log
	@echo "$(GREEN)TUS test cleanup completed$(NC)"

# Run all tests
test: check-binary
	@echo "$(YELLOW)Running all TUS integration tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) ./test/tus/...
	@echo "$(GREEN)All TUS tests completed$(NC)"

# Run basic upload tests
test-basic: check-binary
	@echo "$(YELLOW)Running basic TUS upload tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestTusBasicUpload|TestTusOptionsHandler" ./test/tus/...
	@echo "$(GREEN)Basic TUS tests completed$(NC)"

# Run chunked upload tests
test-chunked: check-binary
	@echo "$(YELLOW)Running chunked TUS upload tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestTusChunkedUpload" ./test/tus/...
	@echo "$(GREEN)Chunked TUS tests completed$(NC)"

# Run resume tests
test-resume: check-binary
	@echo "$(YELLOW)Running TUS upload resume tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestTusResumeAfterInterruption|TestTusHeadRequest" ./test/tus/...
	@echo "$(GREEN)TUS resume tests completed$(NC)"

# Run error handling tests
test-errors: check-binary
	@echo "$(YELLOW)Running TUS error handling tests...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) -run "TestTusInvalidOffset|TestTusUploadNotFound|TestTusDeleteUpload" ./test/tus/...
	@echo "$(GREEN)TUS error tests completed$(NC)"

# Run tests with automatic server management
test-with-server: build-s3
	@echo "$(YELLOW)Running TUS tests with automatic server management...$(NC)"
	@$(MAKE) -C $(TEST_DIR) start-hanzo && \
	sleep 3 && \
	cd $(SEAWEEDFS_ROOT) && go test -v -timeout=$(TEST_TIMEOUT) ./test/tus/...; \
	TEST_RESULT=$$?; \
	$(MAKE) -C $(TEST_DIR) stop-hanzo; \
	$(MAKE) -C $(TEST_DIR) clean; \
	if [ $$TEST_RESULT -eq 0 ]; then echo "$(GREEN)All TUS tests passed!$(NC)"; fi; \
	exit $$TEST_RESULT

# Debug targets
debug-logs:
	@echo "$(YELLOW)=== Master Log ===$(NC)"
	@tail -n 50 /tmp/hanzo-tus-master.log 2>/dev/null || echo "No master log found"
	@echo "$(YELLOW)=== Volume Log ===$(NC)"
	@tail -n 50 /tmp/hanzo-tus-volume.log 2>/dev/null || echo "No volume log found"
	@echo "$(YELLOW)=== Filer Log ===$(NC)"
	@tail -n 50 /tmp/hanzo-tus-filer.log 2>/dev/null || echo "No filer log found"

debug-status:
	@echo "$(YELLOW)=== Process Status ===$(NC)"
	@ps aux | grep -E "(s3|hanzo)" | grep -v grep || echo "No Hanzo processes found"
	@echo "$(YELLOW)=== Port Status ===$(NC)"
	@lsof -i :$(MASTER_PORT) -i :$(VOLUME_PORT) -i :$(FILER_PORT) 2>/dev/null || echo "No ports in use"

# Manual testing targets
manual-start: start-hanzo
	@echo "$(GREEN)Hanzo is now running for manual TUS testing$(NC)"
	@echo ""
	@echo "TUS Endpoints:"
	@echo "  OPTIONS /.tus/           - Capability discovery"
	@echo "  POST    /.tus/{path}     - Create upload"
	@echo "  HEAD    /.tus/.uploads/{id} - Get offset"
	@echo "  PATCH   /.tus/.uploads/{id} - Upload data"
	@echo "  DELETE  /.tus/.uploads/{id} - Cancel upload"
	@echo ""
	@echo "Example curl commands:"
	@echo "  curl -X OPTIONS http://localhost:$(FILER_PORT)/.tus/ -H 'Tus-Resumable: 1.0.0'"
	@echo ""
	@echo "Run 'make manual-stop' when finished"

manual-stop: stop-hanzo clean

# CI targets
ci-test: test-with-server

# Skip integration tests (short mode)
test-short:
	@echo "$(YELLOW)Running TUS tests in short mode (skipping integration tests)...$(NC)"
	@cd $(SEAWEEDFS_ROOT) && go test -v -short ./test/tus/...
	@echo "$(GREEN)Short tests completed$(NC)"

