#!/bin/bash

# Wait for services to be ready
set -e

OPENBAO_ADDR=${OPENBAO_ADDR:-"http://127.0.0.1:8200"}
SEAWEEDFS_S3_ENDPOINT=${SEAWEEDFS_S3_ENDPOINT:-"http://127.0.0.1:8333"}
MAX_WAIT=120 # 2 minutes

echo "🕐 Waiting for services to be ready..."

# Wait for OpenBao
echo "   Waiting for OpenBao at $OPENBAO_ADDR..."
for i in $(seq 1 $MAX_WAIT); do
    if curl -s "$OPENBAO_ADDR/v1/sys/health" >/dev/null 2>&1; then
        echo "   [OK] OpenBao is ready!"
        break
    fi
    if [ $i -eq $MAX_WAIT ]; then
        echo "   [FAIL] Timeout waiting for OpenBao"
        exit 1
    fi
    sleep 1
done

# Wait for Hanzo Master
echo "   Waiting for Hanzo Master at http://127.0.0.1:9333..."
for i in $(seq 1 $MAX_WAIT); do
    if curl -s "http://127.0.0.1:9333/cluster/status" >/dev/null 2>&1; then
        echo "   [OK] Hanzo Master is ready!"
        break
    fi
    if [ $i -eq $MAX_WAIT ]; then
        echo "   [FAIL] Timeout waiting for Hanzo Master"
        exit 1
    fi
    sleep 1
done

# Wait for Hanzo Volume Server
echo "   Waiting for Hanzo Volume Server at http://127.0.0.1:8080..."
for i in $(seq 1 $MAX_WAIT); do
    if curl -s "http://127.0.0.1:8080/status" >/dev/null 2>&1; then
        echo "   [OK] Hanzo Volume Server is ready!"
        break
    fi
    if [ $i -eq $MAX_WAIT ]; then
        echo "   [FAIL] Timeout waiting for Hanzo Volume Server"
        exit 1
    fi
    sleep 1
done

# Wait for Hanzo S3 API
echo "   Waiting for Hanzo S3 API at $SEAWEEDFS_S3_ENDPOINT..."
for i in $(seq 1 $MAX_WAIT); do
    if curl -s "$SEAWEEDFS_S3_ENDPOINT/" >/dev/null 2>&1; then
        echo "   [OK] Hanzo S3 API is ready!"
        break
    fi
    if [ $i -eq $MAX_WAIT ]; then
        echo "   [FAIL] Timeout waiting for Hanzo S3 API"
        exit 1
    fi
    sleep 1
done

echo "🎉 All services are ready!"

# Show service status
echo ""
echo "📊 Service Status:"
echo "   OpenBao:             $(curl -s $OPENBAO_ADDR/v1/sys/health | jq -r '.initialized // "Unknown"')"
echo "   Hanzo Master:    $(curl -s http://127.0.0.1:9333/cluster/status | jq -r '.IsLeader // "Unknown"')"
echo "   Hanzo Volume:    $(curl -s http://127.0.0.1:8080/status | jq -r '.Version // "Unknown"')"
echo "   Hanzo S3 API:    Ready"
echo ""
