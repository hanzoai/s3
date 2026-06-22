#!/bin/bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration - assumes script is run from s3-rdma-sidecar directory
SEAWEEDFS_DIR="$(realpath ..)"
SIDECAR_DIR="$(pwd)"
MOUNT_POINT="/tmp/hanzo-rdma-mount"
FILER_ADDR="localhost:8888"
SIDECAR_ADDR="localhost:8081"

# PIDs for cleanup
MASTER_PID=""
VOLUME_PID=""
FILER_PID=""
SIDECAR_PID=""
MOUNT_PID=""

cleanup() {
    echo -e "\n${YELLOW}🧹 Cleaning up processes...${NC}"
    
    # Unmount filesystem
    if mountpoint -q "$MOUNT_POINT" 2>/dev/null; then
        echo "📤 Unmounting $MOUNT_POINT..."
        fusermount -u "$MOUNT_POINT" 2>/dev/null || umount "$MOUNT_POINT" 2>/dev/null || true
        sleep 1
    fi
    
    # Kill processes
    for pid in $MOUNT_PID $SIDECAR_PID $FILER_PID $VOLUME_PID $MASTER_PID; do
        if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
            echo "🔪 Killing process $pid..."
            kill "$pid" 2>/dev/null || true
        fi
    done
    
    # Wait for processes to exit
    sleep 2
    
    # Force kill if necessary
    for pid in $MOUNT_PID $SIDECAR_PID $FILER_PID $VOLUME_PID $MASTER_PID; do
        if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
            echo "💀 Force killing process $pid..."
            kill -9 "$pid" 2>/dev/null || true
        fi
    done
    
    # Clean up mount point
    if [[ -d "$MOUNT_POINT" ]]; then
        rmdir "$MOUNT_POINT" 2>/dev/null || true
    fi
    
    echo -e "${GREEN}✅ Cleanup complete${NC}"
}

trap cleanup EXIT

wait_for_service() {
    local name=$1
    local url=$2
    local max_attempts=30
    local attempt=1
    
    echo -e "${BLUE}⏳ Waiting for $name to be ready...${NC}"
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -s "$url" >/dev/null 2>&1; then
            echo -e "${GREEN}✅ $name is ready${NC}"
            return 0
        fi
        echo "   Attempt $attempt/$max_attempts..."
        sleep 1
        ((attempt++))
    done
    
    echo -e "${RED}❌ $name failed to start within $max_attempts seconds${NC}"
    return 1
}

echo -e "${BLUE}🚀 SEAWEEDFS RDMA MOUNT DEMONSTRATION${NC}"
echo "======================================"
echo ""
echo "This demo shows Hanzo mount with RDMA acceleration:"
echo "  • Standard Hanzo cluster (master, volume, filer)"
echo "  • RDMA sidecar for acceleration"
echo "  • FUSE mount with RDMA fast path"
echo "  • Performance comparison tests"
echo ""

# Create mount point
echo -e "${BLUE}📁 Creating mount point: $MOUNT_POINT${NC}"
mkdir -p "$MOUNT_POINT"

# Start Hanzo Master
echo -e "${BLUE}🎯 Starting Hanzo Master...${NC}"
cd "$SEAWEEDFS_DIR"
./s3 master -port=9333 -mdir=/tmp/hanzo-master &
MASTER_PID=$!
wait_for_service "Master" "http://localhost:9333/cluster/status"

# Start Hanzo Volume Server
echo -e "${BLUE}💾 Starting Hanzo Volume Server...${NC}"
./s3 volume -master=localhost:9333 -port=8080 -dir=/tmp/hanzo-volume &
VOLUME_PID=$!
wait_for_service "Volume Server" "http://localhost:8080/status"

# Start Hanzo Filer
echo -e "${BLUE}📂 Starting Hanzo Filer...${NC}"
./s3 filer -master=localhost:9333 -port=8888 &
FILER_PID=$!
wait_for_service "Filer" "http://localhost:8888/"

# Start RDMA Sidecar
echo -e "${BLUE}⚡ Starting RDMA Sidecar...${NC}"
cd "$SIDECAR_DIR"
./bin/demo-server --port 8081 --rdma-socket /tmp/rdma-engine.sock --volume-server-url http://localhost:8080 --enable-rdma --debug &
SIDECAR_PID=$!
wait_for_service "RDMA Sidecar" "http://localhost:8081/health"

# Check RDMA capabilities
echo -e "${BLUE}🔍 Checking RDMA capabilities...${NC}"
curl -s "http://localhost:8081/stats" | jq . || curl -s "http://localhost:8081/stats"

echo ""
echo -e "${BLUE}🗂️  Mounting Hanzo with RDMA acceleration...${NC}"

# Mount with RDMA acceleration
cd "$SEAWEEDFS_DIR"
./s3 mount \
    -filer="$FILER_ADDR" \
    -dir="$MOUNT_POINT" \
    -rdma.enabled=true \
    -rdma.sidecar="$SIDECAR_ADDR" \
    -rdma.fallback=true \
    -rdma.maxConcurrent=64 \
    -rdma.timeoutMs=5000 \
    -debug=true &
MOUNT_PID=$!

# Wait for mount to be ready
echo -e "${BLUE}⏳ Waiting for mount to be ready...${NC}"
sleep 5

# Check if mount is successful
if ! mountpoint -q "$MOUNT_POINT"; then
    echo -e "${RED}❌ Mount failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Hanzo mounted successfully with RDMA acceleration!${NC}"
echo ""

# Demonstrate RDMA-accelerated operations
echo -e "${BLUE}🧪 TESTING RDMA-ACCELERATED FILE OPERATIONS${NC}"
echo "=============================================="

# Create test files
echo -e "${BLUE}📝 Creating test files...${NC}"
echo "Hello, RDMA World!" > "$MOUNT_POINT/test1.txt"
echo "This file will be read via RDMA acceleration!" > "$MOUNT_POINT/test2.txt"

# Create a larger test file
echo -e "${BLUE}📝 Creating larger test file (1MB)...${NC}"
dd if=/dev/zero of="$MOUNT_POINT/large_test.dat" bs=1024 count=1024 2>/dev/null

echo -e "${GREEN}✅ Test files created${NC}"
echo ""

# Test file reads
echo -e "${BLUE}📖 Testing file reads (should use RDMA fast path)...${NC}"
echo ""

echo "📄 Reading test1.txt:"
cat "$MOUNT_POINT/test1.txt"
echo ""

echo "📄 Reading test2.txt:"
cat "$MOUNT_POINT/test2.txt"
echo ""

echo "📄 Reading first 100 bytes of large file:"
head -c 100 "$MOUNT_POINT/large_test.dat" | hexdump -C | head -5
echo ""

# Performance test
echo -e "${BLUE}🏁 PERFORMANCE COMPARISON${NC}"
echo "========================="

echo "🔥 Testing read performance with RDMA acceleration..."
time_start=$(date +%s%N)
for i in {1..10}; do
    cat "$MOUNT_POINT/large_test.dat" > /dev/null
done
time_end=$(date +%s%N)
rdma_time=$((($time_end - $time_start) / 1000000)) # Convert to milliseconds

echo "✅ RDMA-accelerated reads: 10 x 1MB file = ${rdma_time}ms total"
echo ""

# Check RDMA statistics
echo -e "${BLUE}📊 RDMA Statistics:${NC}"
curl -s "http://localhost:8081/stats" | jq . 2>/dev/null || curl -s "http://localhost:8081/stats"
echo ""

# List files
echo -e "${BLUE}📋 Files in mounted filesystem:${NC}"
ls -la "$MOUNT_POINT/"
echo ""

# Interactive mode
echo -e "${BLUE}🎮 INTERACTIVE MODE${NC}"
echo "=================="
echo ""
echo "The Hanzo filesystem is now mounted at: $MOUNT_POINT"
echo "RDMA acceleration is active for all read operations!"
echo ""
echo "Try these commands:"
echo "  ls $MOUNT_POINT/"
echo "  cat $MOUNT_POINT/test1.txt"
echo "  echo 'New content' > $MOUNT_POINT/new_file.txt"
echo "  cat $MOUNT_POINT/new_file.txt"
echo ""
echo "Monitor RDMA stats: curl http://localhost:8081/stats | jq"
echo "Check mount status:  mount | grep hanzo"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop the demo and cleanup${NC}"

# Keep running until interrupted
while true; do
    sleep 5
    
    # Check if mount is still active
    if ! mountpoint -q "$MOUNT_POINT"; then
        echo -e "${RED}❌ Mount point lost, exiting...${NC}"
        break
    fi
    
    # Show periodic stats
    echo -e "${BLUE}📊 Current RDMA stats ($(date)):${NC}"
    curl -s "http://localhost:8081/stats" | jq '.rdma_enabled, .total_reads, .rdma_reads, .http_fallbacks' 2>/dev/null || echo "Stats unavailable"
    echo ""
done
