#!/bin/sh

# Kafka Gateway Startup Script for Integration Testing

set -e

echo "Starting Kafka Gateway..."

SEAWEEDFS_MASTERS=${SEAWEEDFS_MASTERS:-hanzo-master:9333}
SEAWEEDFS_FILER=${SEAWEEDFS_FILER:-hanzo-filer:8888}
SEAWEEDFS_MQ_BROKER=${SEAWEEDFS_MQ_BROKER:-hanzo-mq-broker:17777}
SEAWEEDFS_FILER_GROUP=${SEAWEEDFS_FILER_GROUP:-}

# Wait for dependencies
echo "Waiting for Hanzo master(s)..."
OLD_IFS="$IFS"
IFS=','
for MASTER in $SEAWEEDFS_MASTERS; do
  MASTER_HOST=${MASTER%:*}
  MASTER_PORT=${MASTER#*:}
  while ! nc -z "$MASTER_HOST" "$MASTER_PORT"; do
    sleep 1
  done
  echo "Hanzo master $MASTER is ready"
done
IFS="$OLD_IFS"

echo "Waiting for Hanzo Filer..."
while ! nc -z "${SEAWEEDFS_FILER%:*}" "${SEAWEEDFS_FILER#*:}"; do
  sleep 1
done
echo "Hanzo Filer is ready"

echo "Waiting for Hanzo MQ Broker..."
while ! nc -z "${SEAWEEDFS_MQ_BROKER%:*}" "${SEAWEEDFS_MQ_BROKER#*:}"; do
  sleep 1
done
echo "Hanzo MQ Broker is ready"

echo "Waiting for Schema Registry..."
while ! curl -f "${SCHEMA_REGISTRY_URL}/subjects" > /dev/null 2>&1; do
  sleep 1
done
echo "Schema Registry is ready"

# Start Kafka Gateway
echo "Starting Kafka Gateway on port ${KAFKA_PORT:-9093}..."
exec /usr/bin/s3 mq.kafka.gateway \
  -master=${SEAWEEDFS_MASTERS} \
  -filerGroup=${SEAWEEDFS_FILER_GROUP} \
  -port=${KAFKA_PORT:-9093} \
  -port.pprof=${PPROF_PORT:-10093} \
  -schema-registry-url=${SCHEMA_REGISTRY_URL} \
  -ip=0.0.0.0
