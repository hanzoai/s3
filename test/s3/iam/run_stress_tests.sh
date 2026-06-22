#!/bin/bash

# Stress Test Runner for Hanzo S3 IAM

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}💪 Running S3 IAM Stress Tests${NC}"

# Enable stress tests
export ENABLE_STRESS_TESTS=true
export TEST_TIMEOUT=60m

# Run stress tests multiple times
STRESS_ITERATIONS=5

echo -e "${YELLOW}🔄 Running stress tests with $STRESS_ITERATIONS iterations...${NC}"

for i in $(seq 1 $STRESS_ITERATIONS); do
    echo -e "${YELLOW}📊 Iteration $i/$STRESS_ITERATIONS${NC}"
    
    if ! go test -v -timeout=$TEST_TIMEOUT -run "TestS3IAMDistributedTests.*concurrent" ./... -count=1; then
        echo -e "${RED}❌ Stress test failed on iteration $i${NC}"
        exit 1
    fi
    
    # Brief pause between iterations
    sleep 2
done

echo -e "${GREEN}[OK] All stress test iterations completed successfully${NC}"
