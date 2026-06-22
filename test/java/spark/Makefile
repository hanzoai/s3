.PHONY: help build test test-local test-docker clean run-example docker-up docker-down

help:
	@echo "Hanzo Spark Integration Tests"
	@echo ""
	@echo "Available targets:"
	@echo "  build          - Build the project"
	@echo "  test           - Run integration tests (requires Hanzo running)"
	@echo "  test-local     - Run tests against local Hanzo"
	@echo "  test-docker    - Run tests in Docker with Hanzo"
	@echo "  run-example    - Run the example Spark application"
	@echo "  docker-up      - Start Hanzo in Docker"
	@echo "  docker-down    - Stop Hanzo Docker containers"
	@echo "  clean          - Clean build artifacts"

build:
	mvn clean package

test:
	@if [ -z "$$SEAWEEDFS_TEST_ENABLED" ]; then \
		echo "Setting SEAWEEDFS_TEST_ENABLED=true"; \
	fi
	SEAWEEDFS_TEST_ENABLED=true mvn test

test-local:
	@echo "Testing against local Hanzo (localhost:8888)..."
	./run-tests.sh

test-docker:
	@echo "Running tests in Docker..."
	docker compose up --build --abort-on-container-exit spark-tests
	docker compose down

docker-up:
	@echo "Starting Hanzo in Docker..."
	docker compose up -d hanzo-master hanzo-volume hanzo-filer
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Hanzo is ready!"
	@echo "  Master: http://localhost:9333"
	@echo "  Filer:  http://localhost:8888"

docker-down:
	@echo "Stopping Hanzo Docker containers..."
	docker compose down -v

run-example:
	@echo "Running example application..."
	@if ! command -v spark-submit > /dev/null; then \
		echo "Error: spark-submit not found. Please install Apache Spark."; \
		exit 1; \
	fi
	spark-submit \
		--class hanzo.spark.SparkHanzoExample \
		--master local[2] \
		--conf spark.hadoop.fs.s3.impl=hanzo.hdfs.HanzoFileSystem \
		--conf spark.hadoop.fs.hanzo.filer.host=localhost \
		--conf spark.hadoop.fs.hanzo.filer.port=8888 \
		--conf spark.hadoop.fs.hanzo.filer.port.grpc=18888 \
		--conf spark.hadoop.fs.hanzo.replication="" \
		target/hanzo-spark-integration-tests-1.0-SNAPSHOT.jar \
		hanzo://localhost:8888/spark-example-output

clean:
	mvn clean
	@echo "Build artifacts cleaned"

verify-hanzo:
	@echo "Verifying Hanzo connection..."
	@curl -f http://localhost:8888/ > /dev/null 2>&1 && \
		echo "✓ Hanzo filer is accessible" || \
		(echo "✗ Hanzo filer is not accessible at http://localhost:8888"; exit 1)

.DEFAULT_GOAL := help

