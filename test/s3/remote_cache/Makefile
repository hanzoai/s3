# Remote Storage Cache Integration Tests
# Tests the remote object caching functionality with singleflight deduplication
# Uses two Hanzo instances: primary (with caching) and secondary (as remote storage)

.PHONY: all help build-s3 check-deps start-remote stop-remote start-primary stop-primary \
        setup-remote test test-with-server clean logs logs-primary logs-remote health

all: test-with-server

# Configuration
WEED_BINARY := ../../../s3/s3_binary

ACCESS_KEY ?= some_access_key1
SECRET_KEY ?= some_secret_key1

# Primary Hanzo (the one being tested - has remote caching)
PRIMARY_S3_PORT := 8333
PRIMARY_MASTER_PORT := 9333
PRIMARY_FILER_PORT := 8888
PRIMARY_VOLUME_PORT := 9340
PRIMARY_WEBDAV_PORT := 7333
PRIMARY_METRICS_PORT := 9324
PRIMARY_DIR := ./test-primary-data

# Secondary Hanzo (acts as "remote" S3 storage)
REMOTE_S3_PORT := 8334
REMOTE_MASTER_PORT := 9334
REMOTE_FILER_PORT := 8889
REMOTE_VOLUME_PORT := 9341
REMOTE_WEBDAV_PORT := 7334
REMOTE_METRICS_PORT := 9325
REMOTE_DIR := ./test-remote-data

# Test configuration
TEST_TIMEOUT := 15m
TEST_PATTERN := .

# Buckets
REMOTE_BUCKET := remotesourcebucket

# Default target
help:
	@echo "Remote Storage Cache Integration Tests"
	@echo ""
	@echo "Uses two Hanzo instances:"
	@echo "  - Primary (port $(PRIMARY_S3_PORT)): Being tested, has remote caching"
	@echo "  - Remote  (port $(REMOTE_S3_PORT)): Acts as remote S3 storage"
	@echo ""
	@echo "Available targets:"
	@echo "  help              - Show this help message"
	@echo "  build-s3        - Build the Hanzo binary"
	@echo "  check-deps        - Check dependencies"
	@echo "  start-remote      - Start remote Hanzo (secondary)"
	@echo "  stop-remote       - Stop remote Hanzo"
	@echo "  start-primary     - Start primary Hanzo"
	@echo "  stop-primary      - Stop primary Hanzo"
	@echo "  setup-remote      - Configure remote storage mount"
	@echo "  test              - Run tests (assumes servers are running)"
	@echo "  test-with-server  - Start servers, run tests, stop servers"
	@echo "  clean             - Clean up all resources"
	@echo "  logs              - Show server logs"
	@echo "  health            - Check server health"

# Build the Hanzo binary
build-s3:
	@echo "Building Hanzo binary..."
	@cd ../../../s3 && go build -o s3_binary .
	@chmod +x $(WEED_BINARY)
	@echo "Hanzo binary built"

check-deps: build-s3
	@echo "Checking dependencies..."
	@command -v go >/dev/null 2>&1 || (echo "Go is required" && exit 1)
	@test -f $(WEED_BINARY) || (echo "Hanzo binary not found" && exit 1)
	@echo "All dependencies available"

# Start remote Hanzo (acts as the "remote" S3 storage)
start-remote: check-deps
	@echo "Starting remote Hanzo (secondary instance)..."
	@rm -f remote-server.pid
	@mkdir -p $(REMOTE_DIR)
	@$(WEED_BINARY) mini \
		-s3.port=$(REMOTE_S3_PORT) \
		-master.port=$(REMOTE_MASTER_PORT) \
		-filer.port=$(REMOTE_FILER_PORT) \
		-volume.port=$(REMOTE_VOLUME_PORT) \
		-webdav.port=$(REMOTE_WEBDAV_PORT) \
		-s3.allowDeleteBucketNotEmpty=true \
		-s3.config=s3_config.json \
		-dir=$(REMOTE_DIR) \
		-ip=127.0.0.1 \
		-ip.bind=127.0.0.1 \
		-metricsPort=$(REMOTE_METRICS_PORT) \
		> remote-s3.log 2>&1 & echo $$! > remote-server.pid
	@echo "Waiting for remote Hanzo to start..."
	@for i in $$(seq 1 60); do \
		if curl -s http://localhost:$(REMOTE_S3_PORT) >/dev/null 2>&1; then \
			echo "Remote Hanzo started on port $(REMOTE_S3_PORT)"; \
			exit 0; \
		fi; \
		sleep 3; \
	done; \
	echo "ERROR: Remote Hanzo failed to start"; \
	cat remote-s3.log; \
	exit 1

stop-remote:
	@echo "Stopping remote Hanzo..."
	@if [ -f remote-server.pid ]; then \
		kill -TERM $$(cat remote-server.pid) 2>/dev/null || true; \
		sleep 2; \
		kill -KILL $$(cat remote-server.pid) 2>/dev/null || true; \
		rm -f remote-server.pid; \
	fi
	@echo "Remote Hanzo stopped"

# Start primary Hanzo (the one being tested)
start-primary: check-deps
	@echo "Starting primary Hanzo..."
	@rm -f primary-server.pid
	@mkdir -p $(PRIMARY_DIR)
	@$(WEED_BINARY) mini \
		-s3.port=$(PRIMARY_S3_PORT) \
		-master.port=$(PRIMARY_MASTER_PORT) \
		-filer.port=$(PRIMARY_FILER_PORT) \
		-volume.port=$(PRIMARY_VOLUME_PORT) \
		-webdav.port=$(PRIMARY_WEBDAV_PORT) \
		-s3.allowDeleteBucketNotEmpty=true \
		-s3.config=s3_config.json \
		-volume.allowUntrustedRemoteEndpoints \
		-dir=$(PRIMARY_DIR) \
		-ip=127.0.0.1 \
		-ip.bind=127.0.0.1 \
		-metricsPort=$(PRIMARY_METRICS_PORT) \
		> primary-s3.log 2>&1 & echo $$! > primary-server.pid
	@echo "Waiting for primary Hanzo to start..."
	@for i in $$(seq 1 60); do \
		if curl -s http://localhost:$(PRIMARY_S3_PORT) >/dev/null 2>&1; then \
			echo "Primary Hanzo started on port $(PRIMARY_S3_PORT)"; \
			exit 0; \
		fi; \
		sleep 3; \
	done; \
	echo "ERROR: Primary Hanzo failed to start"; \
	cat primary-s3.log; \
	exit 1

stop-primary:
	@echo "Stopping primary Hanzo..."
	@if [ -f primary-server.pid ]; then \
		kill -TERM $$(cat primary-server.pid) 2>/dev/null || true; \
		sleep 2; \
		kill -KILL $$(cat primary-server.pid) 2>/dev/null || true; \
		rm -f primary-server.pid; \
	fi
	@echo "Primary Hanzo stopped"

# Create bucket on remote and configure remote storage mount on primary
setup-remote:
	@echo "Creating bucket on remote Hanzo..."
	@go run utils/create_bucket.go http://localhost:$(REMOTE_S3_PORT) $(ACCESS_KEY) $(SECRET_KEY) $(REMOTE_BUCKET)
	@sleep 3
	@echo "Configuring remote storage on primary..."
	@printf 'remote.configure -name=hanzoremote -type=s3 -s3.access_key=$(ACCESS_KEY) -s3.secret_key=$(SECRET_KEY) -s3.endpoint=http://localhost:$(REMOTE_S3_PORT) -s3.region=us-east-1\nexit\n' | $(WEED_BINARY) shell -master=localhost:$(PRIMARY_MASTER_PORT)
	@sleep 2
	@echo "Mounting remote bucket on primary..."
	@printf 'remote.mount -dir=/buckets/remotemounted -remote=hanzoremote/$(REMOTE_BUCKET) -nonempty\nexit\n' | $(WEED_BINARY) shell -master=localhost:$(PRIMARY_MASTER_PORT)
	@sleep 5
	@printf 'remote.mount\nexit\n' | $(WEED_BINARY) shell -master=localhost:$(PRIMARY_MASTER_PORT) | grep -q "/buckets/remotemounted" || (echo "Mount failed" && exit 1)
	@echo "Remote storage configured and verified"

# Run tests
test: build-s3
	@echo "Running remote cache tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) -run "$(TEST_PATTERN)" .
	@echo "Tests completed"

# Full test workflow
test-with-server: start-remote start-primary
	@sleep 5
	@$(MAKE) setup-remote || (echo "Remote setup failed" && $(MAKE) stop-primary stop-remote && exit 1)
	@sleep 5
	@echo "Running remote cache tests..."
	@$(MAKE) test || (echo "Tests failed" && tail -50 primary-s3.log && $(MAKE) stop-primary stop-remote && exit 1)
	@$(MAKE) stop-primary stop-remote
	@echo "All tests passed"


# Show logs
logs:
	@echo "=== Primary Hanzo Logs ==="
	@if [ -f primary-s3.log ]; then tail -50 primary-s3.log; else echo "No log file"; fi
	@echo ""
	@echo "=== Remote Hanzo Logs ==="
	@if [ -f remote-s3.log ]; then tail -50 remote-s3.log; else echo "No log file"; fi

logs-primary:
	@if [ -f primary-s3.log ]; then tail -f primary-s3.log; else echo "No log file"; fi

logs-remote:
	@if [ -f remote-s3.log ]; then tail -f remote-s3.log; else echo "No log file"; fi

# Clean up
clean:
	@$(MAKE) stop-primary
	@$(MAKE) stop-remote
	@rm -f primary-s3.log remote-s3.log primary-server.pid remote-server.pid
	@rm -rf $(PRIMARY_DIR) $(REMOTE_DIR)
	@rm -f remote_cache.test
	@go clean -testcache
	@echo "Cleanup completed"

# Health check
health:
	@echo "Checking server status..."
	@curl -s http://localhost:$(PRIMARY_S3_PORT) >/dev/null 2>&1 && echo "Primary S3 ($(PRIMARY_S3_PORT)): UP" || echo "Primary S3 ($(PRIMARY_S3_PORT)): DOWN"
	@curl -s http://localhost:$(REMOTE_S3_PORT) >/dev/null 2>&1 && echo "Remote S3 ($(REMOTE_S3_PORT)): UP" || echo "Remote S3 ($(REMOTE_S3_PORT)): DOWN"
