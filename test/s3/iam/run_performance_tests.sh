#!/bin/bash

# Performance Test Runner for Hanzo S3 IAM

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🏁 Running S3 IAM Performance Tests${NC}"

# Enable performance tests
export ENABLE_PERFORMANCE_TESTS=true
export TEST_TIMEOUT=60m

# Run benchmarks
echo -e "${YELLOW}📊 Running benchmarks...${NC}"
go test -bench=. -benchmem -timeout=$TEST_TIMEOUT ./...

# Run performance tests
echo -e "${YELLOW}🧪 Running performance test suite...${NC}"
go test -v -timeout=$TEST_TIMEOUT -run "TestS3IAMPerformanceTests" ./...

echo -e "${GREEN}[OK] Performance tests completed${NC}"
