# FoundationDB Filer Store Installation Guide

This guide covers the installation and setup of the FoundationDB filer store for Hanzo.

## Prerequisites

### FoundationDB Server

1. **Install FoundationDB Server**
   
   **Ubuntu/Debian:**
   ```bash
   # Add FoundationDB repository
   curl -L https://github.com/apple/foundationdb/releases/download/7.4.5/foundationdb-clients_7.4.5-1_amd64.deb -o foundationdb-clients.deb
   curl -L https://github.com/apple/foundationdb/releases/download/7.4.5/foundationdb-server_7.4.5-1_amd64.deb -o foundationdb-server.deb
   
   sudo dpkg -i foundationdb-clients.deb foundationdb-server.deb
   ```
   
   **CentOS/RHEL:**
   ```bash
   # Install RPM packages
   wget https://github.com/apple/foundationdb/releases/download/7.4.5/foundationdb-clients-7.4.5-1.el7.x86_64.rpm
   wget https://github.com/apple/foundationdb/releases/download/7.4.5/foundationdb-server-7.4.5-1.el7.x86_64.rpm
   
   sudo rpm -Uvh foundationdb-clients-7.4.5-1.el7.x86_64.rpm foundationdb-server-7.4.5-1.el7.x86_64.rpm
   ```
   
   **macOS:**
   ```bash
   # Using Homebrew (if available)
   brew install foundationdb
   
   # Or download from GitHub releases
   # https://github.com/apple/foundationdb/releases
   ```

2. **Initialize FoundationDB Cluster**
   
   **Single Node (Development):**
   ```bash
   # Start FoundationDB service
   sudo systemctl start foundationdb
   sudo systemctl enable foundationdb
   
   # Initialize database
   fdbcli --exec 'configure new single ssd'
   ```
   
   **Multi-Node Cluster (Production):**
   ```bash
   # On each node, edit /etc/foundationdb/fdb.cluster
   # Example: testing:testing@node1:4500,node2:4500,node3:4500
   
   # On one node, initialize cluster
   fdbcli --exec 'configure new double ssd'
   ```

3. **Verify Installation**
   ```bash
   fdbcli --exec 'status'
   ```

### FoundationDB Client Libraries

The Hanzo FoundationDB integration requires the FoundationDB client libraries.

**Ubuntu/Debian:**
```bash
sudo apt-get install libfdb-dev
```

**CentOS/RHEL:**
```bash
sudo yum install foundationdb-devel
```

**macOS:**
```bash
# Client libraries are included with the server installation
export LIBRARY_PATH=/usr/local/lib
export CPATH=/usr/local/include
```

## Building Hanzo with FoundationDB Support

### Download FoundationDB Go Bindings

```bash
go mod init hanzo-foundationdb
go get github.com/apple/foundationdb/bindings/go/src/fdb
```

### Build Hanzo

```bash
# Clone Hanzo repository
git clone https://github.com/hanzoai/s3.git
cd hanzo

# Build with FoundationDB support
go build -tags foundationdb -o s3
```

### Verify Build

```bash
./s3 version
# Should show version information

./s3 help
# Should list available commands
```

## Configuration

### Basic Configuration

Create or edit `filer.toml`:

```toml
[foundationdb]
enabled = true
cluster_file = "/etc/foundationdb/fdb.cluster"
api_version = 740
timeout = "5s"
max_retry_delay = "1s"
directory_prefix = "hanzo"
```

### Environment Variables

Alternative configuration via environment variables:

```bash
export WEED_FOUNDATIONDB_ENABLED=true
export WEED_FOUNDATIONDB_CLUSTER_FILE=/etc/foundationdb/fdb.cluster
export WEED_FOUNDATIONDB_API_VERSION=740
export WEED_FOUNDATIONDB_TIMEOUT=5s
export WEED_FOUNDATIONDB_MAX_RETRY_DELAY=1s
export WEED_FOUNDATIONDB_DIRECTORY_PREFIX=hanzo
```

### Advanced Configuration

For production deployments:

```toml
[foundationdb]
enabled = true
cluster_file = "/etc/foundationdb/fdb.cluster"
api_version = 740
timeout = "30s"
max_retry_delay = "5s"
directory_prefix = "hanzo_prod"

# Path-specific configuration for backups
[foundationdb.backup]
enabled = true
cluster_file = "/etc/foundationdb/fdb.cluster"
directory_prefix = "hanzo_backup"
location = "/backup"
timeout = "60s"
```

## Deployment

### Single Node Deployment

```bash
# Start Hanzo with FoundationDB filer
./s3 server -filer \
  -master.port=9333 \
  -volume.port=8080 \
  -filer.port=8888 \
  -s3.port=8333
```

### Distributed Deployment

**Master Servers:**
```bash
# Node 1
./s3 master -port=9333 -peers=master1:9333,master2:9333,master3:9333

# Node 2  
./s3 master -port=9333 -peers=master1:9333,master2:9333,master3:9333 -ip=master2

# Node 3
./s3 master -port=9333 -peers=master1:9333,master2:9333,master3:9333 -ip=master3
```

**Filer Servers with FoundationDB:**
```bash
# Filer nodes
./s3 filer -master=master1:9333,master2:9333,master3:9333 -port=8888
```

**Volume Servers:**
```bash
./s3 volume -master=master1:9333,master2:9333,master3:9333 -port=8080
```

### Docker Deployment

**docker-compose.yml:**
```yaml
version: '3.9'
services:
  foundationdb:
    image: foundationdb/foundationdb:7.4.5
    ports:
      - "4500:4500"
    volumes:
      - fdb_data:/var/fdb/data
      - fdb_config:/var/fdb/config

  hanzo:
    image: chrislusf/hanzo:latest
    command: "server -filer -ip=hanzo"
    ports:
      - "9333:9333"
      - "8888:8888"
      - "8333:8333"
    environment:
      WEED_FOUNDATIONDB_ENABLED: "true"
      WEED_FOUNDATIONDB_CLUSTER_FILE: "/var/fdb/config/fdb.cluster"
    volumes:
      - fdb_config:/var/fdb/config
    depends_on:
      - foundationdb

volumes:
  fdb_data:
  fdb_config:
```

### Kubernetes Deployment

**FoundationDB Operator:**
```bash
# Install FoundationDB operator
kubectl apply -f https://raw.githubusercontent.com/FoundationDB/fdb-kubernetes-operator/main/config/samples/deployment.yaml
```

**Hanzo with FoundationDB:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hanzo-filer
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hanzo-filer
  template:
    metadata:
      labels:
        app: hanzo-filer
    spec:
      containers:
      - name: hanzo
        image: chrislusf/hanzo:latest
        command: ["s3", "filer"]
        env:
        - name: WEED_FOUNDATIONDB_ENABLED
          value: "true"
        - name: WEED_FOUNDATIONDB_CLUSTER_FILE
          value: "/var/fdb/config/cluster_file"
        ports:
        - containerPort: 8888
        volumeMounts:
        - name: fdb-config
          mountPath: /var/fdb/config
      volumes:
      - name: fdb-config
        configMap:
          name: fdb-cluster-config
```

## Testing Installation

### Quick Test

```bash
# Start Hanzo with FoundationDB
./s3 server -filer &

# Test file operations
echo "Hello FoundationDB" > test.txt
curl -F file=@test.txt "http://localhost:8888/test/"
curl "http://localhost:8888/test/test.txt"

# Test S3 API
curl -X PUT "http://localhost:8333/testbucket"
curl -T test.txt "http://localhost:8333/testbucket/test.txt"
```

### Integration Test Suite

```bash
# Run the provided test suite
cd test/foundationdb
make setup
make test
```

## Performance Tuning

### FoundationDB Tuning

```bash
# Configure for high performance
fdbcli --exec 'configure triple ssd'
fdbcli --exec 'configure storage_engine=ssd-redwood-1-experimental'
```

### Hanzo Configuration

```toml
[foundationdb]
enabled = true
cluster_file = "/etc/foundationdb/fdb.cluster"
timeout = "10s"           # Longer timeout for large operations
max_retry_delay = "2s"    # Adjust retry behavior
directory_prefix = "sw"   # Shorter prefix for efficiency
```

### OS-Level Tuning

```bash
# Increase file descriptor limits
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf

# Adjust network parameters
echo "net.core.rmem_max = 134217728" >> /etc/sysctl.conf
echo "net.core.wmem_max = 134217728" >> /etc/sysctl.conf
sysctl -p
```

## Monitoring and Maintenance

### Health Checks

```bash
# FoundationDB cluster health
fdbcli --exec 'status'
fdbcli --exec 'status details'

# Hanzo health
curl http://localhost:9333/cluster/status
curl http://localhost:8888/statistics/health
```

### Log Monitoring

**FoundationDB Logs:**
- `/var/log/foundationdb/` (default location)
- Monitor for errors, warnings, and performance issues

**Hanzo Logs:**
```bash
# Start with verbose logging
./s3 -v=2 server -filer
```

### Backup and Recovery

**FoundationDB Backup:**
```bash
# Start backup
fdbbackup start -d file:///path/to/backup -t backup_tag

# Monitor backup
fdbbackup status -t backup_tag

# Restore from backup
fdbrestore start -r file:///path/to/backup -t backup_tag --wait
```

**Hanzo Metadata Backup:**
```bash
# Export filer metadata
./s3 shell
> fs.meta.save /path/to/metadata/backup.gz
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Check FoundationDB service status: `sudo systemctl status foundationdb`
   - Verify cluster file: `cat /etc/foundationdb/fdb.cluster`
   - Check network connectivity: `telnet localhost 4500`

2. **API Version Mismatch**
   - Update API version in configuration
   - Rebuild Hanzo with matching FDB client library

3. **Transaction Conflicts**
   - Reduce transaction scope
   - Implement appropriate retry logic
   - Check for concurrent access patterns

4. **Performance Issues**
   - Monitor cluster status: `fdbcli --exec 'status details'`
   - Check data distribution: `fdbcli --exec 'status json'`
   - Verify storage configuration

### Debug Mode

```bash
# Enable FoundationDB client tracing
export FDB_TRACE_ENABLE=1
export FDB_TRACE_PATH=/tmp/fdb_trace

# Start Hanzo with debug logging
./s3 -v=3 server -filer
```

### Getting Help

1. **FoundationDB Documentation**: https://apple.github.io/foundationdb/
2. **Hanzo Community**: https://github.com/hanzoai/s3/discussions
3. **Issue Reporting**: https://github.com/hanzoai/s3/issues

For specific FoundationDB filer store issues, include:
- FoundationDB version and cluster configuration
- Hanzo version and build tags
- Configuration files (filer.toml)
- Error messages and logs
- Steps to reproduce the issue
