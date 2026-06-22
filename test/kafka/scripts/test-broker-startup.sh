#!/bin/bash

# Script to test Hanzo MQ broker startup locally
# This helps debug broker startup issues before running CI

set -e

echo "=== Testing Hanzo MQ Broker Startup ==="

# Build s3 binary
echo "Building s3 binary..."
cd "$(dirname "$0")/../../.."
go build -o /tmp/s3 ./s3

# Setup data directory
WEED_DATA_DIR="/tmp/hanzo-broker-test-$$"
mkdir -p "$WEED_DATA_DIR"
echo "Using data directory: $WEED_DATA_DIR"

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    pkill -f "s3.*server" || true
    pkill -f "s3.*mq.broker" || true
    sleep 2
    rm -rf "$WEED_DATA_DIR"
    rm -f /tmp/s3-*.log
}
trap cleanup EXIT

# Start Hanzo server  
echo "Starting Hanzo server..."
/tmp/s3 -v 1 server \
  -ip="127.0.0.1" \
  -ip.bind="0.0.0.0" \
  -dir="$WEED_DATA_DIR" \
  -master.raftHashicorp \
  -master.port=9333 \
  -volume.port=8081 \
  -filer.port=8888 \
  -filer=true \
  -metricsPort=9325 \
  > /tmp/s3-server-test.log 2>&1 &

SERVER_PID=$!
echo "Server PID: $SERVER_PID"

# Wait for master
echo "Waiting for master..."
for i in $(seq 1 30); do
  if curl -s http://127.0.0.1:9333/cluster/status >/dev/null; then
    echo "✓ Master is up"
    break
  fi
  echo "  Waiting for master... ($i/30)"
  sleep 1
done

# Wait for filer
echo "Waiting for filer..."
for i in $(seq 1 30); do
  if nc -z 127.0.0.1 8888; then
    echo "✓ Filer is up"
    break
  fi
  echo "  Waiting for filer... ($i/30)"
  sleep 1
done

# Start MQ broker  
echo "Starting MQ broker..."
/tmp/s3 -v 2 mq.broker \
  -master="127.0.0.1:9333" \
  -ip="127.0.0.1" \
  -port=17777 \
  > /tmp/s3-mq-broker-test.log 2>&1 &

BROKER_PID=$!
echo "Broker PID: $BROKER_PID"

# Wait for broker
echo "Waiting for broker..."
broker_ready=false
for i in $(seq 1 30); do
  if nc -z 127.0.0.1 17777; then
    echo "✓ MQ broker is up"
    broker_ready=true
    break
  fi
  echo "  Waiting for MQ broker... ($i/30)"
  sleep 1
done

if [ "$broker_ready" = false ]; then
  echo "❌ MQ broker failed to start"
  echo
  echo "=== Server logs ==="
  cat /tmp/s3-server-test.log
  echo
  echo "=== Broker logs ==="
  cat /tmp/s3-mq-broker-test.log
  exit 1
fi

# Broker started successfully - discovery will be tested by Kafka gateway
echo "✓ Broker started successfully and accepting connections"

echo
echo "[OK] All tests passed!"
echo "Server logs: /tmp/s3-server-test.log"  
echo "Broker logs: /tmp/s3-mq-broker-test.log"
