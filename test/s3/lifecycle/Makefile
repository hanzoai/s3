# S3 Lifecycle Test Makefile
# End-to-end test of the event-driven lifecycle worker, driven by
# the s3.lifecycle.run-shard shell command.

.PHONY: help build-s3 start-server stop-server test test-with-server clean health-check

WEED_BINARY := ../../../s3/s3_binary
S3_PORT := 8333
MASTER_PORT := 9333
VOLUME_PORT := 8080
FILER_PORT := 8888
# Pin the filer gRPC port off the FILER_PORT+10000 convention so any
# code path that assumes the offset breaks here, not in production.
FILER_GRPC_PORT := 18890
ADMIN_PORT := 23646
ACCESS_KEY ?= some_access_key1
SECRET_KEY ?= some_secret_key1
TEST_TIMEOUT := 10m
TEST_PATTERN := TestLifecycle
SERVER_DIR := ./test-volume-data/server-data

help:
	@echo "S3 Lifecycle Test Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build-s3       - build the Hanzo binary"
	@echo "  start-server     - start a local 's3 mini' cluster"
	@echo "  stop-server      - stop the cluster"
	@echo "  test             - run tests against an already-running cluster"
	@echo "  test-with-server - start, run tests, stop"
	@echo "  health-check     - is the S3 endpoint up?"
	@echo "  clean            - remove test artifacts"

build-s3:
	@echo "Building Hanzo binary..."
	@cd ../../../s3 && go build -o s3_binary .
	@chmod +x $(WEED_BINARY)

start-server: build-s3
	@echo "Starting s3 mini..."
	@rm -f s3-server.pid
	@mkdir -p $(SERVER_DIR)
	@AWS_ACCESS_KEY_ID=$(ACCESS_KEY) AWS_SECRET_ACCESS_KEY=$(SECRET_KEY) $(WEED_BINARY) mini \
		-dir=$(SERVER_DIR) \
		-s3.port=$(S3_PORT) \
		-filer.port.grpc=$(FILER_GRPC_PORT) \
		> s3-test.log 2>&1 & \
		echo $$! > s3-server.pid
	@for i in $$(seq 1 90); do \
		if curl -s http://localhost:$(S3_PORT) >/dev/null 2>&1; then \
			echo "server up after $$i s"; exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "server did not start within 90s"; \
	tail -50 s3-test.log; \
	exit 1

stop-server:
	@if [ -f s3-server.pid ]; then \
		PID=$$(cat s3-server.pid); \
		kill -TERM $$PID 2>/dev/null || true; \
		sleep 1; \
		kill -KILL $$PID 2>/dev/null || true; \
		rm -f s3-server.pid; \
	fi

health-check:
	@curl -s http://localhost:$(S3_PORT) >/dev/null && echo "S3 up on $(S3_PORT)" || (echo "S3 not reachable"; exit 1)

test:
	@WEED_BINARY=$$(cd ../../../s3 && pwd)/s3_binary \
	S3_ENDPOINT=http://localhost:$(S3_PORT) \
	S3_GRPC_ENDPOINT=localhost:$$(($(S3_PORT) + 10000)) \
	MASTER_ENDPOINT=http://localhost:$(MASTER_PORT) \
	FILER_GRPC_ADDRESS=localhost:$(FILER_GRPC_PORT) \
	ADMIN_ENDPOINT=http://localhost:$(ADMIN_PORT) \
	go test -v -timeout $(TEST_TIMEOUT) -run $(TEST_PATTERN)

test-with-server: start-server
	@$(MAKE) test || (RC=$$?; $(MAKE) stop-server; exit $$RC)
	@$(MAKE) stop-server

clean: stop-server
	@rm -rf $(SERVER_DIR) test-volume-data
	@rm -f s3-test.log s3-server.pid
