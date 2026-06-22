# S3 delete regression test harness

.PHONY: all test help build-s3 check-deps start-server stop-server test-delete test-with-server clean logs

WEED_BINARY := ../../../s3/s3_binary
S3_PORT := 8333
TEST_TIMEOUT := 10m
SERVER_DIR := ./test-volume-data
S3_CONFIG := ./test_s3.json

.DEFAULT_GOAL := help

all: test

test: test-with-server

help:
	@echo "S3 delete regression test harness"
	@echo ""
	@echo "Available targets:"
	@echo "  build-s3      - Build the Hanzo binary"
	@echo "  start-server    - Start Hanzo with S3 enabled"
	@echo "  stop-server     - Stop the test server"
	@echo "  test-delete     - Run delete regression tests"
	@echo "  test-with-server - Start server, run tests, stop server"
	@echo "  logs            - Show server logs"
	@echo "  clean           - Remove test artifacts"

build-s3:
	@echo "Building Hanzo binary..."
	@cd ../../../s3 && go build -o s3_binary .
	@chmod +x $(WEED_BINARY)

check-deps: build-s3
	@command -v go >/dev/null 2>&1 || (echo "Go is required but not installed" && exit 1)
	@test -f $(WEED_BINARY) || (echo "Hanzo binary not found at $(WEED_BINARY)" && exit 1)
	@test -f $(S3_CONFIG) || (echo "S3 config not found at $(S3_CONFIG)" && exit 1)

start-server: check-deps
	@echo "Starting Hanzo server for delete regression tests..."
	@rm -f s3-server.pid
	@mkdir -p $(SERVER_DIR)
	@$(WEED_BINARY) mini \
		-dir=$(SERVER_DIR) \
		-s3 \
		-s3.port=$(S3_PORT) \
		-s3.config=$(S3_CONFIG) \
		-master.peers=none \
		> s3-test.log 2>&1 & echo $$! > s3-server.pid
	@for i in $$(seq 1 60); do \
		if curl -sS -o /dev/null http://127.0.0.1:$(S3_PORT)/ >/dev/null 2>&1; then \
			echo "Hanzo S3 server is ready on port $(S3_PORT)"; \
			exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "Hanzo S3 server failed to start"; \
	echo "=== Server logs ==="; \
	test -f s3-test.log && cat s3-test.log || true; \
	exit 1

stop-server:
	@echo "Stopping Hanzo server..."
	@if [ -f s3-server.pid ]; then \
		SERVER_PID=$$(cat s3-server.pid); \
		kill -TERM $$SERVER_PID 2>/dev/null || true; \
		for i in $$(seq 1 20); do \
			if ! kill -0 $$SERVER_PID >/dev/null 2>&1; then \
				break; \
			fi; \
			sleep 0.1; \
		done; \
		if kill -0 $$SERVER_PID >/dev/null 2>&1; then \
			kill -KILL $$SERVER_PID 2>/dev/null || true; \
		fi; \
		rm -f s3-server.pid; \
	fi

test-delete: check-deps
	@echo "Running delete regression tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) .

test-with-server: start-server
	@trap "$(MAKE) stop-server" EXIT; \
	$(MAKE) test-delete || (echo "=== Last 100 lines of server logs ===" && tail -100 s3-test.log && exit 1)
	@$(MAKE) stop-server

logs:
	@if test -f s3-test.log; then \
		tail -f s3-test.log; \
	else \
		echo "No server log file found."; \
	fi

clean:
	@$(MAKE) stop-server
	@rm -f s3-test.log s3-server.pid
	@rm -rf $(SERVER_DIR)
