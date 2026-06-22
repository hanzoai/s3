# Makefile for Iceberg REST Catalog Integration Tests
#
# This Makefile provides easy-to-use targets for running integration tests
# and starting a local development environment.

.PHONY: build mini test test-pyiceberg clean

# Root directory of the Hanzo project
PROJECT_ROOT = ../../..
WEED_BIN = $(PROJECT_ROOT)/s3/s3

# Build the s3 binary
build:
	@echo "Building Hanzo binary..."
	cd $(PROJECT_ROOT)/s3 && go build -v

# Start Hanzo mini cluster for manual testing
# Uses default ports: S3=8333, Iceberg=8182
mini: build
	@echo "Starting Hanzo mini cluster..."
	$(WEED_BIN) mini -s3.port.iceberg=8182 & echo $$! > s3-mini.pid; wait

# Run all integration tests in this directory
test: build
	@echo "Running all Iceberg catalog integration tests..."
	go test -v .

# Run only the PyIceberg compatibility tests (requires Docker)
test-pyiceberg: build
	@echo "Running PyIceberg compatibility tests..."
	go test -v -run TestPyIcebergRestCatalog .

# Clean up temporary data and stop s3 mini if running
clean:
	@echo "Cleaning up..."
	@test -f s3-mini.pid && kill $$(cat s3-mini.pid) && rm s3-mini.pid || true
	@rm -rf /tmp/hanzo-iceberg-test-*
