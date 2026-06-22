#!/bin/bash

set -e

echo "=== Hanzo Spark Integration Tests Quick Start ==="
echo ""

# Check if Hanzo is running
check_hanzo() {
    echo "Checking if Hanzo is running..."
    if curl -f http://localhost:8888/ > /dev/null 2>&1; then
        echo "✓ Hanzo filer is accessible at http://localhost:8888"
        return 0
    else
        echo "✗ Hanzo filer is not accessible"
        return 1
    fi
}

# Start Hanzo with Docker if not running
start_hanzo() {
    echo ""
    echo "Starting Hanzo with Docker..."
    docker-compose up -d hanzo-master hanzo-volume hanzo-filer
    
    echo "Waiting for Hanzo to be ready..."
    for i in {1..30}; do
        if curl -f http://localhost:8888/ > /dev/null 2>&1; then
            echo "✓ Hanzo is ready!"
            return 0
        fi
        echo -n "."
        sleep 2
    done
    
    echo ""
    echo "✗ Hanzo failed to start"
    return 1
}

# Build the project
build_project() {
    echo ""
    echo "Building the project..."
    mvn clean package -DskipTests
    echo "✓ Build completed"
}

# Run tests
run_tests() {
    echo ""
    echo "Running integration tests..."
    export SEAWEEDFS_TEST_ENABLED=true
    mvn test
    echo "✓ Tests completed"
}

# Run example
run_example() {
    echo ""
    echo "Running example application..."
    
    if ! command -v spark-submit > /dev/null; then
        echo "⚠ spark-submit not found. Skipping example application."
        echo "To run the example, install Apache Spark and try: make run-example"
        return 0
    fi
    
    spark-submit \
        --class hanzo.spark.SparkHanzoExample \
        --master local[2] \
        --conf spark.hadoop.fs.s3.impl=hanzo.hdfs.HanzoFileSystem \
        --conf spark.hadoop.fs.hanzo.filer.host=localhost \
        --conf spark.hadoop.fs.hanzo.filer.port=8888 \
        --conf spark.hadoop.fs.hanzo.filer.port.grpc=18888 \
        target/hanzo-spark-integration-tests-1.0-SNAPSHOT.jar \
        hanzo://localhost:8888/spark-quickstart-output
    
    echo "✓ Example completed"
}

# Cleanup
cleanup() {
    echo ""
    echo "Cleaning up..."
    docker-compose down -v
    echo "✓ Cleanup completed"
}

# Main execution
main() {
    # Check if Docker is available
    if ! command -v docker > /dev/null; then
        echo "Error: Docker is not installed or not in PATH"
        exit 1
    fi

    # Check if Maven is available
    if ! command -v mvn > /dev/null; then
        echo "Error: Maven is not installed or not in PATH"
        exit 1
    fi

    # Check if Hanzo is running, if not start it
    if ! check_hanzo; then
        read -p "Do you want to start Hanzo with Docker? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            start_hanzo || exit 1
        else
            echo "Please start Hanzo manually and rerun this script."
            exit 1
        fi
    fi

    # Build project
    build_project || exit 1

    # Run tests
    run_tests || exit 1

    # Run example if Spark is available
    run_example

    echo ""
    echo "=== Quick Start Completed Successfully! ==="
    echo ""
    echo "Next steps:"
    echo "  - View test results in target/surefire-reports/"
    echo "  - Check example output at http://localhost:8888/"
    echo "  - Run 'make help' for more options"
    echo "  - Read README.md for detailed documentation"
    echo ""
    
    read -p "Do you want to stop Hanzo? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cleanup
    fi
}

# Handle Ctrl+C
trap cleanup INT

# Run main
main



