#!/bin/bash

# Script to set up Hanzo infrastructure for Kafka Gateway testing
# This script will start Master, Filer, and MQ Broker components

set -e

BASE_DIR="/tmp/hanzo"
LOG_DIR="$BASE_DIR/logs"
DATA_DIR="$BASE_DIR/data"

echo "=== Hanzo Infrastructure Setup ==="
echo "Setup Date: $(date)"
echo "Base Directory: $BASE_DIR"
echo ""

# Create directories
mkdir -p "$BASE_DIR/master" "$BASE_DIR/filer" "$BASE_DIR/broker" "$LOG_DIR"

# Function to check if a service is running
check_service() {
    local host_port=$1
    local service_name=$2
    
    if timeout 3 bash -c "</dev/tcp/${host_port//://}" 2>/dev/null; then
        echo "✓ $service_name is already running on $host_port"
        return 0
    else
        echo "✗ $service_name is NOT running on $host_port"
        return 1
    fi
}

# Function to start a service in background
start_service() {
    local cmd="$1"
    local service_name="$2"
    local log_file="$3"
    local check_port="$4"
    
    echo "Starting $service_name..."
    echo "Command: $cmd"
    echo "Log: $log_file"
    
    # Start in background
    nohup $cmd > "$log_file" 2>&1 &
    local pid=$!
    echo "PID: $pid"
    
    # Wait for service to be ready
    local retries=30
    while [ $retries -gt 0 ]; do
        if check_service "$check_port" "$service_name" 2>/dev/null; then
            echo "✓ $service_name is ready"
            return 0
        fi
        retries=$((retries - 1))
        sleep 1
        echo -n "."
    done
    echo ""
    echo "❌ $service_name failed to start within 30 seconds"
    return 1
}

# Stop any existing processes
echo "=== Cleaning up existing processes ==="
pkill -f "s3 master" || true
pkill -f "s3 filer" || true  
pkill -f "s3 mq.broker" || true
sleep 2

echo ""
echo "=== Starting Hanzo Components ==="

# Start Master
if ! check_service "localhost:9333" "Hanzo Master"; then
    start_service \
        "s3 master -defaultReplication=001 -mdir=$BASE_DIR/master -peers=none" \
        "Hanzo Master" \
        "$LOG_DIR/master.log" \
        "localhost:9333"
    echo ""
fi

# Start Filer
if ! check_service "localhost:8888" "Hanzo Filer"; then
    start_service \
        "s3 filer -master=localhost:9333 -filer.dir=$BASE_DIR/filer" \
        "Hanzo Filer" \
        "$LOG_DIR/filer.log" \
        "localhost:8888"
    echo ""
fi

# Start MQ Broker
if ! check_service "localhost:17777" "Hanzo MQ Broker"; then
    start_service \
        "s3 mq.broker -filer=localhost:8888 -master=localhost:9333" \
        "Hanzo MQ Broker" \
        "$LOG_DIR/broker.log" \
        "localhost:17777"
    echo ""
fi

echo "=== Infrastructure Status ==="
check_service "localhost:9333" "Master (gRPC)"
check_service "localhost:9334" "Master (HTTP)"
check_service "localhost:8888" "Filer (HTTP)"
check_service "localhost:18888" "Filer (gRPC)"
check_service "localhost:17777" "MQ Broker"

echo ""
echo "=== Infrastructure Ready ==="
echo "Log files:"
echo "  Master: $LOG_DIR/master.log"
echo "  Filer:  $LOG_DIR/filer.log"
echo "  Broker: $LOG_DIR/broker.log"
echo ""
echo "To view logs in real-time:"
echo "  tail -f $LOG_DIR/master.log"
echo "  tail -f $LOG_DIR/filer.log"
echo "  tail -f $LOG_DIR/broker.log"
echo ""
echo "To stop all services:"
echo "  pkill -f \"s3 master\""
echo "  pkill -f \"s3 filer\""
echo "  pkill -f \"s3 mq.broker\""
echo ""
echo "[OK] Hanzo infrastructure is ready for testing!"

